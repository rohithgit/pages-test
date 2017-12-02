// Package utils contains helper utility functions used by other packages in the repo
package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"strconv"
	"strings"
	"unsafe"
	"fmt"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
)

// PrintFields prints type/struct fields and values
func PrintFields(name string, i interface{}) {
	s := reflect.ValueOf(i).Elem()
	typeOfT := s.Type()
	SpectreLog.Debugln("%s fields\n", name)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		SpectreLog.Debugln("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

// SprintFields returns a string of fields and values
func SprintFields(name string, i interface{}) string {
	s := reflect.ValueOf(i).Elem()
	typeOfT := s.Type()
	rtn := fmt.Sprintf("%s fields\n", name)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		rtn += fmt.Sprintf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	return rtn
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func FetchGet(url string, auth models.AuthStruct, headers []models.RequestHeader ,alwaysTrust bool) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if auth.Type == "Token" {
		if auth.AuthHeader != "" {
			req.Header.Set(auth.AuthHeader, auth.AuthToken)
		}
	} else if auth.Type == "Basic" {
		req.SetBasicAuth(auth.User, auth.Password)
	} else if auth.Type == "access_Token" {
		req.Header.Set("access_token", auth.AccessToken)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("DNT", "1")
		SpectreLog.Println("access token", auth.AccessToken)
	} else if auth.Type == "X-Cloupia-Request-Key" {
		req.Header.Set("X-Cloupia-Request-Key", auth.AccessToken)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("DNT", "1")
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}
	cli := &http.Client{}
	if alwaysTrust {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cli = &http.Client{Transport: tr}
	}
	SpectreLog.Println("url = ", req.URL)
	resp, err := cli.Do(req)
	if err != nil {
		SpectreLog.Println(err)
	}
	if resp != nil {
		SpectreLog.Println(resp)
	}
	if err == nil {
		content, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		//defer req.Body.Close()
		contentType := resp.Header.Get("Content-Type")
		//SpectreLog.Println(BytesToString(content))
		//SpectreLog.Println(contentType)
		return content, contentType, nil
	} else {
		return nil, "", err
	}
}

func FetchPost(url string, body string, auth models.AuthStruct) ([]byte, string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if auth.Type == "Token" {
		if auth.AuthHeader != "" {
			req.Header.Set(auth.AuthHeader, auth.AuthToken)
		}
	} else if auth.Type == "Basic" {
		req.SetBasicAuth(auth.User, auth.Password)
	}
	/*req.Header.Set( "Content-Type", "application/x-www-form-urlencoded")*/
	req.Header.Set("Content-Type", auth.ContentType)
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		SpectreLog.Println(err)
	}
	if resp != nil {
		SpectreLog.Println(resp)
	}
	if err == nil {

		defer req.Body.Close()
		defer responseClose(resp)
		content, _ := ioutil.ReadAll(resp.Body)
		contentType := resp.Header.Get("Content-Type")
		return content, contentType, nil
	} else {
		return nil, "", err
	}
}

func responseClose(w *http.Response) {
	if w != nil {
		w.Body.Close()
	}
}
func PrintError(w http.ResponseWriter, errorCode int, errorMessage string) {
	err1 := models.ErrorCode{}
	err1.Error = errorMessage
	err1.Code = strconv.Itoa(errorCode)
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(err1); err != nil {
		panic(err)
	}
}

func InArray(val string, array []string) bool {
	for _, str := range array {
		if strings.ToLower(val) == strings.ToLower(str) {
			return true
		}
	}
	return false
}

func ConvertToint( val string, defaultVal int) int{
	returnVal := defaultVal
	if s, err := strconv.Atoi(val); err == nil {
		returnVal = s
	}
	fmt.Println("%T, %v", returnVal, returnVal)
	return returnVal

}

func ConvertIntToString( val int) string{
	return strconv.Itoa(val)

}

func ApplicationUrlCountCacheKey( applicationId string ) string {
	return constants.CACHE_PREFIX + applicationId + constants.URL_COUNT
}