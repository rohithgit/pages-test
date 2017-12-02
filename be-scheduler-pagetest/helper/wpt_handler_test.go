package helper

import (
	"net/http"
	"testing"

	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
)

type FakeTester struct{
	AppId string
	AppErr error
	AppStatus int
}
// This function fetch the content of a URL will return it as an
// array of bytes if retrieved successfully.
func (ft FakeTester) GetContent(url string) ([]byte, error, int) {
	// Build the request
	if (url == "http://128.107.18.61/runtest.php?url=www.cisco.com&ignoreSSL=1&f=json&fvonly=1&location=Test San Jose") {
		const strContent = "{\"statusCode\":200,\"statusText\":\"Ok\",\"data\":{\"testId\":\"160829_R4_1\",\"ownerKey\":\"c4fb40a0275d07dd5d9c4cbba4f8b16491a5cd56\",\"jsonUrl\":\"http://128.107.18.61/jsonResult.php?test=160829_R4_1\",\"xmlUrl\":\"http://128.107.18.61/xmlResult.php?test=160829_R4_1\",\"userUrl\":\"http://128.107.18.61/results.php?test=160829_R4_1\",\"summaryCSV\":\"http://128.107.18.61/csv.php?test=160829_R4_1\",\"detailCSV\":\"http://128.107.18.61/csv.php?test=160829_R4_1&amp;requests=1\"}}"
		content := make([]byte, len(strContent))
		copy(content[:], strContent)
		return content, nil, http.StatusOK
	}
	//else if (url == "http://128.107.18.61/jsonResult.php?test=160829_R4_1") {
	//		return content, nil, http.StatusOK
	//	}

	return nil, errors.New("Response not set for url: " + url), http.StatusNotFound
}

func (ft FakeTester) GetAppId( appname string) (string, error, int) {
	return ft.AppId, ft.AppErr, ft.AppStatus
}

func TestGetAppId(t *testing.T) {
	f := FakeTester{
		AppId: "57c73e4e636f7700017d744f",
		AppErr: nil,
		AppStatus: http.StatusCreated,
	}

	appName := "test40"
	appID, err, status := f.GetAppId(appName)

	if (status == 409)  {
		utils.SpectreLog.Println("Duplicate app name: ", appName)
		//utils.SpectreLog.Println("app ID is ", appID)
	} else if  err != nil  {
		t.Error("Expected no error, got: ", err)
	} else {
		utils.SpectreLog.Println("Test status: ", status)
		utils.SpectreLog.Println("Got app ID:", appID)
	}
}

func TestDupAppId(t *testing.T) {
	f := FakeTester{
		AppId: "57c73e4e636f7700017d744f",
		AppErr: errors.New("Duplicate app name"),
		AppStatus: http.StatusConflict,
	}

	appName := "test40"
	appID, err, status := f.GetAppId(appName)

	if (status == 409)  {
		utils.SpectreLog.Println("Duplicate app name: ", appName)
		//utils.SpectreLog.Println("app ID is ", appID)
	} else if  err != nil  {
		t.Error("Expected no error, got: ", err)
	} else {
		utils.SpectreLog.Println("Test status: ", status)
		utils.SpectreLog.Println("Got app ID:", appID)
	}
}
