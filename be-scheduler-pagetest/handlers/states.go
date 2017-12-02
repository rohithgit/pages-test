package handlers

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"net/http"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"golang.org/x/net/context"
	"errors"
	"strconv"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// API for content type and page load returns summary URL. Json URL is used to lookup  results for content type load and page size
func States(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		utils.SpectreLog.Errorln("Token is not valid.")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return

	}
	utils.SpectreLog.Debugln("Token is Valid")

	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), CtxESKey, resultsES)

	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}

	found := false
	wpt := helper.WebPageTester{}
	content, err1, status, statusText, lookupUrl, isNew := GetStatesFromCache(ctx, url, location)
	var states []models.State
	if(content != nil) {
		states, err, statusText, status =  UnMarshalStatesFromContent(content)
		if(states != nil) {
			if isNew {
				SaveContentToES(ctx, lookupUrl, content, nil, nil, states)
				if err == nil && status == 200 {
					cache.SetStringWithExpiration(constants.CACHE_PREFIX + url, string(content), constants.TESTRESULTS_EXPIRATION)
				}
			}
			if err := json.NewEncoder(w).Encode(states); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		found = true
	}
	if !found {
		result, err := getLatestFromDB(ctx, url, lookupUrl, location, "")
		if err == nil && len(result.States) > 0 {
			if err := json.NewEncoder(w).Encode(result.States); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
	}

	//If not in Cache do real time test and return summary URL
	if !found {
		lookupUrl, _, err1, status = GetSummaryUrl(wpt, buildTestURL(url, location))
		if (err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
		err = cache.SetStringWithExpiration(constants.CACHE_PREFIX + url + "::" + location, lookupUrl, constants.PAGELOOKUP_EXPIRATION)
		_, err, statusText, status = GetLoadTimePerState(helper.WebPageTester{}, lookupUrl)
	}
	utils.SpectreLog.Debugln("Status: %d\n ", status)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupUrl)
	utils.SpectreLog.Debug("Retrieved summary URL: " + lookupUrl)
	results := new(models.Result)
	results.TestStatusCode = status
	results.LookupUrl = lookupUrl
	results.StatusText = statusText

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(status)
}

// API to lookup load time per content type and page size
func StatesResult(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		utils.SpectreLog.Errorln("Token is not valid.")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return

	}
	utils.SpectreLog.Debugln("Token is Valid")
	lookupurl := r.URL.Query().Get("lookupurl")

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), CtxESKey, resultsES)
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	var states []models.State
	var err1 error
	var statusText string
	var status int

	content, err, status, statusText, _, isNew := LookupStatesFromCache(ctx, lookupurl)
	if err == nil && content != nil {
		states, err1, statusText, status = UnMarshalStatesFromContent(content)
		if(err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
		if isNew {
			SaveContentToES(ctx, lookupurl, content, nil, nil, states)
			if err == nil && status == 200 {
				cache.SetStringWithExpiration(constants.CACHE_PREFIX + lookupurl, string(content), constants.TESTRESULTS_EXPIRATION)
			}
		}
	} else {

		result, err := getLatestFromDB(ctx, "", lookupurl, "", "")
		if err == nil && len(result.States) > 0 {
			if err := json.NewEncoder(w).Encode(result.States); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		states, err1, statusText, status = GetLoadTimePerState(helper.WebPageTester{}, lookupurl)
		if (err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
	}

	utils.SpectreLog.Debugln("Content Type Results found: ", len(states))
	utils.SpectreLog.Debugln("Status: %d\n ", status)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupurl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if(states == nil) {
		results := new(models.Result)
		results.TestStatusCode = status
		results.LookupUrl = lookupurl
		results.StatusText = statusText

		if err := json.NewEncoder(w).Encode(results); err != nil {
			utils.PrintError(w, 500, "Error processing data.")
			return
		}
	} else {
		// valid results
		if err := json.NewEncoder(w).Encode(states); err != nil {
			utils.PrintError(w, 500, "Error processing data.")
			return
		}
	}
}

func GetStatesFromCache(ctx context.Context, url, location string) ([]byte, error, int, string, string, bool) {

	utils.SpectreLog.Debugln("Get content types from cache")
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		utils.SpectreLog.Debugln("Lookup URL found in cache", lookupUrlFromCache)
		return LookupStatesFromCache(ctx, lookupUrlFromCache)
	}

	utils.SpectreLog.Debugln("Url is %s \n", url)
	utils.SpectreLog.Debugln(err)
	return nil, err, http.StatusNotFound, "Not found in cache", url, false
}

func LookupStatesFromCache(ctx context.Context, url string) ([]byte, error, int, string, string, bool) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
		// http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		// return
	}

	var content []byte
	var status int
	var statusText string
	if results, _ := cache.GetString(constants.CACHE_PREFIX + url); results == "" {
		content, err, status := helper.WebPageTester{}.GetContent(url)
		if err == nil && status == 200 {
			return content, nil, status, "", url, true
		}
		return nil, errors.New("Not found in Cache"), http.StatusNotFound, "Not found in Cache", url, false
	} else {
		content = []byte(results)
		status = 200
		statusText = "Cached Response"
	}
	return content, nil, status, statusText, url, false
}

