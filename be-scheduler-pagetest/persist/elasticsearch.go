package persist

import (
	elastic "gopkg.in/olivere/elastic.v3"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"errors"
	"encoding/json"
	"reflect"
	"strings"
	"strconv"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/service_discovery"
	"fmt"
)

type ITestResultES interface {
	FindTests(query map[string]interface{}) ([]models.TestResultsES, error)
	GetUrl(url, location, timeDuration string) ([]models.TestResultsES, error)
	GetLatestForUrl(url, location, timeDuration string) (*models.TestResultsES, error)
	GetTest(testurl string) (*models.TestResultsES, error)
	GetAllTestsForUrl(url string) ([]models.TestResultsES, error)
	GetAllTestResults(userId string) ([]models.TestResultsES, error)
	HasTest(test models.TestResultsES) bool
	HasTestId(id string) bool
	InsertTest(results models.TestResultsES) (string, error)
	RemoveTest(id string) error
	UpdateTest(id string, update models.TestResultsES) error
	GetSuccessfulRuns(url, location, timeDuration string) (int64, error)
	GetSlowestRun(url, location, timeDuration string) (*models.TestResultsES, error)
}

type TestResultES struct {
	testResults *elastic.Client
}

func NewTestResultES() (ITestResultES, error) {
	var err error
	resDB := new(TestResultES)
	if resDB.testResults, err = initES(); err != nil {
		return nil, err
	}
	return resDB, nil
}


func initES() (*elastic.Client, error) {
	var (
		es *elastic.Client
		err error
	)

	utils.SpectreLog.Debugln("Connecting to ES")
	var esUrl string
	useAuth := false
	if global.Options.Environment == "dev" || !global.Options.EnableSvcDiscovery {
		esUrl = global.Options.Elastic
		useAuth = true
	} else {
		// use service discovery
		db_addr := service_discovery.NewServiceAddress(nil)
		err1 := db_addr.GetServiceAddress(global.Options.Elastic)
		if err1 == nil {
			esUrl = fmt.Sprintf("http://%s:%d", db_addr.Name, 9200)
		} else {
			esUrl = global.Options.Elastic
		}
		es, err = elastic.NewSimpleClient(elastic.SetURL(global.Options.Elastic))
	}
	utils.SpectreLog.Debugln("Returned from ES")

	if useAuth {
		utils.SpectreLog.Debugln("Setting basic auth...")
		es, err = elastic.NewSimpleClient(elastic.SetURL(esUrl),
			elastic.SetBasicAuth(global.Options.EsUser, global.Options.EsPassword))
	} else {
		es, err = elastic.NewSimpleClient(elastic.SetURL(esUrl))
	}

	if err != nil {
		utils.SpectreLog.Errorf("Elastic Search init Error: " + err.Error())
		return nil, err
	}

	version, err := es.ElasticsearchVersion(esUrl)

	if err != nil {
		// Handle error
		utils.SpectreLog.Errorf("ES Client ping errored: %s", err.Error())
		return nil, err
	}

	utils.SpectreLog.Debugf("Elasticsearch version %s\n", version)

	err = setupIndex(es)
	if err != nil {
		utils.SpectreLog.Errorf("Elastic Search init Error: " + err.Error())
		return nil, err
	}
	return es, nil
}

func setupIndex(client *elastic.Client) error {
	if ok, err := client.IndexExists("pagetestresults").Do(); !ok && err == nil {
		client.CreateIndex("pagetestresults").Do()
		_, err := client.PutMapping().Index(constants.ES_RESULTS_INDEX).
		Type(constants.ES_RESULTS_TYPE).BodyJson(
			map[string]interface{}{"properties": map[string]interface{}{
				"testurl": map[string]string{"type":"string", "index": "not_analyzed"},
				"pagetesturl": map[string]string{"type":"string", "index": "not_analyzed"},
			}},
		).Do()
		if err != nil {
			utils.SpectreLog.Error("Could not create index mapping.")
			return err
		}
	}
	return nil
}

func (es *TestResultES) FindTests(query map[string]interface{}) ([]models.TestResultsES, error) {
	stringQuery := "{\"query\" :{ \"bool\": {\"must\" : ["
	for key, val := range query { //TODO add more capabilities if needed
		stringQuery = stringQuery + "{\"term\" :{ \"" + key + "\" : \"" + val.(string) + "\" }},"
	}
	stringQuery = strings.TrimSuffix(stringQuery,",")
	stringQuery = stringQuery + "]}}}"
	rawQuery := elastic.NewRawStringQuery(stringQuery)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Query(rawQuery).
	Sort("runtime", true).
	From(0).Size(50).
	Do()
	if err != nil {
		// Handle error
		return nil, err
	}
	var results []models.TestResultsES
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var result models.TestResultsES
			err := json.Unmarshal(*hit.Source, &result)
			if err != nil {
				return results, err
			}
			result.ID = hit.Id
			results = append(results, result)
		}
	}
	return results, nil
}

