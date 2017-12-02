package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/logging"
)

var log = logging.Log.Logger

// Validates common parameters
func checkParams(key, iv_string, data string) error {
	if len(key) == 0 {
		return errors.New("Key cannot be empty")
	}
	if len(key) != 32 {
		return errors.New("Key should be 32 bytes in length.")
	}
	if len(iv_string) == 0 {
		return errors.New("IV cannot be empty")
	}
	if len(iv_string) != aes.BlockSize {
		return fmt.Errorf("IV should be %d bytes in length.", aes.BlockSize)
	}
	if len(data) == 0 {
		return errors.New("Data cannot be empty")
	}
	return nil
}

// Encryptes given data using the given key and iv
func EncryptData(key, iv_string, data string) (string, error) {
	if err := checkParams(key, iv_string, data); err != nil {
		return "", err
	}

	str, err := pad_pkcs7_for_aes([]byte(data))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := []byte(iv_string)

	encVal := make([]byte, len(str))

	encInst := cipher.NewCBCEncrypter(block, iv)
	encInst.CryptBlocks(encVal, str)

	log.Debugf("Encrypted value: %v", encVal)

	return base64.StdEncoding.EncodeToString(encVal), nil
}

// Decryptes given data using the given key and iv
func DecryptData(key, iv_string, data string) (string, error) {

	if err := checkParams(key, iv_string, data); err != nil {
		return "", err
	}

	mainData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := []byte(iv_string)

	decrypted := make([]byte, len(mainData))

	defer func() {
		if r := recover(); r != nil {
			log.Debug("Recovered:", r)
		}
	}()

	decInst := cipher.NewCBCDecrypter(block, iv)
	decInst.CryptBlocks(decrypted, mainData)

	log.Debugf("Decrypt data: %s", decrypted)

	decrypted, err = remove_pad_pkcs7_for_aes(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// Pads the given data
func pad_pkcs7_for_aes(data_to_pad []byte) ([]byte, error) {
	if len(data_to_pad) == 0 {
		return nil, errors.New("Pad error: data is empty")
	}
	padding_length := 1
	for data_len := len(data_to_pad); ((data_len + padding_length) % aes.BlockSize) != 0; {
		padding_length = padding_length + 1
	}

	pad := bytes.Repeat([]byte{byte(padding_length)}, padding_length)
	padded_data := append(data_to_pad, pad...)
	return padded_data, nil
}

// Checks the padding on the given data and removes it, if valid
func remove_pad_pkcs7_for_aes(data_to_unpad []byte) ([]byte, error) {
	if len(data_to_unpad)%aes.BlockSize != 0 || len(data_to_unpad) == 0 {
		return nil, errors.New("Unpad error: data not valid")
	}
	padded_length := int(data_to_unpad[len(data_to_unpad)-1])
	if padded_length > aes.BlockSize || padded_length == 0 {
		return nil, errors.New("Unpad error: padding not valid")
	}
	padding := data_to_unpad[len(data_to_unpad)-padded_length:]
	for i := 0; i < padded_length; i++ {
		if padding[i] != byte(padded_length) {
			return nil, errors.New("Unpad error: padding not valid")
		}
	}
	unpadded_data := data_to_unpad[:len(data_to_unpad)-padded_length]
	return unpadded_data, nil
}
