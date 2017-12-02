package handlers

import (

	"fmt"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"golang.org/x/net/context"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"errors"
	"time"
	"strings"
	neturl "net/url"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

type dbKey string
type cacheKey string
type esKey string

var (
	ctxDBKey dbKey = "resultsDB"
	CtxESKey esKey = "resultsES"
	ctxCacheKey cacheKey = "resultsCache"
)

func FetchDBFromContext(ctx context.Context) (persist.ITestResultDB, bool){
	if(ctx == nil) {
		utils.SpectreLog.Debugln("Context is nil")
		return nil, false
	}
	utils.SpectreLog.Debugln("Context is not null")
	resultsDB, ok := ctx.Value(ctxDBKey).(persist.ITestResultDB)
	utils.SpectreLog.Debugln("retrieved resultsDb ", ok)
	return resultsDB, ok
}

func FetchESFromContext(ctx context.Context) (persist.ITestResultES, bool){
	if(ctx == nil) {
		return nil, false
	}
	resultsDB, ok := ctx.Value(CtxESKey).(persist.ITestResultES)
	return resultsDB, ok
}

func Hello(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w,  "Ack Hello")
}

func PageTest(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debugln("Entering PageTestHandler()")

	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_WRITE) {
		return
	}

	utils.SpectreLog.Debugln("Token is Valid")


	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		utils.SpectreLog.Errorln("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorln("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)
	_, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorln("Failed to init redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	rawurl := r.URL.Query().Get("url")
	utils.SpectreLog.Debugln("Test URL %s \n", rawurl)
	if rawurl == "" {
		utils.PrintError(w, http.StatusBadRequest, "Invalid value in `value` parameter.")
		return
	}

	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}

	testResults, err1, statusCode, statusText, lookupUrl := GetTestResultsFromCache(helper.WebPageTester{}, ctx, url, location)
	if(err1 != nil) {
		utils.PrintError(w, http.StatusBadRequest, err1.Error())
		return
	}

	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)
	utils.SpectreLog.Debugln("Status: %s\n ", statusCode)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)

	results := new(models.Result)
	results.LookupUrl = lookupUrl
	results.StatusText = statusText
	results.TestStatusCode = statusCode
	results.TestResults = testResults

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LookupResults(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering LookupResultsHandler()")
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}

	utils.SpectreLog.Debugln("Token is Valid")

	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)
	_, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init Redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	url := r.URL.Query().Get("url")
	utils.SpectreLog.Debugln("Test URL ", url)

	testResults, err, status, statusText, lookupUrl := LookupTestResultsFromCache(helper.WebPageTester{}, ctx, url)
	if(err != nil) {
		utils.PrintError(w, status, err.Error())
		return
	}

	utils.SpectreLog.Debugln("Test Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.TestResults = testResults

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}

func GetTestResultsFromCache(wpt helper.Tester, ctx context.Context, url, location string) (string, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		results, err, status, statusText, lookupUrl := LookupTestResultsFromCache(wpt, ctx, lookupUrlFromCache)
		return results, err, status, statusText, lookupUrl
	}
	utils.SpectreLog.Debugln("Url is ", url)

	results, err, statusCode, statusText, lookupUrl := GetTestResults(wpt, ctx, url, location)
	if err != nil  {
		return results, err, statusCode, statusText, lookupUrl
	}
	cache.SetStringWithExpiration(constants.CACHE_PREFIX + url + "::" + location, lookupUrl, constants.PAGELOOKUP_EXPIRATION)
	if results != "" && statusCode == 200{
		cache.SetStringWithExpiration(constants.CACHE_PREFIX + lookupUrl, results, constants.TESTRESULTS_EXPIRATION)
	}
	return results, err, statusCode, statusText, lookupUrl
}

func LookupTestResultsFromCache(wpt helper.Tester, ctx context.Context, url string) (string, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	if results, _ := cache.GetString(constants.CACHE_PREFIX + url); results == ""  {

		utils.SpectreLog.Debugln("Lookup Url is ", url)
		results, err, status, statusText := GetResultsFromDetails(wpt, ctx, url)
		if results != "" && status == 200 {
			//cache.SetWithExpiration(constants.CACHE_PREFIX + url, results, constants.TESTRESULTS_EXPIRATION)
		}
		return results, err, status, statusText, url
	} else {
		return results, nil, http.StatusOK, "Cached Response", url
	}
}

