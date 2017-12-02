package handlers

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"errors"
	"strconv"
	"strings"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

//// PerfScore API returns summary URL. Json URL is used to lookup detailed results in LookupPerfScore API

func PerfScoreDetails(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering PerfScoreDetails()")

	if !security.ValidateRequest(w, r, constants.SCOPE_WRITE) {
		utils.SpectreLog.Debugln("token not valid")
		utils.SpectreLog.Errorf("Token is not valid -> ")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	utils.SpectreLog.Debugln("Token is Valid")

	url, location, err := extractParamsFromReq(r)
	if err != nil {
		utils.SpectreLog.Errorf("Failed to extract url query parameters: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	ctx, err := getContext()
	if (err != nil) {
		utils.SpectreLog.Errorf("Error getting context for performance score details ")
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	scoreDetailResults, err, status, statusText, lookupUrl := GetPerfScoreDetails(ctx, url, location)
	if (err != nil) {
		utils.SpectreLog.Errorf("Error getting performace details for url %s and location %s\n ", url, location )
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	utils.SpectreLog.Debugf("Status: %d\n ", status)
	utils.SpectreLog.Debugf("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugf("Lookup Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.ScoreDetailsResult = scoreDetailResults
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}

func GetPerfScoreDetails(ctx context.Context, url string, location string) (models.ScoreDetailResult, error, int, string, string) {
	wpt := helper.WebPageTester{}
	//TODO get lookup URL from Cache
	// If not in Cache get from ES
	scoreDetailResults, err, status, statusText, lookupUrl := GetPerfScoreDetailsFromCache(ctx, wpt, url, location)

	if err != nil  {
		utils.SpectreLog.Debug("Not found in Cache. Getting score details from ElasticSearch")
		utils.SpectreLog.Debugln("Not found in Cache. Getting details from ElasticSearch")
		scoreDetailResults, err, status, statusText, lookupUrl = GetPerfScoreDetailsFromES(ctx, wpt, lookupUrl, url)

		if err != nil || scoreDetailResults.Overallscore == "" {
			utils.SpectreLog.Debugf("Not found in ElasticSearch.")

			if lookupUrl == "" {
				utils.SpectreLog.Debugf("Lookup Url and results not found in ElasticSearch. Running live test for url %s \n", url)
				testUrl := buildTestURL(url, location)
				lookupUrl, _, err, status = GetSummaryUrl(wpt, testUrl)
				statusText = "Live Test"
				if (err != nil) {
					utils.SpectreLog.Errorf("Failed to get summary URL %s for url ", testUrl)
					return models.ScoreDetailResult{}, err, status, "Unable to live results from WPT", lookupUrl
				}
			}
			utils.SpectreLog.Debugf("Retrieving live results for lookupUrl %s from WPT\n", lookupUrl)
			var content []byte
			content, err, status = wpt.GetContent(lookupUrl)
			if err == nil && content != nil {
				utils.SpectreLog.Debugf("Got content for lookup URL %s", lookupUrl)
				scoreDetailResults, err, status, statusText = GetScoreDetailsAndSetCache(ctx, lookupUrl, content, true)
				SetLookupUrlInCache(ctx, url, location, lookupUrl)
				return scoreDetailResults, err, status, statusText, lookupUrl
			} else {
				utils.SpectreLog.Errorf("Failed to get results from WPT for Lookup URL %s ", lookupUrl)
				return models.ScoreDetailResult{}, err, status, statusText, lookupUrl
			}
		} else if scoreDetailResults.Overallscore != "" {
			//Set  Cache if result found in ElasticSearch
			utils.SpectreLog.Debugf("Setting Cache after finding results in ES")
			if lookupUrl == "" {
				utils.SpectreLog.Debugf("Lookup Url and results not found. Running live test for url %s \n", url)
				testUrl := buildTestURL(url, location)
				lookupUrl, _, err, status = GetSummaryUrl(wpt, testUrl)
				statusText = "Live Test"
				if (err != nil) {
					utils.SpectreLog.Errorf("Failed to get summary URL %s for url ", testUrl)
					return models.ScoreDetailResult{}, err, status, "Unable to live results from WPT", lookupUrl
				}
			}
			var content []byte
			SetLookupUrlInCache(ctx, url, location, lookupUrl)
			content, err, status = wpt.GetContent(lookupUrl)
			if (err == nil) {
				SetResultsInCache(ctx, lookupUrl, content)
				err = nil
				statusText = "Found result for lookup url " + lookupUrl
			} else {
				utils.SpectreLog.Warnf("Cannot get content from WPT for lookupURL %s. Unable to set cache with results")
			}
			return scoreDetailResults, err, status, statusText, lookupUrl
		} else {
			utils.SpectreLog.Error("Error getting performance for url %s ", url)
			return scoreDetailResults, err, status, statusText, lookupUrl
		}
	}

	return scoreDetailResults, err, status, statusText, lookupUrl
}

func extractParamsFromReq(req *http.Request) (string, string, error) {
	rawurl := req.URL.Query().Get("url")

	if(rawurl == "") {
		return "", "", errors.New("Url not set")
	}
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		return "", "", parseerr
	}

	location := req.URL.Query().Get("location")
	if(location == "") {
		return "", "", errors.New("location not set")
	}
	return url, location, nil
}

// Lookup Status of test using JSON URL returned from PageLoad API and get test results if test completed
func LookupPerfScoreDetails(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering LookupPerfScoreDetails()")
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}
	utils.SpectreLog.Debugln("Token is Valid")
	ctx, err := getContext()

	if(err != nil) {
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	lookupUrl := r.URL.Query().Get("lookupurl")
	utils.SpectreLog.Debugf("Test URL %s \n", lookupUrl)

	if lookupUrl == "" {
		utils.PrintError(w, http.StatusBadRequest, "Lookupurl not set in query parameter")
		utils.SpectreLog.Debugln("Leaving LookupPerfScoreDetails()")
		return
	}
	wpt := helper.WebPageTester{}

	scoreDetailResult, err, status, statusText := GetPerfScoreDetailsWithLookUrl(ctx, wpt, lookupUrl)
	utils.SpectreLog.Debugln("Test Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("lookupUrl Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.ScoreDetailsResult = scoreDetailResult
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		utils.SpectreLog.Debug("Leaving LookupPerfScoreHandler()")
		return
	}
}

func GetPerfScoreDetailsWithLookUrl(ctx context.Context, wpt helper.WebPageTester, lookupUrl string) (models.ScoreDetailResult, error, int, string) {
	scoreDetailResult, err, status, statusText := GetScoreDetailsWithLookupUrlFromCache(wpt, ctx, lookupUrl)

	if err != nil || status != http.StatusOK {
		//Get Details from ES
		scoreDetailResult, err, status, statusText = GetScoreDetailsWithLookupUrlFromES(ctx, wpt, lookupUrl)
		if err != nil  || status == http.StatusNotFound {
			//Get Live from WPT
			utils.SpectreLog.Debugf("Getting results for lookupurl %s live from web page test", lookupUrl)
			content, err, _ := wpt.GetContent(lookupUrl)
			if err == nil && content != nil {
				scoreDetailResult, err, status, statusText = GetScoreDetailsAndSetCache(ctx, lookupUrl, content, true)
				return scoreDetailResult, err, status, statusText
			} else {
				utils.SpectreLog.Errorf("Failed to get response for Lookup URL %s ",  lookupUrl)
				//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
				return scoreDetailResult, err, status, statusText
			}
		} else if status == http.StatusOK && (scoreDetailResult.Overallscore != "" || scoreDetailResult.Overallscore != "0") {
			utils.SpectreLog.Debugf("Retrieving live results for lookupUrl %s from WPT\n to set cache", lookupUrl)
			var content []byte
			content, err, status = wpt.GetContent(lookupUrl)
			if err == nil && content != nil {
				utils.SpectreLog.Debugf("Retrieved content for lookup URL %s", lookupUrl)
				utils.SpectreLog.Debug("Setting Cache")
				SetResultsInCache(ctx, lookupUrl, content)
			} else {
				utils.SpectreLog.Warn("Cannot get live results from WPT for lookupUrl %s. Cannot set Cache with results", lookupUrl)
				//utils.SpectreLog.Errorf("Failed to get results from WPT for Lookup URL %s ", lookupUrl)
			}
			return scoreDetailResult, err, status, statusText
		} else {
			utils.SpectreLog.Errorf("Error getting results from elastic search for lookup url %s\n", lookupUrl)

		}
	}
	return scoreDetailResult, err, status, statusText
}

func getContext() (context.Context, error){
	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return nil, err
	}
	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return nil, err
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)
	 helper.FetchRCacheFromContext(ctx)
	return ctx, nil
}

func GetPerfScoreDetailsFromCache(ctx context.Context, wpt helper.Tester,  url, location string) (models.ScoreDetailResult, error, int, string, string) {
	utils.SpectreLog.Debug("Entering GetPerfScoreDetailsFromCache()")
	utils.SpectreLog.Debugf("Getting lookupurl from Cache for url %s\n", url)
	lookupUrlFromCache, err := GetLookupUrlFromCache(ctx, url, location )
	if err == nil && lookupUrlFromCache != "" {
		utils.SpectreLog.Debugf("Found lookup in Cache %s\n", lookupUrlFromCache)
		if strings.Contains(lookupUrlFromCache, "jsonResult.php?test=") {
			results, err, status, statusText := GetScoreDetailsWithLookupUrlFromCache(wpt, ctx, lookupUrlFromCache)
			utils.SpectreLog.Debugf("Status %d for getting score details with lookupurl %s\n", status, lookupUrlFromCache)
			return results, err, status, statusText, lookupUrlFromCache
		} else {
			return models.ScoreDetailResult{}, errors.New("Invalid lookupurl"), http.StatusNotFound, "Lookup URL Not Found in Cache", lookupUrlFromCache
		}
	}
	utils.SpectreLog.Debugf("Url not found in Cache")
	utils.SpectreLog.Debug("Leaving GetPerfScoreDetailsFromCache()")
	return models.ScoreDetailResult{}, err, http.StatusNotFound, "Url not found in Cache: " + url, lookupUrlFromCache
}

func GetLookupUrlFromCache(ctx context.Context,  url string, location string) (string, error) {
	cache, _ := helper.FetchRCacheFromContext(ctx)
	utils.SpectreLog.Debugf("Getting lookupurl from cache %s\n", url)
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	utils.SpectreLog.Debugf("Retrieved lookupurl from cache: %s\n", lookupUrlFromCache)
	utils.SpectreLog.Debug("Leaving GetLookupUrlFromCache()")
	return lookupUrlFromCache, err
}

func SetLookupUrlInCache(ctx context.Context,  url string, location string, lookupUrl string) {
	utils.SpectreLog.Debugf("Setting lookup URl %s in Cache for Url %s\n", lookupUrl, url)
	utils.SpectreLog.Debugf("Setting lookup URl %s in Cache for Url %s", lookupUrl, url)
	cache, _ := helper.FetchRCacheFromContext(ctx)
	cache.SetStringWithExpiration(constants.CACHE_PREFIX + url + "::" + location, lookupUrl, constants.PAGELOOKUP_EXPIRATION)
}

func SetResultsInCache(ctx context.Context,  lookupUrl string, content []byte) {
	utils.SpectreLog.Debugf("Setting result  in Cache for LookupUrl %s\n", lookupUrl)
	utils.SpectreLog.Debugf("Setting result  in Cache for LookupUrl %s\n", lookupUrl, lookupUrl)
	cache, _ := helper.FetchRCacheFromContext(ctx)
	cache.SetStringWithExpiration(constants.CACHE_PREFIX + lookupUrl, string(content), constants.TESTRESULTS_EXPIRATION)
}

func GetPerfScoreDetailsFromES(ctx context.Context, wpt helper.Tester,  lookupUrl string, url string) (models.ScoreDetailResult, error,  int , string, string) {
	utils.SpectreLog.Debug("Entering GetPerfScoreDetailsFromES()")
	utils.SpectreLog.Debugf("Getting results from ElasticSearch for url %s and lookup %s \n", url, lookupUrl)
	results, err := getLatestFromDB(ctx, url, lookupUrl, "", "")
	if err == nil  && results != nil  {
		utils.SpectreLog.Debugln("Retrieving score details from elasticsearch result")
		scoreDetailResult, _ := getPerfScoreDetailsFromESResult(results)
		return scoreDetailResult, nil, http.StatusOK, "Saved Response", results.TestUrl
	}

	return models.ScoreDetailResult{}, err, http.StatusNotFound, "Not Found", ""
}

func getPerfScoreDetailsFromESResult(esResult *models.TestResultsES)  (models.ScoreDetailResult, string) {
	var scoreDetailResult models.ScoreDetailResult
	scoreDetails := make([]models.ScoreDetail, 0)
	scoreDetails = buildScoreDetail(helper.KEEPALIVE, esResult.PerformanceScore.KeepAlive, scoreDetails )
	scoreDetails = buildScoreDetail(helper.CACHE, esResult.PerformanceScore.Cache, scoreDetails )
	scoreDetails = buildScoreDetail(helper.GZIP, esResult.PerformanceScore.Gzip, scoreDetails )
	scoreDetails = buildScoreDetail(helper.CDN, esResult.PerformanceScore.CDN, scoreDetails )
	scoreDetails = buildScoreDetail(helper.COMBINE, esResult.PerformanceScore.Combine, scoreDetails )
	scoreDetails = buildScoreDetail(helper.COMPRESS, esResult.PerformanceScore.Compress, scoreDetails )
	scoreDetails = buildScoreDetail(helper.COOKIES, esResult.PerformanceScore.Cookies, scoreDetails )
	scoreDetails = buildScoreDetail(helper.ETAGS, esResult.PerformanceScore.Etags, scoreDetails )
	scoreDetails = buildScoreDetail(helper.MINIFY, esResult.PerformanceScore.Minify, scoreDetails )

	scoreDetailResult.Overallscore = strconv.Itoa(esResult.PerformanceScore.Overall) + "/100"
	scoreDetailResult.Individualscore = scoreDetails
	return scoreDetailResult, esResult.TestUrl
}

func buildScoreDetail(category string, percent int, scoreDetails []models.ScoreDetail) ([]models.ScoreDetail){
	if(percent != -1) {
		var scoreDetail models.ScoreDetail
		scoreDetail.Category = category
		scoreDetail.Score = percent
		scoreDetail.Suggestion = helper.BuildSuggestion(scoreDetail.Category, percent)
		scoreDetail.SuggestionDetails = helper.BuildSuggestionDetails(scoreDetail.Category, percent)
		scoreDetails = append(scoreDetails, scoreDetail)
	}
	return scoreDetails
}

func GetResultsFromCache(wpt helper.Tester, ctx context.Context, lookupUrl string) ([]byte,  error, int, string) {
	cache, _ := helper.FetchRCacheFromContext(ctx)
	var content []byte
	var err error
	status := http.StatusNotFound
	var statusText string

	if results, _ := cache.GetString(constants.CACHE_PREFIX + lookupUrl); results != "" {
		utils.SpectreLog.Debugf("Found result for lookup url %s in Cache \n", lookupUrl)

		content = []byte(results)
		status = http.StatusOK
		statusText = "Cached Response"
	} else {
		utils.SpectreLog.Debugf("Results not found in Cache for lookup url %s\n",lookupUrl)
		err = errors.New("Not Found in Cache")
		status = http.StatusNotFound
		statusText = "Not found in Cache"
	}

	return content, err, status, statusText
}

func GetScoreDetailsWithLookupUrlFromCache(wpt helper.Tester, ctx context.Context, lookupUrl string) (models.ScoreDetailResult, error, int, string) {

	utils.SpectreLog.Debug("Entering LookupPerfScoreFromCache()")
	utils.SpectreLog.Debugf("Getting Results for LookupUrl %s from Cache\n", lookupUrl)
	var jsonResults models.JsonResults1
	content, err, status, statusText := GetResultsFromCache(wpt, ctx, lookupUrl)
	if err == nil {
		jsonResults, err, status, statusText = GetDetailsFromContent(content)
		scoredetails := getScoreDetails(jsonResults)
		utils.SpectreLog.Debug("Leaving LookupPerfScoreFromCache()")
		return  scoredetails, nil, http.StatusOK, "Cached response"
	}
	return models.ScoreDetailResult{}, err, status, statusText
}

func GetScoreDetailsWithLookupUrlFromES(ctx context.Context, wpt helper.Tester,  lookupUrl string) (models.ScoreDetailResult, error,  int , string) {
	utils.SpectreLog.Debug("Entering GetPerfScoreDetailsFromES()")
	utils.SpectreLog.Debugf("Getting results from ElasticSearch for  lookup url %s \n", lookupUrl)
	results, err := getLatestFromDB(ctx, "", lookupUrl, "", "")
	if err == nil  && results != nil  {
		utils.SpectreLog.Debugf("Found results in ElasticSearch for lookup url %s\n", lookupUrl)
		scoreDetailResult, _ := getPerfScoreDetailsFromESResult(results)
		return scoreDetailResult, nil, http.StatusOK, "Saved Response"
	}
	utils.SpectreLog.Debugf("Results not found for lookup url\n")

	return models.ScoreDetailResult{}, err, http.StatusNotFound, "Not Found"
}

func GetDetailsFromContent(content []byte) (models.JsonResults1, error, int, string) {
	err, testingJsonResponse := UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		if err != nil {
			return models.JsonResults1{}, err, http.StatusInternalServerError, "Error UnMarshalling Json Results"
		}

		if details.Data.SuccessfulFVRuns != 0   {
			utils.SpectreLog.Debugln("Successful Runs found!!")

			return details, nil, details.StatusCode, details.StatusText
		} else {
			utils.SpectreLog.Debugln("Successful Runs not found!!")
			return models.JsonResults1{}, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text), http.StatusBadRequest, details.Data.Runs.One.FirstView.ConsoleLog[0].Text
		}
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return models.JsonResults1{}, err, testingJsonResponse.Data.StatusCode, testingJsonResponse.Data.StatusText
	}
}

func GetScoreDetailsAndSetCache(ctx context.Context, summaryUrl string, content []byte,  isNew bool) (models.ScoreDetailResult, error, int, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		//return 0, errors.New("Unable to find cache in context."), http.StatusInternalServerError, ""
	}

	utils.SpectreLog.Debugf("Getting perfscore details from content for lookupUrl %s\n", summaryUrl)
	var scoreDetailResult models.ScoreDetailResult

	jsonResults, err, status, statusText := GetDetailsFromContent(content)
	if err == nil && jsonResults.Data.SuccessfulFVRuns != 0 {
		utils.SpectreLog.Debugf("Found successful runs for lookup url %s. Getting score details.", summaryUrl)
		scoreDetailResult = getScoreDetails(jsonResults)

		if (isNew) {
			//Update or Insert in DB
			utils.SpectreLog.Debugf("Saving results to ElasticSearch...\n")
			SaveToES(ctx, summaryUrl, jsonResults, nil, nil, nil) //TODO user and app ids
		}
		utils.SpectreLog.Debugf("Setting Cache for summaryUrl")
		cache.SetStringWithExpiration(constants.CACHE_PREFIX + summaryUrl, string(content), constants.TESTRESULTS_EXPIRATION)
	}
	return scoreDetailResult, err, status, statusText
}

