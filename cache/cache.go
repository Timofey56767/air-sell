package cache

import (
	"sync"

	"github.com/TOMMy-Net/air-sell/models"
)

var mutex sync.RWMutex


var airCache map[int]models.AirPorts

func AddAirport(a models.AirPorts) {
	mutex.Lock()
	defer mutex.Unlock()
	airCache[a.ID] = a
}

func GetAirport(i int) models.AirPorts {
	mutex.RLock()
	defer mutex.RUnlock()
	return airCache[i]
}

func DeleteAirport(i int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(airCache, i)
}

func NewCache() {
	airCache = make(map[int]models.AirPorts, 20)
}