func buildTestURL(url string, location string) (string) {

	pageTestUrl := global.Options.WptServer + "/runtest.php"
	resultUrl := pageTestUrl + "?" + "url=" + url + "&ignoreSSL=1&f=json&fvonly=1&location=" + location
	utils.SpectreLog.Debugln("Result URL: ", resultUrl)
	return resultUrl
}

func GetTestResults(wpt helper.Tester, ctx context.Context, url, location string) (string, error, int, string, string) {
	utils.SpectreLog.Debugln("Get Page load from URL: %s", url)

	summaryUrl, _, err, status := GetSummaryUrl(wpt, buildTestURL(url, location))
	utils.SpectreLog.Debugln("Status from get summary Url", status)
	if(err != nil) {
		return "", err, status, "Error while getting Summary", summaryUrl
	}

	results, err, statusCode, statusText  := GetResultsFromDetails(wpt, ctx, summaryUrl)
	return results, err, statusCode, statusText, summaryUrl
}

func GetResultsFromDetails(wpt helper.Tester, ctx context.Context, url string) (string, error, int, string) {
	utils.SpectreLog.Debugln("Details url ", url)
	results, err := getLatestFromDB(ctx, "", url, "", "")
	if err == nil  && results.SuccessfulRuns > 0 {
		utils.SpectreLog.Debugln("found results for: " + results.TestUrl)
		bytes, err := json.Marshal(results)
		if err != nil {
			return "", err, 500, "Error Marshalling test data"
		}
		return string(bytes), nil, http.StatusOK, "Saved Response"
	}
	content, err, status := wpt.GetContent(url)
	utils.SpectreLog.Debugln("Status from get summary Url  ", status)
	if err != nil {
		return "", err, status, "Error Getting Content from details Url"
	}

	err, testingJsonResponse := UnmarshalTestingResponse(content)

	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		_, esErr := SaveToES(ctx, url, details, nil, nil, nil)
		if esErr != nil {
			utils.SpectreLog.Errorf("Failed to save to Elastic Search")
		}
		return string(content), err, details.StatusCode, details.StatusText
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return string(content), err, testingJsonResponse.StatusCode, testingJsonResponse.StatusText
	}
}

func GetSummaryUrl(wpt helper.Tester, url string) (string, string, error, int) {
	//Get Summary Json
	//Parse Json result
	//Get Details URL
	//Get Details from WPT
	//Parse result
	//Get PageLoad
	content, err, status := wpt.GetContent(url)
	if(err != nil) {
		return "", "", err, status
	}

	// Note that this will require you add fmt to your list of imports.
	var summary models.Summary
	err1 := json.Unmarshal(content, &summary)

	if err1 != nil {
		return "", "", err1, status
	}
	if summary.StatusCode >= 400 {
		return "", "", errors.New(summary.StatusText), summary.StatusCode
	}

	utils.SpectreLog.Debugln("Summary URL: " + summary.Data.JsonUrl)
	return summary.Data.JsonUrl, summary.Data.TestId, nil, status
}

func UnmarshalNormalResponse(content []byte) (error, models.JsonResults1) {
	var details models.JsonResults1
	err := json.Unmarshal(content, &details)
	return err, details
}

func UnmarshalTestingResponse(content []byte) (error, models.JsonForTestStarted) {
	var details models.JsonForTestStarted
	err := json.Unmarshal(content, &details)
	return err, details
}

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	client := http.Client{}
	for _, url := range urls {
		go func(url string) {
			utils.SpectreLog.Debugln("Fetching %s \n", url)
			resp, err := client.Get(url)
			ch <- &HttpResponse{url, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			utils.SpectreLog.Debugln("%s was fetched\n", r.url)
			if r.err != nil {
				utils.SpectreLog.Debugln("with an error", r.err)
			}
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			utils.SpectreLog.Debugln(".")
		}
	}
	return responses
}

