package cache

import (
	"challenge2016/internal/model"
	"sync"
)

var DistributorsCache = make(map[string]map[string]*model.Country)

var DistributorMutex sync.RWMutex
