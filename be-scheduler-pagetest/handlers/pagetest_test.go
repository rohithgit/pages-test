package handlers

import (
	log "github.com/Sirupsen/logrus"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"

	"testing" //see: https://github.com/stretchr/testify/assert
	"os"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"

	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/specnl/spectre-base-microservice/redis"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"

	"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
)

var ctx = context.TODO()
var url = "www.cisco.com"
var ctx1  = context.TODO()

func TestMain(m *testing.M) {
	//data.CreateTransactions()
	global.Options.Environment = "test"
	global.Options.WptServer = "http://spectre-webpagetest.infra-01.cisco-appanalytics.com"
	//global.Options.TopoUrl = "spectre-topocrudsvc.service.consul:31030"
	//global.Options.TopoPath = "/api/v1/topo/apps"
	global.Options.SchedulerEnabled = false
	global.Options.SchedulerInt = 1
	global.Options.SchedulerSleep = 2
	global.Options.Elastic = "http://localhost:9200"
	global.Options.RedisServers = "172.23.212.130:6379"

	mockDB := persist.NewTestResultDBMock().(*persist.TestResultDBMock)
	mockDB.MockInsertTest = func(test *models.ResultsStorage) (string, error) {
		return "1", nil
	}
	mockDB.MockGetTest = func(id string) (*models.ResultsStorage, error) {
		return nil, errors.New("not found")
	}
	ctx = context.WithValue(ctx, ctxDBKey, mockDB)
	cacheConn, _ :=  redis.CreateConn("","test")
	utils.SpectreLog.Debugln("Redis: %v", cacheConn)
	ctx = context.WithValue(ctx, ctxCacheKey, persist.NewRedisMock())
//<<<<<<< HEAD
//	hazelcastConn, _ :=  hazelcast.Connect("")
//	fmt.Printf("HAZELCAST: %v", hazelcastConn)
//	ctx = context.WithValue(ctx, ctxCacheKey, persist.NewHazelcastMock())

	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer persist.EndSession(session)
	// initialize the mongodb session
	persist.InitSession(session)
	time.Sleep(time.Second*10)
	//session := initMdbSession()
	ctx1 = getStorageContext()



	os.Exit(m.Run())
}

func getStorageContext() (context.Context)  {
	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		log.Errorf("Failed to init Test Result DB: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return nil
	}
	ctx1 := context.WithValue(context.TODO(), ctxDBKey, resultsDB)
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		log.Errorf("Failed to init Test Result ES: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return nil
	}
	ctx1 = context.WithValue(ctx1, CtxESKey, resultsES)
	cache, errC := helper.FetchRCacheFromContext(ctx1)

	if errC != nil {
		log.Errorf("Failed to init Hazelcast Cache: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return nil
	}
	ctx1 = context.WithValue(ctx1, ctxCacheKey, cache)
	return ctx1

}


func addMockDBFunctionality() {
	mockDB := persist.NewTestResultDBMock().(*persist.TestResultDBMock)
	mockDB.MockInsertTest = func(test *models.ResultsStorage) (string, error) {
		return test.Results, nil
	}
	mockDB.MockGetTest = func(id string) (*models.ResultsStorage, error) {
		return &models.ResultsStorage{Results:id}, nil
	}
	ctx = context.WithValue(ctx, ctxDBKey, mockDB)
}

// TestLocationHandler tests the TestLocationHandler
/*func TestLocationHandler(t *testing.T) {
	//t.SkipNow()

	global.Testing = true

	// create a router
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	// subrouter
	s.HandleFunc("/pagetest", PageTest).Name("pagetest")

	// test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// create path to call
	path := "/api/v1/pagetest"

	// send request
	res, err := http.Get(ts.URL + path)
	assert.NoError(t, err, "request errored")

	// read the response
	res2, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err, "response read errored")

	// test response
	assert.Contains(t, string(res2), "pagetest")
}*/

func TestPageLoadSummary(t *testing.T) {
	loadTime, err, status, statusText, lookupUrl := GetPageLoad(FakeTester{}, ctx, url, "")

	if err != nil {
		utils.SpectreLog.Debugln(" Error loading page: %s \n", err.Error())
	} else {
		// Note that this will require you add fmt to your list of imports.
		utils.SpectreLog.Debugln("")
		utils.SpectreLog.Debugln("Load Time %d \n", loadTime)
		utils.SpectreLog.Debugln("Look up Url %s \n", lookupUrl)
	}

	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln("Status: %d \n", status)
	utils.SpectreLog.Debugln("Status text: %s \n", statusText)

	//Test location

	loadTime, err, status, statusText, lookupUrl = GetPageLoad(FakeTester{}, ctx, url, "Test")

	if err != nil {
		utils.SpectreLog.Debugln(" Error loading page: %s \n", err.Error())
	} else {
		// Note that this will require you add fmt to your list of imports.
		utils.SpectreLog.Debugln("")
		utils.SpectreLog.Debugln("Load Time %d \n", loadTime)
		utils.SpectreLog.Debugln("Look up Url %s \n", lookupUrl)
	}

	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln("Status: %d \n", status)
	utils.SpectreLog.Debugln("Status text: %s \n", statusText)
}


func TestGetDetails(t *testing.T) {
	content, err, status := FakeTester{}.GetContent("http://128.107.18.61/jsonResult.php?test=160829_R4_1")
	pageLoad, err, status, statusText := GetPageLoadFromDetails(ctx, "http://128.107.18.61/jsonResult.php?test=160629_BY_4", content, status, true)
	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln("Status : %d ", status)
	utils.SpectreLog.Debugln("Status Text : %s ", statusText)
	utils.SpectreLog.Debugln("")
	if err != nil {
		utils.SpectreLog.Debugln("")
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln("PageLoadTime : %d ", pageLoad)
}

func TestUnmarshalNormalResponse2(t *testing.T) {

}

func TestGetLoadTimeFromCache(t *testing.T) {
	pageLoad, err, status, statusText, lookupUrl := GetPageLoadFromCache(FakeTester{}, ctx, "ui-portalserver.sjdev2.cloud.cisco.com", "", "")
	utils.SpectreLog.Debugln("Status : %d \n", status)
	utils.SpectreLog.Debugln("StatusText : %s \n", statusText)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}

	utils.SpectreLog.Debugln(" PageLoadTime : %d ", pageLoad)
}

func TestUpdateCustomerInDb(t *testing.T) {
	global.Options.Mongo = ""
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer persist.EndSession(session)
	// initialize the mongodb session
	persist.InitSession(session)
	_,  err := persist.NewTestResultDB()
	if(err != nil) {
		return
	}
	//ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)
	//TODO uncomment when audit trail done
	//updateCustomerInDb(ctx, "https://ui-portalserver.sjdev2.cloud.cisco.com", "http://128.107.18.61/jsonResult.php?test=160629_H7_E", "tenantId1","customerId1")

}

func TestGetTestResultsFromCache(t *testing.T) {
	results, err, statusCode, statusText, lookupUrl := GetTestResultsFromCache(FakeTester{}, ctx, url, "")
	utils.SpectreLog.Debugln("Test Results : %s \n", results)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)
	utils.SpectreLog.Debugln("Status code : %d \n", statusCode)
	utils.SpectreLog.Debugln("Status text : %s \n", statusText)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
}