func CanonicalizeUrl(rawurl string) (string, error){
	utils.SpectreLog.Debugln("Test URL %s \n", rawurl)
	if rawurl == "" {
		return rawurl, errors.New("Invalid url.")
	}
	canonicalURL, parseerr := neturl.Parse(rawurl)
	if parseerr != nil {
		return rawurl, errors.New(parseerr.Error())
	}
	if canonicalURL.Scheme == "" {
		canonicalURL.Scheme = "http"
	}
	return canonicalURL.String(), nil
}

func SaveToES(ctx context.Context, testUrl string, details models.JsonResults1, domains []models.DomainResult, contentTypes []models.ContentType, states []models.State) (string, error){

	utils.SpectreLog.Debug("Entering handlers.pagetest.SaveToES .....")

	dbResults, ok := FetchESFromContext(ctx)
	if ok {
		results, err := dbResults.GetTest(testUrl)
		if err == nil && results.SuccessfulRuns > 0 {
			return results.ID, nil
		}
		perfscore := models.PerformanceScore{
			Cache: details.Data.Average.FirstView.ScoreCache,
			CDN: details.Data.Average.FirstView.ScoreCdn,
			Gzip: details.Data.Average.FirstView.ScoreGzip,
			Cookies: details.Data.Average.FirstView.ScoreCookies,
			KeepAlive: details.Data.Average.FirstView.ScoreKeep_alive,
			Minify: details.Data.Average.FirstView.ScoreMinify,
			Combine: details.Data.Average.FirstView.ScoreCombine,
			Compress: details.Data.Average.FirstView.ScoreCompress,
			Etags: details.Data.Average.FirstView.ScoreEtags,
			Overall: int(getAverageScore(details)),
		}
		if domains == nil {
			domains = getDomains(details)
		}
		if contentTypes == nil {
			contentTypes = getContentTypes(details)
		}
		if states == nil {
			states = getStates(details)
		}
		result := models.TestResultsES{
			ApplicationId: "", //TODO review
			PageTestUrl: details.Data.URL,
			TestUrl: testUrl,
			Runtime: time.Now(),
			User: models.User{UserId: ""}, //TODO review
			Location: strings.Split(details.Data.Location, ":")[0],
			EndpointLocation: "Unknown", //TODO
			PerformanceScore: perfscore,
			Pageload: models.Pageload{Loadtime: details.Data.Average.FirstView.LoadTime},
			Domains:  domains,
			ContentTypes: contentTypes,
			States: states,
			SuccessfulRuns: details.Data.SuccessfulFVRuns,
		}

		utils.SpectreLog.Debug("Inserting  pagetest results for testUrl %s", testUrl)
		id, err := dbResults.InsertTest(result)
		return id, err
	}
	return "", errors.New("Failed to save to DB")
}

func SaveContentToES(ctx context.Context, testUrl string, content []byte, domains []models.DomainResult, contentTypes []models.ContentType, states []models.State) (string, error) {
	err, details := UnmarshalNormalResponse(content)
	if err != nil {
		return "", errors.New("Error UnMarshalling Json Results")
	}
	return SaveToES(ctx, testUrl, details, domains, contentTypes, states)
}

func getLatestFromDB(ctx context.Context, url, lookupUrl, location, interval string) (*models.TestResultsES, error){
	dbResults, ok := FetchESFromContext(ctx)
	fmt.Printf("Looking up results in ES for url %s or lookupUrl %s\n", url, lookupUrl)
	if ok {
		if lookupUrl != "" {
			result, err := dbResults.GetTest(lookupUrl)
			if err == nil {
				fmt.Printf("Found result in ES for lookupURL %s\n", lookupUrl)
				return result, nil
			}
		}
		if url != "" {
			result, err := dbResults.GetLatestForUrl(url, location, interval)
			if err == nil {
				fmt.Printf("Found result in ES for url %s\n", url)
				return result, nil
			}
		}
	}
	fmt.Println("Not found in elastic")
	return nil, errors.New("Not found in elastic.")
}


func getDomains(details models.JsonResults1) []models.DomainResult{
	domains, err, _,_ := ParseDomainInfo(details)
	if err == nil {
		return domains
	}
	return nil
}