func getScoreDetails(details models.JsonResults1) (models.ScoreDetailResult) {
	scoreDetails := make([]models.ScoreDetail, 0)
	var scoreDetailResult models.ScoreDetailResult

	scoreDetails = buildScoreDetail(helper.KEEPALIVE, details.Data.Average.FirstView.ScoreKeep_alive, scoreDetails)
	scoreDetails = buildScoreDetail(helper.GZIP, details.Data.Average.FirstView.ScoreGzip, scoreDetails)
	scoreDetails = buildScoreDetail(helper.COMPRESS, details.Data.Average.FirstView.ScoreCompress, scoreDetails)
	scoreDetails = buildScoreDetail(helper.CACHE, details.Data.Average.FirstView.ScoreCache, scoreDetails)
	scoreDetails = buildScoreDetail(helper.CDN, details.Data.Average.FirstView.ScoreCdn, scoreDetails)
	scoreDetails = buildScoreDetail(helper.COMBINE, details.Data.Average.FirstView.ScoreCombine, scoreDetails)
	scoreDetails = buildScoreDetail(helper.ETAGS, details.Data.Average.FirstView.ScoreEtags, scoreDetails)
	scoreDetails = buildScoreDetail(helper.COOKIES, details.Data.Average.FirstView.ScoreCookies, scoreDetails)
	scoreDetails = buildScoreDetail(helper.MINIFY, details.Data.Average.FirstView.ScoreMinify, scoreDetails)

	overallScore := getAverageScore(details)
	scoreDetailResult.Overallscore = strconv.Itoa(int(overallScore)) + "/100"
	scoreDetailResult.Individualscore = scoreDetails

	return  scoreDetailResult
}
