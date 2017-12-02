package models

import (
	"testing"

	"github.com/stretchr/testify/assert" //see: https://github.com/stretchr/testify
)

// TestNewProduct performs a unit test on NewProduct
func TestNewOptions(t *testing.T) {
	options := NewOptions()
	assert.NotNil(t, options, "NewOptions() returned nil")
}

func TestInitOptions(t *testing.T) {
	_, err := InitOptions()
	assert.NoError(t, err, "InitOptions() returned an error")
}