func getContentTypes(details models.JsonResults1) []models.ContentType {
	contentTypes, err, _,_ := ParseContentTypesInfo(details)
	if err == nil {
		return contentTypes
	}
	return nil
}

func getStates(details models.JsonResults1) []models.State {
	states, err, _,_ := ParseStatesInfo(details)
	if err == nil {
		return states
	}
	return nil
}

func GetPageTestsForApp( w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debugln("Inside GetPageTestsForApp")
	applicationId := r.URL.Query().Get("applicationId")
	if applicationId == "" {
		utils.PrintError(w, http.StatusBadRequest, "ApplicationId cannot be null")
		return
	}
	utils.SpectreLog.Debugln("ApplicationId %T %s \n", applicationId, applicationId)
	startPage := utils.ConvertToint( r.URL.Query().Get("startPage"), 1)
	utils.SpectreLog.Debugln("start page %T %v \n", startPage, startPage)
	numRows := utils.ConvertToint( r.URL.Query().Get("numRows"), 10)
	utils.SpectreLog.Debugln("Number of Rows %T %v \n", numRows, numRows)
	if startPage < 1 {
		startPage = 1
	}
	startRow := (startPage - 1) * numRows
	utils.SpectreLog.Debugln("Start Row %T %v \n", startRow, startRow)
	endRow := startRow + numRows
	utils.SpectreLog.Debugln("End Row %T %v \n", endRow, endRow)
	interval := r.URL.Query().Get("interval")

	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)
	utils.SpectreLog.Debugln("After NewTestResultDB %T %v \n", endRow, endRow)
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)
	utils.SpectreLog.Debugln("After NewTestResultES %T %v \n", endRow, endRow)
	_, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	utils.SpectreLog.Debugln("After FetchRCacheFromContext %T %v \n", endRow, endRow)
	urlsDB,  _ := persist.NewUrlDB()
	//get time period. set to default, when there is none
	// uniqueUrls,_ := urlsDB.GetAllUrls()
	uniqueUrls,_ := urlsDB.GetAllUrlsForApp( applicationId)
	var testOverview models.TestOverview
	testOverview.UrlStats = make([]models.UrlStat, 0)
	testOverview.Count = len(uniqueUrls)
	utils.SpectreLog.Debugln("After testOverview.Count %T %v \n", testOverview.Count, testOverview.Count)
	// overviewObjs := make([]models.TestOverview, 0,1)
	for i, rawurl := range uniqueUrls {
		if( i > endRow){
			break
		}
		if( i >= startRow && i< endRow) {
			// fmt.Println( "v " + rawurl)
			url, parseerr := CanonicalizeUrl(rawurl.Url)
			if( parseerr != nil){
				continue
			}

			//for every url get performance Score
			perfScore, err1, status, statusText, lookupUrl := GetPerfScoreFromCache(helper.WebPageTester{}, ctx, url, constants.DEFAULT_LOCATION)
			if(err1 != nil) {
				utils.PrintError(w, status, err1.Error())
				return
			}
			// overview.PerformanceScore = perfScore
			utils.SpectreLog.Debugln( "perfScore %T %d",perfScore, perfScore)
			utils.SpectreLog.Debugln("Status: %d\n ", status)
			utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
			utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)
			// pageLoadTime, err1, status, statusText, lookupUrl := GetPageLoadFromCache(helper.WebPageTester{}, ctx, url, constants.DEFAULT_LOCATION)
			_, pageLoadTime, availability, err1, status, statusText := GetAvailability(helper.WebPageTester{}, ctx, url, constants.DEFAULT_LOCATION, interval)
			if(err1 != nil) {
				utils.PrintError(w, status, err1.Error())
				return
			}
			// overview.PerformanceScore = perfScore
			utils.SpectreLog.Debugln( "availability %T %d",availability, availability)
			utils.SpectreLog.Debugln( "pageLoadTime %T %d",pageLoadTime, pageLoadTime)
			utils.SpectreLog.Debugln("Status: %d\n ", status)
			utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
			utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)

			//for every url get page load
			results, _ := resultsES.GetUrl(url, constants.DEFAULT_LOCATION, interval)

			var times []models.LoadtimeHistory
			var perfScores []models.ScoreHistory
			var totalLoadtime int
			var averageLoadtime int
			count := len(results)
			utils.SpectreLog.Debug( "count!!!!!! ", count)
			if count > 0 {
				for _, result := range results {
					times = append(times, models.LoadtimeHistory{Loadtime: result.Loadtime, Runtime: result.Runtime.Unix()})
					totalLoadtime += result.Loadtime
					utils.SpectreLog.Debugln( "result.PerformanceScore %T %s",result.PerformanceScore.Overall, result.PerformanceScore.Overall)
					perfScores =  append(perfScores, models.ScoreHistory{Score: result.PerformanceScore.Overall, Runtime: result.Runtime.Unix()})
				}
				averageLoadtime = totalLoadtime / len(results)

			}
			//for every url performance score history
			//for every url get Availability history
			urlStat := models.UrlStat{
    				URL: url,
				PerformanceScore: perfScore,
				PageLoadTime: pageLoadTime,
				PageLoadHistory: times,
				PageScoreHistory: perfScores,
				AverageLoadTime: averageLoadtime,
				Availability:availability,
				DashboardId: rawurl.DashboardId,
			}
			testOverview.UrlStats = append(testOverview.UrlStats, urlStat)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(testOverview); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(http.StatusOK)

}

