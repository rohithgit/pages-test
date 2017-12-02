// Package data contains data elements that you might need in other packages within the service. For example,
// when running this executable in stand-alone test mode you might need JSON test data.
package data

// DefaultOptions contains the default option values. These are used to
// initialize the Options object
var DefaultOptions = []byte(`
{
	"debug": false,
	"verbose": false,
	"port": 8000,
	"mock": false
}
`)
