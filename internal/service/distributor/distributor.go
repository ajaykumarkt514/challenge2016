package distributor

import (
	"challenge2016/internal/errors"
	"challenge2016/internal/model"
	"challenge2016/internal/store/cache"
	"fmt"
	"net/http"
	"strings"
)

func (c client) Add(distributor *model.Distributor) (*model.DistributorResponse, error) {
	// Validates distributor struct fields
	err := distributor.Validate()
	if err != nil {
		return nil, err
	}

	// Acquires lock to avoid race conditions
	cache.DistributorMutex.Lock()
	defer cache.DistributorMutex.Unlock()

	err = addDistributor(distributor)
	if err != nil {
		return nil, err
	}

	response := &model.DistributorResponse{
		Name:      distributor.Name,
		Locations: cache.DistributorsCache[distributor.Name],
	}

	return response, nil
}

func (c client) Get(name string) (*model.DistributorResponse, error) {
	cache.DistributorMutex.RLock()
	defer cache.DistributorMutex.RUnlock()

	if _, ok := cache.DistributorsCache[name]; !ok {
		return nil, &errors.Response{
			Code:   http.StatusNotFound,
			Reason: "distributor not found",
		}
	}

	response := &model.DistributorResponse{
		Name:      name,
		Locations: cache.DistributorsCache[name],
	}

	return response, nil
}

func (c client) CheckAccess(name string, region string) (string, error) {
	cache.DistributorMutex.RLock()
	defer cache.DistributorMutex.RUnlock()

	// Fetch distributor's permissions from cache
	distributorCache, distributorExists := cache.DistributorsCache[name]
	if !distributorExists {
		return "", &errors.Response{
			Code:   http.StatusNotFound,
			Reason: fmt.Sprintf("Distributor %s not found", name),
		}
	}

	// Parse region into components: country, province, city
	keys := strings.Split(region, "-")
	var country, province, city string

	switch len(keys) {
	case 3:
		city = model.Sanitize(keys[0])
		province = model.Sanitize(keys[1])
		country = model.Sanitize(keys[2])
	case 2:
		province = model.Sanitize(keys[0])
		country = model.Sanitize(keys[1])
	case 1:
		country = model.Sanitize(keys[0])
	default:
		return "", &errors.Response{
			Code:   http.StatusBadRequest,
			Reason: "Invalid region format",
		}
	}

	// Check if the region is included
	if isIncluded, err := isRegionIncluded(distributorCache, country, province, city); err != nil {
		return "", err
	} else if !isIncluded {
		return "NO", nil
	}

	return "YES", nil
}

// Helper to check if a region is included in the distributor's permissions
func isRegionIncluded(distributorCache map[string]*model.Country, country, province, city string) (bool, error) {
	// Check country
	countryCache, countryExists := distributorCache[country]
	if !countryExists {
		return false, nil
	}

	// Check province if specified
	if province != "" {
		provinceCache, provinceExists := countryCache.Provinces[province]
		if !provinceExists {
			return false, nil
		}

		// Check city if specified
		if city != "" {
			_, cityExists := provinceCache.Cities[city]
			if !cityExists {
				return false, nil
			}
		}
	}

	return true, nil
}

func addDistributor(distributor *model.Distributor) error {
	if _, ok := cache.DistributorsCache[distributor.Name]; ok {
		return &errors.Response{
			Code:   http.StatusBadRequest,
			Reason: "distributor already exists",
		}
	}

	err := validateAgainstParentDistributor(distributor)
	if err != nil {
		return err
	}

	err = addIncludeRegions(distributor)
	if err != nil {
		return err
	}

	err = removeExcludeRegions(distributor)
	if err != nil {
		return err
	}

	return nil
}

func validateAgainstParentDistributor(distributor *model.Distributor) error {
	// If there is no parent distributor, validation is not required
	if distributor.Parent == nil {
		return nil
	}

	parentName := *distributor.Parent
	parentPermissions, parentExists := cache.DistributorsCache[parentName]
	if !parentExists {
		return &errors.Response{
			Code:   http.StatusBadRequest,
			Reason: fmt.Sprintf("Parent distributor %s not found in cache", parentName),
		}
	}

	// Validate each included region against the parent distributor's permissions
	for _, includeRegion := range distributor.Include {
		if !isRegionWithinParentPermissions(includeRegion, parentPermissions) {
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: fmt.Sprintf("Region %s is outside the authorized area of parent distributor %s", includeRegion, parentName),
			}
		}
	}

	// Ensure excluded regions are valid within the parent's authorized regions
	for _, excludeRegion := range distributor.Exclude {
		if !isRegionWithinParentPermissions(excludeRegion, parentPermissions) {
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: fmt.Sprintf("Excluded region %s is outside the authorized area of parent distributor %s", excludeRegion, parentName),
			}
		}
	}

	return nil
}

