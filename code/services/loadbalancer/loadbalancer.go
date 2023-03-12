package loadbalancer

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type API struct {
	Key       string
	Times     uint32
	Available bool
}

type LoadBalancer struct {
	apis []*API
	mu   sync.RWMutex
}

func NewLoadBalancer(keys []string) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, key := range keys {
		lb.apis = append(lb.apis, &API{Key: key})
	}
	//SetAvailabilityForAll true
	lb.SetAvailabilityForAll(true)
	return lb
}

func (lb *LoadBalancer) GetAPI() *API {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	var availableAPIs []*API
	for _, api := range lb.apis {
		if api.Available {
			availableAPIs = append(availableAPIs, api)
		}
	}
	if len(availableAPIs) == 0 {
		//随机复活一个
		fmt.Printf("No available API, revive one randomly\n")
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(lb.apis))
		lb.apis[index].Available = true
		return lb.apis[index]
	}

	selectedAPI := availableAPIs[0]
	minTimes := selectedAPI.Times
	for _, api := range availableAPIs {
		if api.Times < minTimes {
			selectedAPI = api
			minTimes = api.Times
		}
	}
	selectedAPI.Times++
	//fmt.Printf("API Availability:\n")
	//for _, api := range lb.apis {
	//	fmt.Printf("%s: %v\n", api.Key, api.Available)
	//	fmt.Printf("%s: %d\n", api.Key, api.Times)
	//}

	return selectedAPI
}
func (lb *LoadBalancer) SetAvailability(key string, available bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, api := range lb.apis {
		if api.Key == key {
			api.Available = available
			return
		}
	}
}

func (lb *LoadBalancer) RegisterAPI(key string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.apis == nil {
		lb.apis = make([]*API, 0)
	}

	lb.apis = append(lb.apis, &API{Key: key})
}

func (lb *LoadBalancer) SetAvailabilityForAll(available bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, api := range lb.apis {
		api.Available = available
	}
}

func (lb *LoadBalancer) GetAPIs() []*API {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	apis := make([]*API, len(lb.apis))
	copy(apis, lb.apis)
	return apis
}
