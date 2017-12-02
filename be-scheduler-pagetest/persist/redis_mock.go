package persist

import (
	"errors"
)

type MockCache struct {
	vals   map[string]string
}

func NewRedisMock()  MockCache {
	return MockCache{}
}

func (cache *MockCache) GetString(key string) (string, error){
	if val := cache.vals[key]; val != "" {
		return val, nil
	}
	return "", errors.New("not found")
}

func (cache *MockCache) SetString(key string, val string) {
	cache.vals[key] = val
}

func (cache *MockCache) SetWithExpiration(key string, val string,  expiration int32) {
	cache.vals[key] = val
}
