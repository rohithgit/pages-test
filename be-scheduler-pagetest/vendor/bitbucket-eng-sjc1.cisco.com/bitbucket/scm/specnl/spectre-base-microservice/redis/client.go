package redis

import (
	"errors"

	"strings"

	"github.com/garyburd/redigo/redis"
	// "fmt"
)

var pool *redis.Pool

type RedisCache struct {
	conn *redis.Conn
	serviceName string
	customerId string
}

func newPool( connStr string) (*redis.Pool) {
	return &redis.Pool{
    MaxIdle: 80,
    MaxActive: 12000, // max number of connections
    Dial: func() (redis.Conn, error) {
			// fmt.Println( "connStr ", connStr)
      c, err := redis.Dial("tcp", connStr)
      // if err != nil {
      //   panic(err.Error())
      // }
      return c, err
    },
  }
}

func createPool( connStr string) (*redis.Pool, error){
	if strings.TrimSpace(connStr) == "" {
		return nil, errors.New("Connection string is empty")
	}

	if pool == nil {
		pool = newPool(connStr)
	}

	return pool, nil
}

func CreateConn( connStr, serviceName string) (connCache *RedisCache, err error){
	if strings.TrimSpace(connStr) == "" {
		return nil, errors.New("Connection string is empty")
	}
	connPool,err := createPool(connStr)
	if err != nil {
		return nil, err
	}
	if connPool == nil  {
		return nil, errors.New("Redis connection pool is not active")
	}
	connCache = new(RedisCache)
	conn := connPool.Get()
	connCache.conn = &conn
	connCache.serviceName = serviceName

	if connCache.conn == nil  {
		return nil, errors.New("Redis connection is not active")
	}
	return connCache, nil
}

func (connCache *RedisCache) GetHash( key string) (map[string]string, error){
	if connCache == nil {
		return nil, errors.New("Cache is empty")
	}

	conn := *connCache.conn
	if key == "" {
		// fmt.Println("No key provided to get value")
		return nil, errors.New("Key is empty")
	}
	// var out map[string]string
	_, err := redis.Values(conn.Do("HGETALL", key))
	// if  err != nil {
	// 	fmt.Println("GETHash ==================", err)
	// } else{
	// 	fmt.Println("Value GET %T", value)
	// }
	return nil, err
}

func (connCache *RedisCache)  SetHash( key string, values map[string]string) error{
/*	fmt.Println("inside SetHash")
	// conn, err := createConn()
	conn := *connCache.conn
	if key == "" {
		fmt.Println("No key provided to get value")
		return errors.New("Key is empty")
	}
	var args []interface{}{string}
	for k, v := range values {
		fmt.Println( "K, V = ", k,v)
		args = append(args, k, v)
	}

	if _, err := conn.Do("HMSET", key, args); err != nil {
		fmt.Println("SET  Hash==================", err)
		return err
	}*/
	return nil
}

func modifiedKey( customerId, serviceName, key string) (string){

	modifiedKeyName  := "";

	if strings.TrimSpace(customerId) != "" {
		modifiedKeyName = customerId + "."
	}

	if strings.TrimSpace(serviceName) != "" {
		modifiedKeyName += serviceName + "."
	}

	modifiedKeyName += key
	// fmt.Println("modifiedKeyName = ", modifiedKeyName)
	return modifiedKeyName;
}

// This function Gets value from hazelcast cache depending on key that is passed
func (connCache *RedisCache) GetString(key string) (string, error) {

	if connCache == nil {
		return "", errors.New("Cache is empty")
	}

	conn := *connCache.conn
	if key == "" {
		// fmt.Println("No key provided to get value")
		return "", errors.New("Key is empty")
	}

	// var out string
	value, err := redis.String(conn.Do("GET", modifiedKey(connCache.customerId, connCache.serviceName, key)))
	return value, err
}