func (es *TestResultES) GetUrl(url, location, timeDuration string) ([]models.TestResultsES, error) {
	//termQuery := elastic.NewTermQuery("pagetesturl", url)
	//rangeQuery := elastic.NewRangeQuery("runtime").From(timeEnd.Add(-timeDuration)).To(timeEnd)
	timeRange, rangeError := fixTimeDurationForES(timeDuration)
	if rangeError != nil {
		return nil, errors.New("Invalid time range.")
	}
	location = strings.ToLower(location)
	queryStr := "{\"query\" :{ \"bool\": {\"must\" : [" +
	"{\"term\" :{ \"pagetesturl\" : \"" + url + "\" }}"
	if location != "" {
		queryStr += ",{\"term\" :{ \"location\" : \"" + location + "\" }}"
	}
	if timeRange != "" {
		queryStr += ",{\"range\" :{\"runtime\":{\"gte\": \"now-" + timeRange + "\", \"lte\": \"now\" }}}"
	}
	queryStr += "]}}}"
	rawQuery := elastic.NewRawStringQuery(queryStr)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).// search in index "twitter"
	Query(rawQuery).  // specify the query
	Sort("runtime", true). // sort by "user" field, ascending
	From(0).Size(10).   // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var results []models.TestResultsES
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var result models.TestResultsES
			err := json.Unmarshal(*hit.Source, &result)
			if err != nil {
				return results, err
			}
			result.ID = hit.Id
			results = append(results, result)
		}
	}
	if results != nil && len(results) > 0 {
		return results, nil
	} else {
		return make([]models.TestResultsES, 0), nil
	}

}

func (es *TestResultES) GetLatestForUrl(url, location, timeDuration string) (*models.TestResultsES, error) {
	timeRange, rangeError := fixTimeDurationForES(timeDuration)
	if rangeError != nil {
		return nil, errors.New("Invalid time range.")
	}
	location = strings.ToLower(location)
	queryStr := "{\"query\" :{ \"bool\": {\"must\" : [" +
	"{\"term\" :{ \"pagetesturl\" : \"" + url + "\" }}"
	if location != "" {
		queryStr += ",{\"term\" :{ \"location\" : \"" + location + "\" }}"
	}
	if timeRange != "" {
		queryStr += ",{\"range\" :{\"runtime\":{\"gte\": \"now-" + timeRange + "\", \"lte\": \"now\" }}}"
	}
	queryStr += "]}}}"
	rawQuery := elastic.NewRawStringQuery(queryStr)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).// search in index "twitter"
	Query(rawQuery).  // specify the query
	Sort("runtime", false).
	From(0).Size(1).   // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var id string
	if searchResult.TotalHits() > 0 {
		id = searchResult.Hits.Hits[0].Id
	}
	var estype models.TestResultsES
	for _, item := range searchResult.Each(reflect.TypeOf(estype)) {
		if t, ok := item.(models.TestResultsES); ok {
			t.ID = id
			return &t, nil
		}
	}
	return nil, errors.New("No results found.")

}

func (es *TestResultES) GetTest(testurl string) (*models.TestResultsES, error) {
	termQuery := elastic.NewTermQuery("testurl", testurl)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).// search in index "twitter"
	Query(termQuery).  // specify the query
	From(0).Size(1).   // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var id string
	if searchResult.TotalHits() > 0 {
		id = searchResult.Hits.Hits[0].Id
	}
	var estype models.TestResultsES
	for _, item := range searchResult.Each(reflect.TypeOf(estype)) {
		if t, ok := item.(models.TestResultsES); ok {
			t.ID = id
			return &t, nil
		}
	}
	return nil, errors.New("No results found.")

}

func (es *TestResultES) GetAllTestsForUrl(url string) ([]models.TestResultsES, error) {
	termQuery := elastic.NewTermQuery("pagetesturl", url)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).   // search in index "twitter"
	Type(constants.ES_RESULTS_TYPE).
	Query(termQuery).  // specify the query
	Sort("runtime", true). // sort by "user" field, ascending
	From(0).Size(50).   // take documents 0-49
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var results []models.TestResultsES
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var result models.TestResultsES
			err := json.Unmarshal(*hit.Source, &result)
			if err != nil {
				return results, err
			}
			result.ID = hit.Id
			results = append(results, result)
		}
	}
	return results, nil
}

