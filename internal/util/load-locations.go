package util

import (
	"challenge2016/internal/model"
	"challenge2016/internal/store/cache"
	"encoding/csv"
	"gofr.dev/pkg/gofr/logging"
	"io"
	"os"
)

func LoadLocations(logger logging.Logger, file *os.File) error {
	reader := csv.NewReader(file)

	// Skip the header row
	_, err := reader.Read()
	if err != nil {
		logger.Errorf("Error while reading header: %v", err.Error())
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // Break the loop on EOF
			}

			return err
		}

		// Logs error if number of fields in row is not equals to 6
		if len(record) != 6 {
			logger.Errorf("Invalid number of fields in row: %v", record)
			continue
		}

		// Read fields from the CSV file
		cityCode := model.Sanitize(record[0]) // City Code
		city := model.Sanitize(record[3])     // City Name

		provinceCode := model.Sanitize(record[1]) // Province Code
		province := model.Sanitize(record[4])     // Province Name

		countryCode := model.Sanitize(record[2]) // Country Code
		country := model.Sanitize(record[5])     // Country Name

		// Adds new country to cache
		c, ok := cache.Locations[country]
		if !ok {
			c = &model.Country{
				Code:      countryCode,
				Provinces: make(map[string]*model.Province),
			}

			cache.Locations[country] = c
		}

		// Adds new provinces to cache
		p, ok := c.Provinces[province]
		if !ok {
			p = &model.Province{
				Code:   provinceCode,
				Cities: make(map[string]*model.City),
			}

			c.Provinces[province] = p
		}

		// Adds new city to cache
		ci, ok := p.Cities[city]
		if !ok {
			ci = &model.City{
				Code: cityCode,
			}

			p.Cities[city] = ci
		}
	}

	return nil
}
