package server

import (
	"sync"
)

type Address struct {
	serviceName string
	addr        string
}

type nameStore struct {
	data       map[string]map[string]*Address
	dataLocker sync.RWMutex
}

var serviceNameData *nameStore

func init() {
	serviceNameData = &nameStore{
		data: map[string]map[string]*Address{},
	}

}

func Register(serviceName, address string) {
	ns := serviceNameData
	addr := &Address{
		serviceName: serviceName,
		addr:        address,
	}

	ns.dataLocker.Lock()
	_, ok := ns.data[serviceName]
	if !ok {
		ns.data[serviceName] = make(map[string]*Address, 0)
	}
	ns.data[serviceName][address] = addr
	ns.dataLocker.Unlock()

}

func GetAllData() *nameStore {
	return serviceNameData
}

// 根据服务名称获取地址信息
func GetByServiceName(serviceName string) []string {
	ns := serviceNameData
	ns.dataLocker.RLock()
	defer ns.dataLocker.RUnlock()
	if _, ok := ns.data[serviceName]; ok {
		address := make([]string, 0)
		for _, mapv := range ns.data[serviceName] {
			address = append(address, mapv.addr)
		}
		return address
	}

	return []string{}
}
