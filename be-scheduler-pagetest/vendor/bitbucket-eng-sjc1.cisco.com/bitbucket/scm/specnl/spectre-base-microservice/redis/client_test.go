package redis

import (
	"testing"
	"time"
)

const (
	KEY_1   string = "testkey-1"
	KEY_2   string = "testkey-2"
	VALUE_1 string = "testval1"
	VALUE_2 string = "testval2"
	KEY_3   string = "testkey-3"
	VALUE_3 string = "testkey-3"
)

//var args []string = []string{"172.17.0.2:5701", "172.17.0.3:5701"}

var connectString string = "localhost:6379"
var serviceName string = "Collector"

func TestConnect(t *testing.T) {
	_,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}
}

func TestSetGetString(t *testing.T) {
	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}

	err = cache.SetString(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}

	value, err := cache.GetString(KEY_1)
	if err != nil {
		t.Error("Get Error", err)
	}

	if value == "" || value != VALUE_1 {
		t.Error("Get Error", err)
	}
}

func TestSetWithExpirationValidaCase(t *testing.T) {

	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}
	err = cache.SetStringWithExpiration("testkey", "testvalue", 1)
	if err != nil {
		t.Error("Set Error", err)
	}
}

func TestSetWithExpirationCheckExpiration(t *testing.T) {

	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}
	err = cache.SetStringWithExpiration("mkey", "mvalue", 1)
	if err != nil {
		t.Error("Set Error", err)
	}

	duration := time.Duration(1) * time.Second
	time.Sleep(duration)

	_, err = cache.GetString("mkey")
	if err == nil {
		t.Error("Get error", err)
	}

}

func TestDelete(t *testing.T) {

	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}

	err = cache.SetString(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}

	err = cache.Delete(KEY_1)
	ok, err := cache.Lookup(KEY_1)
	if (ok) {
		t.Error("Error Deleting", err)
	}
}

func TestLookup(t *testing.T) {

	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}

	err = cache.SetString(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}

	ok, err := cache.Lookup(KEY_1)
	if ok != true {
		t.Error("Value not found", err)
	}
}
/*
func TestGetMulti(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	err = c.Set(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set error for KEY_1", err)
	}

	// Check if KEY_1 is set
	ok, err := c.Lookup(KEY_1)
	if ok != true {
		t.Error("Lookup failed for KEY_1", err)
	}

	if err := c.Set(KEY_2, VALUE_2); err != nil {
		t.Error("Set error for KEY_2", err)
	}
	// Check if KEY_2 is set
	ok, err = c.Lookup(KEY_2)
	if ok != true {
		t.Error("Lookup failed for KEY_2", err)
	}

	var vv []string
	vv, err = c.GetMulti([]string{KEY_1, KEY_2})
	//fmt.Println(vv[0])
	//fmt.Println(vv[1])
	if len(vv) != 2 {
		t.Error("GetMulti Error", err)
	}
	if vv[0] != VALUE_1 {
		t.Error("GetMulti Error on KEY_1")
	}
	if vv[1] != VALUE_2 {
		t.Error("GetMulti Error on KEY_2")
	}

}
*/
func TestClearAll(t *testing.T) {

	cache,err := CreateConn(connectString, serviceName)
	if err != nil {
		t.Error(err)
	}

	err = cache.SetString(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set error for KEY_1", err)
	}

	err = cache.ClearAll()
	val, err := cache.GetString(KEY_1)
	if val != "" {
		t.Error("Value still not cleared", err)
	}
}

func TestEmptyConnect(t *testing.T) {

	_,err := CreateConn("", serviceName)
	if err == nil {
		t.Error(err)
	}
}

func TestWrongValues(t *testing.T) {

	cache,err := CreateConn("", serviceName)
	if err == nil {
		t.Error(err)
	}

	cache,err = CreateConn(connectString, serviceName)
	err = cache.SetString("", "")
	if err == nil {
		t.Error("Nil Value Set")
	}

	err = cache.SetString("foo", "")
	if err == nil {
		t.Error("Empty string Set for Value")
	}

	_, err = cache.GetString("")
	var ok bool = false
	ok, err = cache.Lookup("")
	if ok == true {
		t.Error("Blank string found")
	}

	_, err = cache.GetString("foo")

	ok, err = cache.Lookup("")
	if ok == true {
		t.Error("Blank string found")
	}

	err = cache.Delete("")
	if err == nil {
		t.Error("Blank Key string found")
	}
	ok, err = cache.Lookup("")
	if ok == true {
		t.Error("Error Deleting")
	}
	err = cache.SetString(KEY_3, VALUE_3)
	err = cache.ClearAll()

	_, err = cache.GetString(KEY_3)
	if err == nil {
		t.Error("Value still not cleared", err)
	}
}

/*func TestNilGetMulti(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	_, err = c.GetMulti(nil)
	if err == nil {
		t.Error("No Keys passed in MultiGet")
	}

}*/
