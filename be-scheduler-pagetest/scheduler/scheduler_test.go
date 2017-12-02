package scheduler

import (
	"testing"
	"net/http"
	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"os"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"

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

func TestMain(m *testing.M) {
	//data.CreateTransactions()
	global.Options.Environment = "test"
	global.Options.WptServer = "128.107.18.61"
	global.Options.TopoUrl = "spectre-topocrudsvc.service.consul:31030"
	global.Options.TopoPath = "/api/v1/topo/apps"
	global.Options.SchedulerEnabled = true
	global.Options.SchedulerInt = 1
	global.Options.SchedulerSleep = 2
	global.Options.Elastic = "127.0.0.1:9200"


	os.Exit(m.Run())
}
func TestSubmitBatch(t *testing.T) {
	//f := FakeTester{
	//	AppId: "57c73e4e636f7700017d744f",
	//	AppErr: nil,
	//	AppStatus: http.StatusCreated,
	//}
	urls := []string{"www.cisco.com", "www.msn.com"}
	testIds, err := SubmitBatch(helper.WebPageTester{}, urls, "", "")

	if(err != nil ) {
		utils.SpectreLog.Println(err)
	}
	for _, testId := range testIds {
		utils.SpectreLog.Println("URL: ", testId)
	}
}

func TestGetBatchResults(t *testing.T) {
	//urls := []string{"www.cisco.com", "www.msn.com"}
	//testIds, err := SubmitBatch(urls, "", "")

	//f := FakeTester{
	//	AppId: "57c73e4e636f7700017d744f",
	//	AppErr: nil,
	//	AppStatus: http.StatusCreated,
	//}
	var testIds map[string]string
	testIds = make(map[string]string)
	testIds["160722_GK_J"] = "www.cisco.com"
	testIds["160722_M9_K"] = "www.msn.com"
	testresults, err := GetBatchResults(helper.WebPageTester{}, testIds, "")

	if(err != nil ) {
		utils.SpectreLog.Println(err)
	}
	for testId, results := range testresults {
		utils.SpectreLog.Println("--------------------")
		utils.SpectreLog.Println("URL: ", testId)
		utils.SpectreLog.Println("Results: ", results)
		utils.SpectreLog.Println()
		utils.SpectreLog.Println("--------------------")
	}
}

func TestScheduler(t *testing.T) {
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer persist.EndSession(session)
	// initialize the mongodb session
	persist.InitSession(session)
	StartScheduler()
}
