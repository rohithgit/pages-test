package handlers

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"net/http"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"golang.org/x/net/context"
	"errors"
	"strings"
	"strconv"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// API for content type and page load returns summary URL. Json URL is used to lookup  results for content type load and page size
func ContentTypes(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_WRITE) {
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
		utils.SpectreLog.Errorf("Failed to init Redis Cache: %v", errC)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}
	//TODO get access token
	// if access token not there get UserID and customerID from query params
	//add user and customer ID to DB
	found := false
	wpt := helper.WebPageTester{}
	content, err1, status, statusText, lookupUrl, isNew := GetContentTypesFromCache(ctx, url, location)
	var contentTypes []models.ContentType
	if(content != nil) {
		contentTypes, err, statusText, status =  UnMarshalContentTypesFromContent(content)
		if(contentTypes != nil) {
			if isNew {
				SaveContentToES(ctx, lookupUrl, content, nil, contentTypes, nil)
				if err == nil && status == 200 {
					cache.SetStringWithExpiration(constants.CACHE_PREFIX + url, string(content), constants.TESTRESULTS_EXPIRATION)
				}
			}
			if err := json.NewEncoder(w).Encode(contentTypes); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		found = true
	}
	if !found {
		result, err := getLatestFromDB(ctx, url, lookupUrl, location, "")
		if err == nil && len(result.ContentTypes) > 0 {
			if err := json.NewEncoder(w).Encode(result.ContentTypes); err != nil {
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
		_, err, statusText, status = GetLoadTimePerContentType(helper.WebPageTester{}, lookupUrl)
	}

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
func ContentTypesResult(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}

	utils.SpectreLog.Debugln("Token is Valid")


	//TODO get access token
	// if access token not there get UserID and customerID from query params
	//add user and customer ID to DB
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
		utils.SpectreLog.Errorf("Failed to init Redis Cache: %v", errC)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
	}

	var contentTypes []models.ContentType
	var err1 error
	var statusText string
	var status int

	content, err, status, statusText, _, isNew := LookupContentTypesFromCache(ctx, lookupurl)
	if err == nil && content != nil {
		contentTypes, err1, statusText, status = UnMarshalContentTypesFromContent(content)
		if(err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
		if isNew {
			SaveContentToES(ctx, lookupurl, content, nil, contentTypes, nil)
			if err == nil && status == 200 {
				cache.SetStringWithExpiration(constants.CACHE_PREFIX + lookupurl, string(content), constants.TESTRESULTS_EXPIRATION)
			}
		}
	} else {

		result, err := getLatestFromDB(ctx, "", lookupurl, "", "")
		if err == nil && len(result.ContentTypes) > 0 {
			if err := json.NewEncoder(w).Encode(result.ContentTypes); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		contentTypes, err1, statusText, status = GetLoadTimePerContentType(helper.WebPageTester{}, lookupurl)
		if (err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if(contentTypes == nil) {
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
		if err := json.NewEncoder(w).Encode(contentTypes); err != nil {
			utils.PrintError(w, 500, "Error processing data.")
			return
		}
	}
}

func GetContentTypesFromCache(ctx context.Context, url, location string) ([]byte, error, int, string, string, bool) {

	utils.SpectreLog.Debugln("Get content types from cache")
	cache, err := helper.FetchRCacheFromContext(ctx)
	if (err != nil) {
		//return 0, errors.New("Unable to find cache in context."), http.StatusInternalServerError, "", "", false
	}
	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		utils.SpectreLog.Debugln("Lookup URL found in cache", lookupUrlFromCache)
		return LookupContentTypesFromCache(ctx, lookupUrlFromCache)
	}

	utils.SpectreLog.Debugln("Url is %s \n", url)
	return nil, err, http.StatusNotFound, "Not found in cache", url, false
}

func LookupContentTypesFromCache(ctx context.Context, url string) ([]byte, error, int, string, string, bool) {
	cache, err := helper.FetchRCacheFromContext(ctx)
	if (err != nil) {
		//return 0, errors.New("Unable to find cache in context."), http.StatusInternalServerError, "", "", false
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

func UnMarshalContentTypesFromContent(content []byte) ([]models.ContentType, error, string, int) {
	err, testingJsonResponse := UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		if err != nil {
			return nil, err, "Error unmarshalling normal response", http.StatusInternalServerError
		}
		return ParseContentTypesInfo(details)
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return nil, err, testingJsonResponse.Data.StatusText, testingJsonResponse.Data.StatusCode
	}
}

func ParseContentTypesInfo(details models.JsonResults1) ([]models.ContentType, error, string, int) {
	if details.Data.SuccessfulFVRuns != 0 {
		utils.SpectreLog.Debugln("Successful Runs found!!")
		contentTypes := make([]models.ContentType, 0)
		breakdown := details.Data.Runs.One.FirstView.Breakdown.(map[string]interface{})
		contentTypeMap := make(map[string]models.ContentType)
		totalLoadTime := 0
		for _, val := range details.Data.Runs.One.FirstView.Requests {
			loadtime, err := strconv.Atoi(val.LoadMs)
			if err != nil {
				utils.SpectreLog.Errorf("Could not convert string %v to int. \n", val.LoadMs)
			} else {
				totalLoadTime += loadtime
				switch {//TODO revisit all contentType strings
				case strings.Contains(val.ContentType, "html"):
					contentType := contentTypeMap["html"]
					contentType.Name = "HTML"
					contentType.LoadTime += loadtime
					contentTypeMap["html"] = contentType
				case strings.Contains(val.ContentType, "css"):
					contentType := contentTypeMap["css"]
					contentType.Name = "CSS"
					contentType.LoadTime += loadtime
					contentTypeMap["css"] = contentType
				case strings.Contains(val.ContentType, "javascript"):
					contentType := contentTypeMap["js"]
					contentType.Name = "JS"
					contentType.LoadTime += loadtime
					contentTypeMap["js"] = contentType
				case strings.Contains(val.ContentType, "flash"):
					contentType := contentTypeMap["flash"]
					contentType.Name = "Flash"
					contentType.LoadTime += loadtime
					contentTypeMap["flash"] = contentType
				case strings.Contains(val.ContentType, "font"):
					contentType := contentTypeMap["font"]
					contentType.Name = "Font"
					contentType.LoadTime += loadtime
					contentTypeMap["font"] = contentType
				case strings.Contains(val.ContentType, "image"):
					contentType := contentTypeMap["image"]
					contentType.Name = "Image"
					contentType.LoadTime += loadtime
					contentTypeMap["image"] = contentType
				default:
					contentType := contentTypeMap["other"]
					contentType.Name = "Other"
					contentType.LoadTime += loadtime
					contentTypeMap["other"] = contentType
				}
			}
		}
		for key, val := range breakdown {
			valMap := val.(map[string]interface{})
			contentType := contentTypeMap[key]
			contentType.Size = int(valMap["bytes"].(float64))
			if totalLoadTime > 0 {
				contentType.Percent = float64(contentType.LoadTime) / float64(totalLoadTime)
			}
			if contentType.Size > 0 {
				contentTypes = append(contentTypes, contentType)
			}
		}
		return contentTypes, nil, details.StatusText, http.StatusOK
	} else {
		utils.SpectreLog.Debugln("Successful Runs not found!!")
		return nil, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text),
		"Successful runs not found",
		http.StatusInternalServerError
	}
}

func GetLoadTimePerContentType(wpt helper.WebPageTester, summaryUrl string) ([]models.ContentType, error, string, int) {
	content, err, status := wpt.GetContent(summaryUrl)
	if err != nil {
		return nil, err, "Error getting content from summary url status", status
	}

	utils.SpectreLog.Debugln("Getting load time per content types from content: ", summaryUrl)
	return UnMarshalContentTypesFromContent(content)
}
