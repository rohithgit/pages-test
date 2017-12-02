package crypt

import (
	"testing"
	"github.com/stretchr/testify/require"
)

var data = "Spectre is the next gen cloud management platform"
var key = "$spectre*_*project*_2016_filler!" // 32 bytes!
var iv = "spectre_project1" // 16 bytes
var enc = "iYqWzi8jQcGQwjTHPn9624F7GFjGB2o9GF/P9uczQpmh3erPZFbeYCmCDs8r45ug30v5aIFpO6kHODzh3g4wnA=="

func TestEncryptdata(t *testing.T) {
	enc_data, err := EncryptData(key, iv, data)
	require.Nil(t, err)
	require.NotNil(t, enc_data)
	require.Equal(t, enc, enc_data)
}

func TestEncInvalidKey(t *testing.T) {
	enc_data, err := EncryptData("", iv, data)
	require.NotNil(t, err)
	require.EqualError(t, err, "Key cannot be empty")
	require.Equal(t, "", enc_data)

	enc_data, err = EncryptData("blah", iv, data)
	require.NotNil(t, err)
	require.EqualError(t, err, "Key should be 32 bytes in length.")
	require.Equal(t, "", enc_data)
}

func TestEncInvalidIV(t *testing.T) {
	enc_data, err := EncryptData(key, "", data)
	require.NotNil(t, err)
	require.EqualError(t, err, "IV cannot be empty")
	require.Equal(t, "", enc_data)

	enc_data, err = EncryptData(key, "blah", data)
	require.NotNil(t, err)
	require.EqualError(t, err, "IV should be 16 bytes in length.")
	require.Equal(t, "", enc_data)
}

func TestEncInvalidData(t *testing.T) {
	enc_data, err := EncryptData(key, iv, "")
	require.NotNil(t, err)
	require.EqualError(t, err, "Data cannot be empty")
	require.Equal(t, "", enc_data)
}

func TestDecryptdata(t *testing.T) {
	dec, err := DecryptData(key, iv, enc)
	require.Nil(t, err)
	require.NotNil(t, dec)
	require.Equal(t, data, dec)
}

func TestInvalidDecryptdata(t *testing.T) {
	dec, err := DecryptData(key, iv, "")
	require.NotNil(t, err)
	require.EqualError(t, err, "Data cannot be empty")
	require.Equal(t, "", dec)
}

func TestInvalidUnPad(t *testing.T){
	dec, err := DecryptData(key, iv, "blah")
	require.Nil(t, err)
	require.Equal(t, "", dec)
}

