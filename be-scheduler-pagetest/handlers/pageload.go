package handlers

import (
	"net/http"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"
	"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	// "strings"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// Page Load API returns summary URL. Json URL is used to lookup detailed results in LookupPageLoad API
func PageLoad(w http.ResponseWriter, r *http.Request) {
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

	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}

	//TODO change this to get from RBAC layer
	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	appId := r.URL.Query().Get("applicationId")
	dashboardId := r.URL.Query().Get("dashboardId")
	//Insert URL in MongoDB System URLS
	insertSystemUrl(url, appId, dashboardId)
	//count is reset with new url. so remove object from cache.
	cache.Delete(utils.ApplicationUrlCountCacheKey( appId))
	wpt := helper.WebPageTester{}

	var status int

	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}
	pageLoadTime, err1, status, statusText, lookupUrl := GetPageLoadFromCache(wpt, ctx, url, location, "")

	if(err1 != nil) {
		utils.PrintError(w, status, err1.Error())
		return
	}

	utils.SpectreLog.Debugln("Page Load Time: %d\n ", pageLoadTime)
	utils.SpectreLog.Debugln("Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)
	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.PageLoadTime = pageLoadTime
	results.AppId = appId
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}

// Lookup Status of test using JSON URL returned from PageLoad API and get test results if test completed
func LookupPageLoad(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering LookupPageLoadHandler()")
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
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	rawurl := r.URL.Query().Get("lookupurl")
	utils.SpectreLog.Debugln("Test URL %s \n", rawurl)
	if rawurl == "" {
		utils.PrintError(w, http.StatusBadRequest, "Invalid value in `value` parameter.")
		return
	}

	pageLoadTime, err, status, statusText, lookupUrl := LookupPageLoadFromCache(helper.WebPageTester{}, ctx, "", rawurl, "")
	if(err != nil) {
		utils.PrintError(w, status, err.Error())
		return
	}
	utils.SpectreLog.Debugln("Page Load Time: %d\n ", pageLoadTime)
	utils.SpectreLog.Debugln("Test Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.PageLoadTime = pageLoadTime

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}

func LookupPageLoadFromCache(wpt helper.Tester, ctx context.Context, url, lookupUrl, interval string) (int, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}

	var content []byte
	var status int
	var statusText string
	if results, _ := cache.GetString(constants.CACHE_PREFIX + lookupUrl); results == "" {
		results, err :=getLatestFromDB(ctx, url, lookupUrl, "", interval)
		if err == nil  && results.Pageload.Loadtime > 0 {
			return results.Pageload.Loadtime, nil, http.StatusOK, "Saved Response", lookupUrl
		}
		content, err, status := wpt.GetContent(lookupUrl)
		if err != nil {
			return 0, err, status, "Error Retrieving Content", lookupUrl
		}
		loadTime, err, status, statusText := GetPageLoadFromDetails(ctx, lookupUrl, content, status, true)
		return loadTime, err, status, statusText, lookupUrl
	} else {
		content = []byte(results)
	}
	status = 200
	statusText = "Cached Response"

	pageLoadTime, err, status, statusTextDet := GetPageLoadFromDetails(ctx, lookupUrl, content, status, false) //TODO base isNew off of FindTest
	if statusText == "" {
		statusText = statusTextDet
	}

	return pageLoadTime, err, status, statusText, lookupUrl
}

func insertSystemUrl(url, appId, dashboardId string)  {
	//Add http if not added
	resultUrl := url
	urlDb,  err := persist.NewUrlDB()
	if err == nil {
		utils.SpectreLog.Debug("Inserting system Url in Db: %s ", url)
		modelToStore := new(models.SystemUrl)
		modelToStore.Url = resultUrl
		modelToStore.AppId = appId
		modelToStore.DashboardId = dashboardId
		id, err1 := urlDb.InsertUrl(modelToStore)
		if (err1 != nil) {
			utils.SpectreLog.Warn("Url %s not inserted due to error: ", url, err1)
		} else {
			utils.SpectreLog.Debug("Successfully inserted url: %s with id ", url, id)
		}
	} else {
		utils.SpectreLog.Error("Error getting handle to DB %s", err)
	}
}

func GetPageLoadFromCache(wpt helper.Tester, ctx context.Context, url, location, interval string) (int, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		utils.SpectreLog.Debugln("Lookup Url is %s \n", lookupUrlFromCache)
		loadTime, err, status, statusText, lookupUrl := LookupPageLoadFromCache(wpt, ctx, url, lookupUrlFromCache, interval)
		//if err == nil {
		return loadTime, err, status, statusText, lookupUrl
		//}
	}
	results, err := getLatestFromDB(ctx, url, lookupUrlFromCache, location, interval)

	if err == nil  && results.Pageload.Loadtime > 0 {
		return results.Pageload.Loadtime, nil, http.StatusOK, "Saved Response", lookupUrlFromCache
	}
	pageLoadTime, err, status, statusText, lookupUrl := GetPageLoad(wpt, ctx, url, location)
	if(err != nil) {
		return pageLoadTime, err, status, statusText, lookupUrl
	}

	cache.SetStringWithExpiration(constants.CACHE_PREFIX + url + "::" + location, lookupUrl, constants.PAGELOOKUP_EXPIRATION)
	return pageLoadTime, err, status, statusText, lookupUrl
}

