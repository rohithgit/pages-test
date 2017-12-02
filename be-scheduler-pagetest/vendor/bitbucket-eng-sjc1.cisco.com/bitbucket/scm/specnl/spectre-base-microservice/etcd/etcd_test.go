package etcd

import (
	"os"
	"path"
	"strings"
	"testing"

	"logging"

	"github.com/stretchr/testify/require"
)

var validtestValues = map[string]string{
	"/key1": "value1",
}

var invalidtestValues = map[string]string{
	"key2": "value2",
}

var validDirValues = map[string]string{
	"/dir1/key1": "dval1",
	"/dir1/key2": "dval2",
}

var invalidDirValues = map[string]string{
	"/key1": "dval1",
}

// This is a global instance to store the etcd client
var SpectreEtcdClient *SpectreEtcd

var log = logging.Log.Logger

func init() {
	// Output to stderr instead of stdout, could also be a file.
	// Only log the warning severity or above.
	logging.Log.LoggingInit("Text", os.Stdout, "Info")

	var err error
	log.Info("Trying to establish etcd connection.")

	SpectreEtcdClient, err = NewETCDClient(
		[]string{"http://etcd0:2379"},
		//[]string{"http://localhost:4001"},
		nil)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("etcd connection established")
	}
}

func TestEtcdConnection(t *testing.T) {
	require := require.New(t)
	require.NotNil(SpectreEtcdClient.EtcdKeysClient)
}

func TestValidSetValue(t *testing.T) {
	require := require.New(t)
	for k, v := range validtestValues {
		_, err := SpectreEtcdClient.SetKeyValue(k, v)
		require.Nil(err, "Unable to set value. Error:", err)
	}
}

func TestInvalidSetValue(t *testing.T) {
	require := require.New(t)
	for k, v := range invalidtestValues {
		_, err := SpectreEtcdClient.SetKeyValue(k, v)
		require.NotNil(err, "Error should have occured trying to store key: %s with value: %s", k, v)
	}
}

func TestValidGetValue(t *testing.T) {
	require := require.New(t)
	for k, v := range validtestValues {
		val, err := SpectreEtcdClient.GetStringValue(k)
		require.Nil(err, "Unable to get value. Error:", err)
		require.Equal(v, *val, "Value obtained %s did not match what was set %s", *val, v)
	}
}

func TestInValidGetValue(t *testing.T) {
	require := require.New(t)
	for k, _ := range invalidtestValues {
		_, err := SpectreEtcdClient.GetStringValue(k)
		require.NotNil(err, "Error should have occured trying to retrieve key: %s", k)
	}
}

func TestValidDirSetValue(t *testing.T) {
	require := require.New(t)
	for k, v := range validDirValues {
		_, err := SpectreEtcdClient.SetKeyValue(k, v)
		require.Nil(err, "Unable to set value. Error:", err)
	}
}

func TestValidDirGetValue(t *testing.T) {
	require := require.New(t)
	for k, v := range validDirValues {
		val, err := SpectreEtcdClient.GetStringValue(k)
		require.Nil(err, "Unable to get value. Error:", err)
		require.Equal(v, *val, "Value obtained %s did not match what was set %s", *val, v)
	}
}

func TestInValidDirRecursive(t *testing.T) {
	require := require.New(t)
	for k, _ := range invalidDirValues {
		_, err := SpectreEtcdClient.GetDir(k)
		require.NotNil(err, "Error should have occured trying to retrieve key: %s", k)
	}
}

func TestValidDirRecursive(t *testing.T) {
	require := require.New(t)
	for k, _ := range validDirValues {
		dirName := strings.Replace(k, path.Base(k), "", -1)
		val, err := SpectreEtcdClient.GetDir(dirName)
		require.Nil(err, "Unable to get value. Error:", err)
		require.NotNil("Value obtained %v is not valid for key: %s", val, dirName)
	}
}

func deleteKeysFromMap(m map[string]string, require *require.Assertions) {
	for k, _ := range m {
		_, err := SpectreEtcdClient.DeleteKey(k)
		require.Nil(err, "Unable to delete value. Error:", err)
	}
}

func TestDelete(t *testing.T) {
	require := require.New(t)
	deleteKeysFromMap(validtestValues, require)
	deleteKeysFromMap(validDirValues, require)
}