func TestLookupTestResultsFromCache(t *testing.T) {
	results, err, statusCode, statusText, lookupUrl := LookupTestResultsFromCache(FakeTester{}, ctx, "http://128.107.18.61/jsonResult.php?test=160829_R4_1")
	utils.SpectreLog.Debugln("Test Results : %s \n", results)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)
	utils.SpectreLog.Debugln("Status code : %d \n", statusCode)
	utils.SpectreLog.Debugln("Status text : %s \n", statusText)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}

	addMockDBFunctionality()

	results, err, statusCode, statusText, lookupUrl = LookupTestResultsFromCache(FakeTester{}, ctx, "http://128.107.18.61/jsonResult.php?test=160629_H7_E")
	utils.SpectreLog.Debugln("Test Results : %s \n", results)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)
	utils.SpectreLog.Debugln("Status code : %d \n", statusCode)
	utils.SpectreLog.Debugln("Status text : %s \n", statusText)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
}

func TestLookupPageLoadFromCache(t *testing.T) {
	pageLoad, err, status, statusText, lookupUrl := LookupPageLoadFromCache(FakeTester{}, ctx, "", "http://128.107.18.61/jsonResult.php?test=160629_R0_A", "")

	utils.SpectreLog.Debugln("Status : %d \n", status)
	utils.SpectreLog.Debugln("StatusText : %s \n", statusText)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)
	utils.SpectreLog.Debugln("")
	if err != nil {
		utils.SpectreLog.Debugln("")
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln(" PageLoadTime : %d ", pageLoad)

	// Test unmarshal from DB

	addMockDBFunctionality()
	pageLoad, err, status, statusText, lookupUrl = LookupPageLoadFromCache(FakeTester{}, ctx, "", "http://128.107.18.61/jsonResult.php?test=160629_R0_A", "")

	utils.SpectreLog.Debugln("Status : %d \n", status)
	utils.SpectreLog.Debugln("StatusText : %s \n", statusText)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)
	utils.SpectreLog.Debugln("")
	if err != nil {
		utils.SpectreLog.Debugln("")
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln(" PageLoadTime : %d ", pageLoad)

}

func TestUnmarshalNormalResponse(t *testing.T) {
	//content, _, _ := GetContent("http://128.107.18.61/jsonResult.php?test=160825_K2_14W")
	var content []byte
	content = []byte {123, 34, 100, 97, 116, 97, 34, 58,
		123, 34, 105, 100, 34, 58, 34, 49, 54, 48, 56, 50, 53, 95,
		75, 50, 95, 49, 52, 87, 34, 44, 34, 117, 114, 108, 34, 58, 34,
		104, 116, 116, 112, 58, 92, 47, 92, 47, 119, 119, 119, 46, 109, 115,
		110, 46, 99, 111, 109, 34, 44, 34, 115, 117, 109, 109, 97, 114, 121, 34,
		58, 34, 104, 116, 116, 112, 58, 92, 47, 92, 47, 49, 50, 56, 46, 49, 48, 55}

	err, _ := UnmarshalNormalResponse(content)

	if(err == nil) {
		t.Error("Expected unexpected end of JSON input, got ", err)
	}
}


func TestGetContent(t *testing.T) {
	content, err, _ := FakeTester{}.GetContent("http://128.107.18.61/runtest.php?f=json&url=www.msn.com")

	if err != nil  {
		t.Error("Expected no error, got: ", err)
	} else if content == nil {
		t.Error("Expected content for www.msn.com run test, got nil instead")
	}
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
		utils.SpectreLog.Debugln("Duplicate app name: ", appName)
		//utils.SpectreLog.Debugln("app ID is ", appID)
	} else if  err != nil  {
		t.Error("Expected no error, got: ", err)
	} else {
		utils.SpectreLog.Debugln("Test status: ", status)
		utils.SpectreLog.Debugln("Got app ID:", appID)
	}
}