func UnMarshalStatesFromContent(content []byte) ([]models.State, error, string, int) {
	err, testingJsonResponse := UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		utils.SpectreLog.Debugln(err)
		if err != nil {
			return nil, err, "Error unmarshalling normal response", http.StatusInternalServerError
		}
		return ParseStatesInfo(details)
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return nil, err, testingJsonResponse.Data.StatusText, testingJsonResponse.Data.StatusCode
	}
}

func ParseStatesInfo(details models.JsonResults1) ([]models.State, error, string, int) {
	if details.Data.SuccessfulFVRuns != 0 {
		utils.SpectreLog.Debugln("Successful Runs found!!")
		states := make([]models.State, 0)
		//stateMap := make(map[string]models.State)
		var ttfb, dns, connect, ssl, dl int
		var ok bool
		var strval string
		var floatval float64
		totalttfb := 0
		totaldns := 0
		totalconnect := 0
		totalssl := 0
		totaldl := 0
		for _, val := range details.Data.Runs.One.FirstView.Requests {
			floatval, ok = val.TtfbMs.(float64)
			if !ok {
				strval, ok = val.TtfbMs.(string)
				if ok {
					ttfb, _ = strconv.Atoi(strval)
				}
			} else {
				ttfb = int(floatval)
			}
			if ttfb > 0 {
				totalttfb += ttfb
			}
			floatval, ok = val.DnsMs.(float64)
			if !ok {
				strval, ok = val.DnsMs.(string)
				if ok {
					dns, _ = strconv.Atoi(strval)
				}
			} else {
				dns = int(floatval)
			}
			if dns > 0 {
				totaldns += dns
			}
			floatval, ok = val.ConnectMs.(float64)
			if !ok {
				strval, ok = val.ConnectMs.(string)
				if ok {
					connect, _ = strconv.Atoi(strval)
				}
			} else {
				connect = int(floatval)
			}
			if connect > 0 {
				totalconnect += connect
			}
			floatval, ok = val.SslMs.(float64)
			if !ok {
				strval, ok = val.SslMs.(string)
				if ok {
					ssl, _ = strconv.Atoi(strval)
				}
			} else {
				ssl = int(floatval)
			}
			if ssl > 0 {
				totalssl += ssl
			}
			floatval, ok = val.DownloadMs.(float64)
			if !ok {
				strval, ok = val.DownloadMs.(string)
				if ok {
					dl, _ = strconv.Atoi(strval)
				}
			} else {
				dl = int(floatval)
			}
			if dl > 0 {
				totaldl += dl
			}
		}
		totalLoadTime := totalttfb + totaldns + totalconnect + totalssl + totaldl
		if totalttfb > 0 {
			state := models.State{}
			state.Name = "Wait"
			state.TimeSpent = totalttfb
			state.Percent = float64(totalttfb) / float64(totalLoadTime)
			states = append(states, state)
		}
		if totaldns > 0 {
			state := models.State{}
			state.Name = "DNS"
			state.TimeSpent = totaldns
			state.Percent = float64(totaldns) / float64(totalLoadTime)
			states = append(states, state)
		}
		if totalconnect > 0 {
			state := models.State{}
			state.Name = "Connect"
			state.TimeSpent = totalconnect
			state.Percent = float64(totalconnect) / float64(totalLoadTime)
			states = append(states, state)
		}
		if totalssl > 0 {
			state := models.State{}
			state.Name = "SSL"
			state.TimeSpent = totalssl
			state.Percent = float64(totalssl) / float64(totalLoadTime)
			states = append(states, state)
		}
		if totaldl > 0 {
			state := models.State{}
			state.Name = "Download"
			state.TimeSpent = totaldl
			state.Percent = float64(totaldl) / float64(totalLoadTime)
			states = append(states, state)
		}
		utils.SpectreLog.Debugln(totalttfb, totaldns, totalconnect, totalssl, totaldl)
		return states, nil, details.StatusText, http.StatusOK
	} else {
		utils.SpectreLog.Debugln("Successful Runs not found!!")
		return nil, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text),
		"Successful runs not found",
		http.StatusInternalServerError
	}
}

func GetLoadTimePerState(wpt helper.WebPageTester, summaryUrl string) ([]models.State, error, string, int) {
	content, err, status := wpt.GetContent(summaryUrl)
	if err != nil {
		return nil, err, "Error getting content from summary url status", status
	}

	utils.SpectreLog.Debugln("Getting load time per content types from content: ", summaryUrl)
	return UnMarshalStatesFromContent(content)
}