func GetPageLoad(wpt helper.Tester, ctx context.Context, url, location string) (int, error, int, string, string) {
	utils.SpectreLog.Debugln("Get Page load from URL: %s", url)
	resultUrl := buildTestURL(url, location)
	utils.SpectreLog.Debugln("")
	utils.SpectreLog.Debugln("Result URL: %s\n", resultUrl)
	utils.SpectreLog.Debugln("")
	summaryUrl, _, err, status := GetSummaryUrl(wpt, resultUrl)

	if(err != nil) {
		return 0, err, status, "Error getting Json Summary Results", summaryUrl
	}

	content, err, status := wpt.GetContent(summaryUrl)
	if err != nil {
		return 0, err, status, "Error Retrieving Content", summaryUrl
	}

	pageLoadTime, err, status, statusText := GetPageLoadFromDetails(ctx, summaryUrl, content, status, true)
	return pageLoadTime, nil, status, statusText, summaryUrl
}

func GetPageLoadFromDetails(ctx context.Context, summaryUrl string, content []byte, status int, isNew bool) (int, error, int, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}

	err, testingJsonResponse := UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		if err != nil {
			return 0, err, http.StatusInternalServerError, "Error UnMarshalling Json Results"
		}

		if details.Data.SuccessfulFVRuns != 0 {
			utils.SpectreLog.Debugln("Successful Runs found!!")
			utils.SpectreLog.Debugln("Average Load Time  is %d \n", details.Data.Average.FirstView.LoadTime)
			//Update or Insert in DB
			SaveToES(ctx, summaryUrl, details, nil, nil, nil) //TODO user and app ids

			cache.SetStringWithExpiration(constants.CACHE_PREFIX + summaryUrl, string(content), constants.TESTRESULTS_EXPIRATION)
			return details.Data.Average.FirstView.LoadTime, nil, details.StatusCode, details.StatusText
		} else {
			utils.SpectreLog.Debugln("Successful Runs not found!!")
			return 0, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text), http.StatusBadRequest, details.Data.Runs.One.FirstView.ConsoleLog[0].Text
		}
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return 0, err, testingJsonResponse.Data.StatusCode, testingJsonResponse.Data.StatusText
	}
}

func upsert(ctx context.Context, summaryUrl string, details models.JsonResults1, isNew bool) (string, error){
	resultsES, ok := FetchESFromContext(ctx)
	if !ok {
		utils.SpectreLog.Debugln("Unable to find DB in context.")
		return "", errors.New("Unable to find DB in context.")
	}
	if isNew && ok {
		results, err := resultsES.GetTest(summaryUrl)
		utils.SpectreLog.Debugln(err)

		if err == nil  && results.TestUrl != "" {
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
			result := models.TestResultsES{
				ApplicationId: results.ApplicationId,
				PageTestUrl: details.Data.URL,
				TestUrl: summaryUrl,
				Runtime: time.Now(), //TODO
				User: models.User{UserId: results.User.UserId},
				Location: details.Data.Location,
				EndpointLocation: "Unknown",//TODO
				PerformanceScore: perfscore,
				Pageload: models.Pageload{Loadtime: details.Data.Average.FirstView.LoadTime},
				Domains: getDomains(details),
				ContentTypes: getContentTypes(details),
				States: getStates(details),
				SuccessfulRuns: details.Data.SuccessfulFVRuns,
			}
			err := resultsES.UpdateTest(results.ID, result)
			if err == nil {
				utils.SpectreLog.Debugln("Updated to db as ", results.TestUrl)
				return result.ID, nil
			} else {
				utils.SpectreLog.Debugln("Failed to save to DB")
				return "", err
			}
		} else if err != nil && isNew {
			id, err := SaveToES(ctx, summaryUrl, details, nil, nil, nil)
			return id, err
		}
	}
	return "", errors.New("Nothing to update.")
}

func PageLoadChart(w http.ResponseWriter, r *http.Request) {

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
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	//TODO change this to get from RBAC layer
	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}

	location := r.URL.Query().Get("location")
	interval := r.URL.Query().Get("interval")
	results, err := resultsES.GetUrl(url, location, interval)
	var times []models.LoadtimeHistory
	var totalLoadtime int
	var averageLoadtime int
	if len(results) > 0 {
		for _, result := range results {
			times = append(times, models.LoadtimeHistory{Loadtime: result.Loadtime, Runtime: result.Runtime.Unix()})
			totalLoadtime += result.Loadtime
		}
		averageLoadtime = totalLoadtime / len(results)
	}

	pageloadHistory := models.PageLoadHistory{LoadtimeHistory: times, AverageLoadtime: averageLoadtime}
	if(err != nil) {
		utils.PrintError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(pageloadHistory); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}
