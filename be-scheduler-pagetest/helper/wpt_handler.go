package helper

import (
	"net/http"
	"io/ioutil"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
)

type Tester interface {
	GetContent(string) ([]byte, error, int)
	//GetAppId(string) (string, error, int)
}

type WebPageTester struct{}

// This function fetch the content of a URL will return it as an
// array of bytes if retrieved successfully.
func (wpt WebPageTester) GetContent(url string) ([]byte, error, int) {
	// Build the request

	resp, err := http.Get(url)
	utils.SpectreLog.Println(resp)
	utils.SpectreLog.Println(err)
	statusCode := 500
	if( resp != nil) {
		statusCode = resp.StatusCode
	}
	if (err == nil) && (statusCode == 200) {
		// Defer the closing of the body
		defer resp.Body.Close()
		// Read the content into a byte array
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err, statusCode
		}

		// At this point we're done - simply return the bytes
		return body, nil, statusCode
	} else {
		return nil, err, statusCode
	}
}

//func (wpt WebPageTester) GetAppId( appName string) (string, error, int) {
//	trusted := true //TODO get from config
//
//	if (global.Options.TopoUrl == "") || global.Options.TopoPath == "" {
//		utils.SpectreLog.Errorf("Topo service url not configured properly: %s", global.Options.TopoUrl)
//		return "", errors.New("Topo Service not configured"), http.StatusInternalServerError
//	}
//	url := "https://" + global.Options.TopoUrl+global.Options.TopoPath
//	headersMap := map[string]string{
//		"ContentType": "application/json",
//		"x-csco-tenantid": "tenant1", //TODO: get this from RBAC
//	}
//	client := &http.Client{}
//	if( trusted){
//		tr := &http.Transport{
//			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//		}
//		client = &http.Client{Transport: tr}
//	}
//	data := "{\"appName\":\""+ appName +"\"}"
//	req, err := http.NewRequest("POST", url, strings.NewReader(data))
//	utils.SpectreLog.Println( "url = " , url, data)
//	for key, value := range headersMap {
//		utils.SpectreLog.Println("Key:", key, "Value:", value)
//		req.Header.Add( key, value)
//	}
//
//	var app1 app_struct
//	resp, err := client.Do(req)
//	if (err == nil) && (resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 409 ) {
//		// Defer the closing of the body
//		defer resp.Body.Close()
//		// Read the content into a byte array
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//
//			return "", err, resp.StatusCode
//		}
//		utils.SpectreLog.Println( "resp = ", body)
//		utils.SpectreLog.Println("status = ", resp.StatusCode)
//		if(resp.StatusCode == 409) {
//			return "", errors.New("Duplicate App name" + appName), resp.StatusCode
//		}
//		// At this point we're done - simply return the bytes
//		err1 := json.Unmarshal(body, &app1)
//		if (err1 == nil) {
//			return app1.AppId, nil, resp.StatusCode
//		} else {
//			utils.SpectreLog.Debugf("Error unmarshalling content into app struct for app name %s", appName)
//			utils.SpectreLog.Errorf("Cannot retrieve appId for app ->%s", appName)
//			return  "", errors.New("Cannot retrieve app Id for app name: " + appName), http.StatusInternalServerError
//		}
//	} else {
//		utils.SpectreLog.Println( "resp = ", nil, err)
//		return "", err, resp.StatusCode
//	}
//}

type app_struct struct {
	AppName string `json:"appName,omitempty"`
	AppId string `json:"appId,omitempty"`
}