// Helper function to validate if a region is within the parent distributor's permissions
func isRegionWithinParentPermissions(region string, parentPermissions map[string]*model.Country) bool {
	keys := strings.Split(region, "-")
	country := model.Sanitize(keys[0])

	parentCountry, countryExists := parentPermissions[country]
	if !countryExists {
		return false
	}

	if len(keys) > 1 {
		province := model.Sanitize(keys[1])
		parentProvince, provinceExists := parentCountry.Provinces[province]
		if !provinceExists {
			return false
		}

		if len(keys) > 2 {
			city := model.Sanitize(keys[2])
			_, cityExists := parentProvince.Cities[city]
			if !cityExists {
				return false
			}
		}
	}

	return true
}

func addIncludeRegions(distributor *model.Distributor) error {
	if len(distributor.Include) == 0 {
		return &errors.Response{
			Code:   http.StatusBadRequest,
			Reason: "Distributor doesn't have access to any region to start with",
		}
	}

	result := make(map[string]*model.Country)

	for _, value := range distributor.Include {
		var country, province, city string

		keys := strings.Split(value, "-")

		switch len(keys) {
		case 3:
			city = model.Sanitize(keys[0])
			province = model.Sanitize(keys[1])
			country = model.Sanitize(keys[2])
		case 2:
			province = model.Sanitize(keys[0])
			country = model.Sanitize(keys[1])
		case 1:
			country = model.Sanitize(keys[0])
		default:
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: "Invalid region format in exclude list",
			}
		}

		// Validate country
		c, countryExists := cache.Locations[country]
		if !countryExists {
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: fmt.Sprintf("Country %v not found", country),
			}
		}

		// Initialize result entry for the country if not already done
		if _, ok := result[country]; !ok {
			result[country] = &model.Country{
				Code:      c.Code,
				Provinces: make(map[string]*model.Province),
			}
		}

		if province != "" {
			// Validate province
			p, provinceExists := c.Provinces[province]
			if !provinceExists {
				return &errors.Response{
					Code:   http.StatusBadRequest,
					Reason: fmt.Sprintf("Province %v not found in country %v", province, country),
				}
			}

			// Initialize result entry for the province if not already done
			if _, ok := result[country].Provinces[province]; !ok {
				result[country].Provinces[province] = &model.Province{
					Code:   p.Code,
					Cities: make(map[string]*model.City),
				}
			}

			if city != "" {
				// Validate city
				cityObj, cityExists := p.Cities[city]
				if !cityExists {
					return &errors.Response{
						Code:   http.StatusBadRequest,
						Reason: fmt.Sprintf("City %v not found in province %v, country %v", city, province, country),
					}
				}
				// Add specific city
				result[country].Provinces[province].Cities[city] = cityObj
			} else {
				// Add all cities in the province
				for cityName, cityObj := range p.Cities {
					result[country].Provinces[province].Cities[cityName] = cityObj
				}
			}
		} else {
			// Add all provinces and cities in the country
			for provinceName, provinceObj := range c.Provinces {
				if _, ok := result[country].Provinces[provinceName]; !ok {
					result[country].Provinces[provinceName] = &model.Province{
						Code:   provinceObj.Code,
						Cities: make(map[string]*model.City),
					}
				}
				for cityName, cityObj := range provinceObj.Cities {
					result[country].Provinces[provinceName].Cities[cityName] = cityObj
				}
			}
		}
	}

	// Update cache with the validated regions
	cache.DistributorsCache[distributor.Name] = result

	return nil
}

func removeExcludeRegions(distributor *model.Distributor) error {
	// Get the current distributor's cache
	resultCache := cache.DistributorsCache[distributor.Name]

	for _, value := range distributor.Exclude {
		var country, province, city string

		keys := strings.Split(value, "-")
		switch len(keys) {
		case 3:
			city = model.Sanitize(keys[0])
			province = model.Sanitize(keys[1])
			country = model.Sanitize(keys[2])
		case 2:
			province = model.Sanitize(keys[0])
			country = model.Sanitize(keys[1])
		case 1:
			country = model.Sanitize(keys[0])
		default:
			return &errors.Response{
				Code:   http.StatusBadRequest,
				Reason: "Invalid region format in exclude list",
			}
		}

		// Check if country exists in current cache
		countryCache, countryExists := resultCache[country]
		if !countryExists {
			continue // Skip if country is not in the cache
		}

		if province != "" {
			// Check if province exists in country cache
			provinceCache, provinceExists := countryCache.Provinces[province]
			if !provinceExists {
				continue // Skip if province is not in the cache
			}

			if city != "" {
				// Check if city exists in province cache and delete it
				if _, cityExists := provinceCache.Cities[city]; cityExists {
					delete(provinceCache.Cities, city)
				}
			} else {
				// Delete the entire province if only country and province are specified
				delete(countryCache.Provinces, province)
			}
		} else {
			// Delete the entire country if only the country is specified
			delete(resultCache, country)
		}
	}

	// Update the distributor's cache
	cache.DistributorsCache[distributor.Name] = resultCache

	return nil
}
