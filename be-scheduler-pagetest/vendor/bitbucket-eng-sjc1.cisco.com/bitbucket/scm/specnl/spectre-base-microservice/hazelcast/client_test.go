package hazelcast

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

var args []string = []string{"172.17.0.2:5701", "172.17.0.3:5701"}

func TestConnect(t *testing.T) {

	_, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}
}

func TestSet(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}
	err = c.Set(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}
}

func TestSetWithExpirationValidaCase(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}
	err = c.SetWithExpiration("testkey", "testvalue", 1)
	if err != nil {
		t.Error("Set Error", err)
	}
}

func TestSetWithExpirationCheckExpiration(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}
	err = c.SetWithExpiration("mkey", "mvalue", 1)
	if err != nil {
		t.Error("Set Error", err)
	}
	duration := time.Duration(1) * time.Second
	time.Sleep(duration)

	_, err = c.Get("mkey")
	if err == nil {
		t.Error("Get error", err)
	}

}

func TestGet(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}
	_, err = c.Get(KEY_1)
	if err != nil {
		t.Error("Get error", err)
	}
}

func TestDelete(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	err = c.Set(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}

	c.Delete(KEY_1)
	ok, err := c.Lookup(KEY_1)
	if ok == true {
		t.Error("Error Deleting", err)
	}
}

func TestLookup(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	err = c.Set(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set Error", err)
	}

	ok, err := c.Lookup(KEY_1)
	if ok != true {
		t.Error("Value not found", err)
	}
}

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

func TestClearAll(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	err = c.Set(KEY_1, VALUE_1)
	if err != nil {
		t.Error("Set error for KEY_1", err)
	}

	err = c.ClearAll()
	_, err = c.Get(KEY_1)
	if err != nil {
		t.Error("Value still not cleared", err)
	}
}

func TestEmptyConnect(t *testing.T) {

	_, err := Connect("")
	if err != nil {
		t.Error(err)
	}
}

func TestWrongValues(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	err = c.Set("", "")
	if err == nil {
		t.Error("Nil Value Set")
	}

	err = c.Set("foo", "")
	if err == nil {
		t.Error("Empty string Set for Value")
	}

	_, err = c.Get("")
	var ok bool = false
	ok, err = c.Lookup("")
	if ok == true {
		t.Error("Blank string found")
	}

	_, err = c.Get("foo")

	ok, err = c.Lookup("")
	if ok == true {
		t.Error("Blank string found")
	}

	c.Delete("")
	if err == nil {
		t.Error("Blank Key string found")
	}
	ok, err = c.Lookup("")
	if ok == true {
		t.Error("Error Deleting")
	}
	err = c.Set(KEY_3, VALUE_3)
	err = c.ClearAll()

	_, err = c.Get(KEY_3)
	if err != nil {
		t.Error("Value still not cleared", err)
	}
}

func TestNilGetMulti(t *testing.T) {

	c, err := Connect(args...)
	if err != nil {
		t.Error(err)
	}

	_, err = c.GetMulti(nil)
	if err == nil {
		t.Error("No Keys passed in MultiGet")
	}

}