func GetPageTestsForAppCount( w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debugln("Inside GetPageTestsForApp")
	applicationId := r.URL.Query().Get("applicationId")
	if applicationId == "" {
		utils.PrintError(w, http.StatusBadRequest, "ApplicationId cannot be null")
		return
	}
	utils.SpectreLog.Debugln("ApplicationId %T %s \n", applicationId, applicationId)
	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)
	utils.SpectreLog.Debugln("After NewTestResultDB %T %v \n")
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	utils.SpectreLog.Debugln("After FetchRCacheFromContext %T %v \n")
	cacheKey := utils.ApplicationUrlCountCacheKey( applicationId)
	count := 0;
	tempurlCount,errC := cache.GetString( cacheKey);
	if( errC == nil && tempurlCount != "") {
		count = utils.ConvertToint(tempurlCount, 0)
	}else {
		urlsDB,_ := persist.NewUrlDB()
		//get time period. set to default, when there is none
		// uniqueUrls,_ := urlsDB.GetAllUrls()
		count,_ = urlsDB.GetUrlCountForApp( applicationId)
		cache.SetStringWithExpiration(cacheKey, utils.ConvertIntToString(count), constants.PAGELOOKUP_EXPIRATION)
	}

	testOverview := models.TestOverview{
		Count: count,
		UrlStats: nil,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(testOverview); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(http.StatusOK)

}

func URLExists( w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		utils.SpectreLog.Errorln("Token is not valid")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	utils.SpectreLog.Debugln("Token is Valid")

	//TODO get access token
	// if access token not there get UserID and customerID from query params
	//add user and customer ID to DB

	//TODO change this to get from RBAC layer
	rawurl := r.URL.Query().Get("url")

	//userId := r.URL.Query().Get("userId")
	//customerId := r.URL.Query().Get("customerId")

	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	applicationId := r.URL.Query().Get("applicationId")
	if applicationId == "" {
		utils.PrintError(w, http.StatusBadRequest, "ApplicationId cannot be null")
		return
	}
	urlsDB,_ := persist.NewUrlDB()
	//get time period. set to default, when there is none
	// uniqueUrls,_ := urlsDB.GetAllUrls()
	hasUrl := urlsDB.HasUrl(url, applicationId)
	var dashboardId string
	if hasUrl {
		urlInfo, err := urlsDB.GetUrl(url, applicationId)
		if err != nil {
			//log.Debugln("Error finding url in DB.")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(models.DashboardExists{Exists: false}); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		dashboardId = urlInfo.DashboardId
	}
	exists := models.DashboardExists{Exists: hasUrl, DashboardId: dashboardId}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(exists); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(http.StatusOK)
}