func (es *TestResultES) GetAllTestResults(userId string) ([]models.TestResultsES, error) {
	termQuery := elastic.NewTermQuery("userid", userId)
	searchResult, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).   // search in index "twitter"
	Type(constants.ES_RESULTS_TYPE).
	Query(termQuery).  // specify the query
	Sort("runtime", true). // sort by "user" field, ascending
	From(0).Size(10).   // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var results []models.TestResultsES
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var result models.TestResultsES
			err := json.Unmarshal(*hit.Source, &result)
			if err != nil {
				return results, err
			}
			result.ID = hit.Id
			results = append(results, result)
		}
	}
	return results, nil
}

func (es *TestResultES) HasTest(test models.TestResultsES) bool {
	termQuery := elastic.NewMatchQuery("testurl", test.TestUrl)
	out, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Query(termQuery).
	Do()
	if err != nil {
		return false
	}
	return out.TotalHits() > 0
}

func (es *TestResultES) HasTestId(id string) bool {
	out, err := es.testResults.Exists().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).
	Id(id).
	Do()
	if err != nil {
		return false
	}
	return out
}

func (es *TestResultES) InsertTest(results models.TestResultsES) (string, error) {
	out, err := es.testResults.Index().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).
	BodyJson(results).
	Refresh(true).
	Do()
	if err != nil {
		// Handle error
		return "", err
	}
	return out.Id, nil
}

func (es *TestResultES) RemoveTest(id string) error {
	res, err := es.testResults.Delete().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).
	Id(id).
	Do()
	if err != nil {
		// Handle error
		return err
	}
	if res.Found {
		return nil
	}
	return errors.New("ID " + id + " not found.")
}

func (es *TestResultES) UpdateTest(id string, update models.TestResultsES) error {
	bytes, err := json.Marshal(update)
	if err != nil {
		return err
	}
	resultsMap := make(map[string]interface{})
	json.Unmarshal(bytes, &resultsMap)
	_, err = es.testResults.
	Update().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).
	Id(id).
	Doc(resultsMap).
	Do()
	if err != nil {
		// Handle error
		return err
	}
	return nil
}

func (es *TestResultES) GetSuccessfulRuns(url, location, timeDuration string) (int64, error) {
	timeRange, rangeError := fixTimeDurationForES(timeDuration)
	if rangeError != nil {
		return 0, errors.New("Invalid time range.")
	}
	location = strings.ToLower(location)
	queryStr := "{\"query\" :{ \"bool\": {\"must\" : [" +
	"{\"term\" :{ \"pagetesturl\" : \"" + url + "\" }}"
	if location != "" {
		queryStr += ",{\"term\" :{ \"location\" : \"" + location + "\" }}"
	}
	if timeRange != "" {
		queryStr += ",{\"range\" :{\"runtime\":{\"gte\": \"now-" + timeRange + "\", \"lte\": \"now\" }}}"
	}
	queryStr += "]}}}"
	rawQuery := elastic.NewRawStringQuery(queryStr)
	count, err := es.testResults.Count().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).// search in index "twitter"
	Query(rawQuery).  // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return 0, err
	}
	return count, nil
}

func (es *TestResultES) GetSlowestRun(url, location, timeDuration string) (*models.TestResultsES, error) {
	timeRange, rangeError := fixTimeDurationForES(timeDuration)
	if rangeError != nil {
		return nil, errors.New("Invalid time range.")
	}
	location = strings.ToLower(location)
	queryStr := "{\"query\" :{ \"bool\": {\"must\" : [" +
	"{\"term\" :{ \"pagetesturl\" : \"" + url + "\" }}"
	if location != "" {
		queryStr += ",{\"term\" :{ \"location\" : \"" + location + "\" }}"
	}
	if timeRange != "" {
		queryStr += ",{\"range\" :{\"runtime\":{\"gte\": \"now-" + timeRange + "\", \"lte\": \"now\" }}}"
	}
	queryStr += "]}}}"
	rawQuery := elastic.NewRawStringQuery(queryStr)
	results, err := es.testResults.Search().
	Index(constants.ES_RESULTS_INDEX).
	Type(constants.ES_RESULTS_TYPE).// search in index "twitter"
	Query(rawQuery).  // take documents 0-9
	Sort("loadtime", false). // sort by "user" field, ascending
	From(0).Size(1).   // take documents 0-9
	Do()                // execute
	if err != nil {
		// Handle error
		return nil, err
	}
	var estype models.TestResultsES
	for _, item := range results.Each(reflect.TypeOf(estype)) {
		if t, ok := item.(models.TestResultsES); ok {
			return &t, nil
		}
	}
	return nil, errors.New("No runs found.")
}


