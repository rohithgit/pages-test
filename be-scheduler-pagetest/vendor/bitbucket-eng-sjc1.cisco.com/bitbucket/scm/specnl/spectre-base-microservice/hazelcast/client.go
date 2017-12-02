package hazelcast

import (
	"errors"
	"strings"
	"github.com/bradfitz/gomemcache/memcache"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/logging"
)

// This is the struct used to store hazelcast client instance
type Cache struct {
	conn *memcache.Client
}

var Key string
var Value []byte
var log = logging.Log.Logger

// This function establishes connection with hazelcast cluster and returns necessary client and an error in case of a problem
func Connect(server ...string) (memc *Cache, err error) {
	a := len(server)
	if a == 0 {
		log.Infof("No hazelcast server to connect")
		return nil, errors.New("No server list provided")
	}
	memc = new(Cache)
	memc.conn = memcache.New(server...)
	if memc.conn == nil {
		return nil, errors.New("Hazelcast connection not established")
	}
	return memc, nil
}

// This function Gets value from hazelcast cache depending on key that is passed
func (memc *Cache) Get(key string) (string, error) {
	if memc == nil || memc.conn == nil {
		return "", errors.New("Hazelcast connection is not active")
	}
	if key == "" {
		log.Infof("No key provided to get value")
		return "", errors.New("Key is empty")
	}
	if item, err := memc.conn.Get(key); err == nil {
		log.Infof("Getting Key %s", key)
		log.Infof("Value %s successfully Get for key %s", string(item.Value), key)
		return string(item.Value), nil
	}
	return "", errors.New("No value found")
}

// This function Sets value to hazelcast cache, only supports string
func (memc *Cache) Set(key string, val string) error {
	item := memcache.Item{Key: key, Value: []byte(val)}
	log.Infof("Setting Key %s", key)
	return memc.SetItem(&item)
}

// This function Sets value to hazelcast cache with expiration, only supports string
func (memc *Cache) SetWithExpiration(key, val string, expiration int32) error {
	item := memcache.Item{Key: key, Value: []byte(val), Expiration: expiration}
	log.Infof("Setting Key %s", key)
	return memc.SetItem(&item)
}

// This function Sets value to hazelcast cache
func (memc *Cache) SetItem(item *memcache.Item) error {
	if memc == nil || memc.conn == nil {
		return errors.New("No active Hazelcast connection")
	}
	if item == nil || strings.TrimSpace(item.Key) == "" || len(item.Value) == 0 {
		return errors.New("Key and Value should be non-empty")
	}

	log.Infof("Setting Key %s", item.Key)
	return memc.conn.Set(item)
}

//Delete value from cache using the Key
func (memc *Cache) Delete(key string) error {
	if memc == nil || memc.conn == nil {
		return errors.New("No active Hazelcast connection")
	}
	if key == "" {
		return errors.New("No key to delete")
	}
	log.Infof("Deleting key %s", key)
	return memc.conn.Delete(key)
}

// This function checks whether the value exists in the  cache
func (memc *Cache) Lookup(key string) (bool, error) {
	if memc == nil || memc.conn == nil {
		return false, errors.New("No active Hazelcast connection")
	}
	if key == "" {
		return false, errors.New("No key to lookup")
	}
	log.Infof("Looking up for key %s", key)
	_, err := memc.conn.Get(key)
	if err != nil {
		return false, err
		log.Infof("Lookup Failed with error", err)
	}
	log.Infof("Lookup successful for key: %s ", key)
	return true, nil
}

// GetMulti get value from cache depending on keys that is passed
func (memc *Cache) GetMulti(keys []string) ([]string, error) {
	size := len(keys)
	var rv []string
	if memc == nil || memc.conn == nil {
		return nil, errors.New("No active Hazelcast connection")
	}
	if size < 1 {
		return nil, errors.New("No keys passed")
	}
	mv, err := memc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, string(v.Value))
		}
		log.Infof("Get on Multiple Keys successful")
		return rv, nil
	}
	log.Infof("Get on Multiple keys failed", err)
	return rv, err
}

// ClearAll function clears all data that is cached.
func (memc *Cache) ClearAll() error {
	if memc == nil || memc.conn == nil {

		return errors.New("No active Hazelcast connection")
	}
	log.Infof("Flushing all data from cache")
	return memc.conn.FlushAll()
}