func (connCache *RedisCache)  SetString( key string, value string) (error){
	// fmt.Println("inside SetString")
	// conn, err := createConn()
	if connCache == nil {
		return errors.New("Cache is empty")
	}
	conn := *connCache.conn
	if key == "" {
		// fmt.Println("No key provided to get value")
		return errors.New("Key is empty")
	}
	if value == "" {
		return errors.New("Value cannot be empty")
	}
	if _, err := conn.Do("SET", modifiedKey(connCache.customerId, connCache.serviceName, key), value); err != nil {
		// fmt.Println("SET ==================", err)
		return err
	}
	return nil
}

func (connCache *RedisCache)  GetKeys( keyPattern string) (values []string, err error){
	if connCache == nil {
		return nil, errors.New("Cache is empty")
	}

	conn := *connCache.conn
	// defer conn.Close()
	if keyPattern == "" {
		// fmt.Println("No keyPattern provided to get value")
		return nil, errors.New("keyPattern is empty")
	}

	values, err = redis.Strings(conn.Do("KEYS", "*"+modifiedKey(connCache.customerId, connCache.serviceName, keyPattern)+"*"))
	// if  err != nil {
	// 	fmt.Println("GetKeys ==================", err)
	// } else{
	// 	// out = value
	// 	fmt.Println("Value GetKeys", values)
	// }
	return values, err
}

func (connCache *RedisCache) SetStringWithExpiration(key, value string, expiration int32) (error) {
	//return errors.New("Redis connection pool is not active")
	if connCache == nil {
		return errors.New("Cache is empty")
	}

	conn := *connCache.conn
	if key == "" {
		// fmt.Println("No key provided to get value")
		return errors.New("Key is empty")
	}

	if value == "" {
		return errors.New("Value cannot be empty")
	}

	modifiedKeyStr := modifiedKey(connCache.customerId, connCache.serviceName, key)
	conn.Send("MULTI")
	conn.Send("SET", modifiedKeyStr, value)
	conn.Send( "EXPIRE", modifiedKeyStr, expiration)

	if _, err := conn.Do("EXEC"); err != nil {
		// fmt.Println("SET ==================", err)
		return err
	}
	return nil
}

func (connCache *RedisCache)  Delete(key string) error {

	if connCache == nil {
		return errors.New("Cache is empty")
	}

	conn := *connCache.conn
	// defer conn.Close()
	if key == "" {
		// fmt.Println("No key provided to get value")
		return errors.New("Key is empty")
	}
	if _, err := conn.Do("DEL", modifiedKey(connCache.customerId, connCache.serviceName, key)); err != nil {
		// fmt.Println("DELETE ==================", err)
		return err
	}
	// fmt.Println("DELETE successful==================")
	return nil
}

func (connCache *RedisCache)  Lookup(key string) (bool, error) {

	if connCache == nil {
		return false, errors.New("Cache is empty")
	}

	_, err := connCache.GetString(key)
	// fmt.Println("Value = ", value)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (connCache *RedisCache)  GetTimeToLive( key string) (int, error){

	if connCache == nil {
		return -1, errors.New("Cache is empty")
	}

	// TTL mykey
	conn := *connCache.conn
	// defer conn.Close()
	if key == "" {
		// fmt.Println("No key provided to get value")
		return -1, errors.New("Key is empty")
	}
	// var out int
	value, err := redis.Int(conn.Do("TTL", modifiedKey(connCache.customerId, connCache.serviceName, key)))
	// if  err != nil {
	// 	fmt.Println("GetTimeToLive ==================", err)
	// } else{
	// 	out = value
	// 	fmt.Println("Value GetTimeToLive", value)
	// }
	return value, err
}


func (connCache *RedisCache)  ClearAll() error {
	if connCache == nil {
		return errors.New("Cache is empty")
	}

	conn := *connCache.conn
	if _, err := conn.Do("FLUSHDB"); err != nil {
		// fmt.Println("ClearAll ==================", err)
		return err
	}
	return nil
}

func (connCache *RedisCache)  Close()  {
	if connCache != nil {
		conn := *connCache.conn
		conn.Close()
	}
}

func (connCache *RedisCache)  SetCustomerId( customerId string) {
	if connCache == nil {
		connCache.customerId = customerId
	}

}

// func GetKeys( keyPattern string) error {
// 	fmt.Println("Indide look up");
// 	return errors.New("Redis connection pool is not active")
// }