func fixTimeDurationForES(timeDuration string) (string, error) {
	switch {
	case timeDuration == "":
		return "", nil
	case strings.Contains(timeDuration, "min"):
		val, err := strconv.Atoi(strings.Replace(timeDuration, "min", "", -1))
		if err == nil {
			return strconv.Itoa(val) + "m", nil
		}
	case strings.Contains(timeDuration, "h"):
		val, err := strconv.Atoi(strings.Replace(timeDuration, "h", "", -1))
		if err == nil {
			return strconv.Itoa(val) + "h", nil
		}
	case strings.Contains(timeDuration, "d"):
		val, err := strconv.Atoi(strings.Replace(timeDuration, "d", "", -1))
		if err == nil {
			return strconv.Itoa(val) + "d", nil
		}
	case strings.Contains(timeDuration, "m"):
		val, err := strconv.Atoi(strings.Replace(timeDuration, "m", "", -1))
		if err == nil {
			return strconv.Itoa(val) + "M", nil
		}
	}
	utils.SpectreLog.Errorf("Invalid time range %s", timeDuration)
	return "6h", errors.New("Invalid time range")
}

type TestResultESMock struct {
	MockFindTests func(query map[string]interface{}) ([]models.TestResultsES, error)
	MockGetUrl func(url, location, timeDuration string) ([]models.TestResultsES, error)
	MockGetLatestForUrl func(url, location, timeDuration string) (*models.TestResultsES, error)
	MockGetTest func(testurl string) (*models.TestResultsES, error)
	MockGetAllTestsForUrl func(url string) ([]models.TestResultsES, error)
	MockGetAllTestResults func(userId string) ([]models.TestResultsES, error)
	MockHasTest func(test models.TestResultsES) bool
	MockHasTestId func(id string) bool
	MockInsertTest func(results models.TestResultsES) (string, error)
	MockRemoveTest func(id string) error
	MockUpdateTest func(id string, update models.TestResultsES) error
	MockGetSuccessfulRuns func (url, location, timeDuration string) (int64, error)
	MockGetSlowestRun func(url, location, timeDuration string) (*models.TestResultsES, error)

}

func NewTestResultESMock() ITestResultES {
	return &TestResultESMock{}
}

func (es *TestResultESMock) FindTests(query map[string]interface{}) ([]models.TestResultsES, error){
	return es.MockFindTests(query)
}

func (es *TestResultESMock) GetUrl(url, location, timeDuration string) ([]models.TestResultsES, error) {
	return es.MockGetUrl(url, location, timeDuration)
}

func (es *TestResultESMock) GetLatestForUrl(url, location, timeDuration string) (*models.TestResultsES, error) {
	return es.MockGetLatestForUrl(url, location, timeDuration)
}

func (es *TestResultESMock) GetTest(testurl string) (*models.TestResultsES, error) {
	return es.MockGetTest(testurl)
}

func (es *TestResultESMock) GetAllTestsForUrl(url string) ([]models.TestResultsES, error) {
	return es.MockGetAllTestsForUrl(url)
}

func (es *TestResultESMock) GetAllTestResults(userId string) ([]models.TestResultsES, error) {
	return es.MockGetAllTestResults(userId)
}

func (es *TestResultESMock) HasTest(test models.TestResultsES) bool {
	return es.MockHasTest(test)
}

func (es *TestResultESMock) HasTestId(id string) bool {
	return es.MockHasTestId(id)
}

func (es *TestResultESMock) InsertTest(results models.TestResultsES) (string, error) {
	return es.MockInsertTest(results)
}

func (es *TestResultESMock) RemoveTest(id string) error {
	return es.MockRemoveTest(id)
}

func (es *TestResultESMock) UpdateTest(id string, update models.TestResultsES) error {
	return es.MockUpdateTest(id, update)
}

func (es *TestResultESMock) GetSuccessfulRuns(url, location, timeDuration string) (int64, error) {
	return 	es.MockGetSuccessfulRuns(url, location, timeDuration)
}

func (es *TestResultESMock) GetSlowestRun(url, location, timeDuration string) (*models.TestResultsES, error) {
	return 	es.MockGetSlowestRun(url, location, timeDuration)

}
