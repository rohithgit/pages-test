package handlers

import (
	"net/http"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// PerfScore API returns summary URL. Json URL is used to lookup detailed results in LookupPerfScore API
func PerfScore(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering PageTestHandler()")
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
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), ctxDBKey, resultsDB)

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
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

	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}

	perfScore, err1, status, statusText, lookupUrl := GetPerfScoreFromCache(helper.WebPageTester{}, ctx, url, location)

	if(err1 != nil) {
		utils.PrintError(w, status, err1.Error())
		return
	}

	utils.SpectreLog.Debugln("Performance Score: %d\n ", perfScore)
	utils.SpectreLog.Debugln("Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.PerformanceScore = perfScore
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(status)
}

func GetPerfScoreFromCache(wpt helper.Tester, ctx context.Context, url, location string) (int, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		results, err, status, statusText, lookupUrl := LookupPerfScoreFromCache(wpt, ctx, url, lookupUrlFromCache)
		//if err == nil {
		return results, err, status, statusText, lookupUrl
		//}
	}
	results, err :=getLatestFromDB(ctx, url, lookupUrlFromCache, location, "")
	if err == nil  && results.PerformanceScore.Overall > 0 {
		return results.PerformanceScore.Overall, nil, http.StatusOK, "Saved Response", lookupUrlFromCache
	}
	perfScore, err, status, statusText, lookupUrl := GetPerfScore(wpt, ctx, url, location)
	if(err != nil) {
		return perfScore, err, status, statusText, lookupUrl
	}
	cache.SetStringWithExpiration(constants.CACHE_PREFIX + url + "::" + location, lookupUrl, constants.PAGELOOKUP_EXPIRATION)
	utils.SpectreLog.Debug("Leaving GetPerfScoreFromCache()")
	return perfScore, err, status, statusText, lookupUrl

}

func GetPerfScore(wpt helper.Tester,ctx context.Context, url, location string) (int, error, int, string, string) {
	utils.SpectreLog.Debug("Entering GetPerfScore()")
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

	perfScore, err, status, statusText := GetPerfScoreFromDetails(ctx, summaryUrl, content, status, true)
	if(err != nil) {
		return perfScore, err, status, statusText, summaryUrl
	}
	utils.SpectreLog.Debug("Leaving GetPerfScore()")
	return perfScore, nil, status, statusText, summaryUrl
}

// Lookup Status of test using JSON URL returned from PageLoad API and get test results if test completed
func LookupPerfScore(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering LookupPerfScoreHandler()")
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
	url := r.URL.Query().Get("lookupurl")
	utils.SpectreLog.Debugln("Test URL %s \n", url)

	if url == "" {
		utils.PrintError(w, http.StatusBadRequest, "Invalid value in `value` parameter.")
		utils.SpectreLog.Debug("Leaving LookupPerfScoreHandler()")
		return
	}
	perfscore, err, status, statusText, lookupUrl := LookupPerfScoreFromCache(helper.WebPageTester{}, ctx, "", url)
	if(err != nil) {
		utils.PrintError(w, status, err.Error())
		utils.SpectreLog.Debug("Leaving LookupPerfScoreHandler()")
		return
	}
	utils.SpectreLog.Debugln("Performance Score: %d\n ", perfscore)
	utils.SpectreLog.Debugln("Test Status: %d\n ", status)
	utils.SpectreLog.Debugln("Status Text: %s\n ", statusText)
	utils.SpectreLog.Debugln("lookupUrl Url: %s\n ", lookupUrl)

	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.LookupUrl = lookupUrl
	results.PerformanceScore = perfscore

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		utils.SpectreLog.Debug("Leaving LookupPerfScoreHandler()")
		return
	}
	w.WriteHeader(status)
	utils.SpectreLog.Debug("Leaving LookupPerfScoreHandler()")
}

func LookupPerfScoreFromCache(wpt helper.Tester, ctx context.Context, url string, lookupUrl string) (int, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}
	utils.SpectreLog.Debug("Entering LookupPerfScoreFromCache()")

	var content []byte
	var status int
	var statusText string
	isNew := false
	if results, _ := cache.GetString(constants.CACHE_PREFIX + lookupUrl); results == "" {
		results, err := getLatestFromDB(ctx, url, lookupUrl, "", "")
		if err == nil  && results.PerformanceScore.Overall > 0 {
			return results.PerformanceScore.Overall, nil, http.StatusOK, "Saved Response", lookupUrl
		}
		content, err, status := wpt.GetContent(lookupUrl)
		if err != nil {
			return 0, err, status, "Error Retrieving Content", lookupUrl
		}
		perfScore, err, status, statusText := GetPerfScoreFromDetails(ctx, lookupUrl, content, status, true)
		return perfScore, err, status, statusText, lookupUrl
	} else {
		content = []byte(results)
	}
	status = 200
	statusText = "Cached Response"
	utils.SpectreLog.Debugln("Lookup Url is %s \n", lookupUrl)
	perfScore, err, status, statusTextDet := GetPerfScoreFromDetails(ctx, lookupUrl, content, status, isNew)
	if statusText == "" {
		statusText = statusTextDet
	}
	if(err != nil) {
		return perfScore, err, status, statusText, lookupUrl
	}
	utils.SpectreLog.Debug("Leaving LookupPerfScoreFromCache()")
	return perfScore, nil, status, statusText, lookupUrl
}

func GetPerfScoreFromDetails(ctx context.Context, summaryUrl string, content []byte, status int, isNew bool) (int, error, int, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}
	utils.SpectreLog.Debug("Entering GetPerfScoreFromDetails()")
	details, err, status, statusText := GetDetailsFromContent(content)

	if (err != nil) || ( status < http.StatusOK) {
		return 0, err, status, statusText
	}

	if details.Data.SuccessfulFVRuns != 0 {
		utils.SpectreLog.Debugln("Successful Runs found!!")
		score := getAverageScore(details)
		utils.SpectreLog.Debugln("Average Perf Score  is %f \n", score)
		if isNew {
			id, err := SaveToES(ctx, summaryUrl, details, nil, nil, nil)

			if err == nil && id != "" {
				utils.SpectreLog.Debugln("Saved to ES as " + id)
			} else {
				utils.SpectreLog.Debugln("Failed to save to ES")
			}
		}
		cache.SetStringWithExpiration(constants.CACHE_PREFIX + summaryUrl, string(content), constants.TESTRESULTS_EXPIRATION)
		utils.SpectreLog.Debug("Leaving GetPerfScoreFromDetails()")
		return int(score), nil, details.StatusCode, details.StatusText
	} else {
		utils.SpectreLog.Debugln("Successful Runs not found!!")
		utils.SpectreLog.Warn("Successful Runs not found for: " + summaryUrl)
		utils.SpectreLog.Debug("Leaving GetPerfScoreFromDetails()")
		return 0, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text), http.StatusBadRequest, details.Data.Runs.One.FirstView.ConsoleLog[0].Text
	}
}

func getAverageScore(details models.JsonResults1) (float64) {
	total := float64(0)
	count := float64(0)
	if details.Data.Runs.One.FirstView.ScoreCombine != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreCombine)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreCache != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreCache)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreCdn != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreCdn)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreCompress != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreCompress)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreCookies != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreCookies)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreEtags != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreEtags)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreGzip != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreGzip)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreKeep_alive != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreKeep_alive)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreMinify != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreMinify)
		count += 1
	}

	if details.Data.Runs.One.FirstView.ScoreProgressiveJpeg != -1 {
		total += float64(details.Data.Runs.One.FirstView.ScoreProgressiveJpeg)
		count += 1
	}
	return total/count
}
