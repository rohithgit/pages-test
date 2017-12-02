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
	"math"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"fmt"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// API for domain and page load returns summary URL. Json URL is used to lookup  results for domain load and page size
func DomainLoadTime(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}

	utils.SpectreLog.Debugln("Token is Valid")

	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	ctx := helper.CreateContext()
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)

	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}
	appId := r.URL.Query().Get("applicationId")
	utils.SpectreLog.Debugln("App Id is ", appId)

	//TODO get access token
	// if access token not there get UserID and customerID from query params
	//add user and customer ID to DB
	found := false
	isNew := false
	wpt := helper.WebPageTester{}
	content, err1, status, statusText, lookupUrl := GetDomainFromCache(ctx, url, location)

	var domainResults []models.DomainResult
	if(content != nil) {
		domainResults, err, statusText, status =  UnMarshalDomainFromContent(content)
		if(domainResults != nil) {
			if isNew {
				SaveContentToES(ctx, lookupUrl, content, domainResults, nil, nil)
				if err == nil && status == 200 {
					cache.SetStringWithExpiration(constants.CACHE_PREFIX + url, string(content), constants.TESTRESULTS_EXPIRATION)
				}
			}
			if err := json.NewEncoder(w).Encode(domainResults); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		found = true
	}
	if !found {
		result, err := getLatestFromDB(ctx, url, lookupUrl, location, "")
		if err == nil && len(result.Domains) > 0 {
			if err := json.NewEncoder(w).Encode(result.Domains); err != nil {
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
		_, err, statusText, status = GetLoadTimePerDomain(helper.WebPageTester{}, lookupUrl)
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
	//w.WriteHeader(status)
}

// API to lookup load time per domain and page size
func DomainLoadTimeResult(w http.ResponseWriter, r *http.Request) {
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

	ctx := helper.CreateContext()
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}

	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx = context.WithValue(ctx, CtxESKey, resultsES)

	var domainResults []models.DomainResult
	var err1 error
	var statusText string
	var status int

	isNew := false
	content, err, status, statusText, _ := LookupContentFromCache(ctx, lookupurl)
	if err == nil && content != nil {
		domainResults, err1, statusText, status = UnMarshalDomainFromContent(content)
		if(err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
		if isNew {
			SaveContentToES(ctx, lookupurl, content, domainResults, nil, nil)
			if err == nil && status == 200 {
				cache.SetStringWithExpiration(constants.CACHE_PREFIX + lookupurl, string(content), constants.TESTRESULTS_EXPIRATION)
			}
		}
	} else {

		result, err := getLatestFromDB(ctx, "", lookupurl, "", "")
		if err == nil && len(result.Domains) > 0 {
			if err := json.NewEncoder(w).Encode(result.Domains); err != nil {
				utils.PrintError(w, 500, "Error processing data.")
			}
			return
		}
		domainResults, err1, statusText, status = GetLoadTimePerDomain(helper.WebPageTester{}, lookupurl)
		if (err1 != nil) {
			utils.PrintError(w, status, err1.Error())
			return
		}
	}

	utils.SpectreLog.Debugln("Domain Results found: ", len(domainResults))
	utils.SpectreLog.Debugln("Status: %d\n ", status)
	utils.SpectreLog.Debugln("Lookup Url: %s\n ", lookupurl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if(domainResults == nil) {
		results := new(models.Result)
		results.TestStatusCode = status
		results.LookupUrl = lookupurl
		results.StatusText = statusText

		if err := json.NewEncoder(w).Encode(results); err != nil {
			utils.PrintError(w, 500, "Error processing data.")
		}
		return
	} else {

		if err := json.NewEncoder(w).Encode(domainResults); err != nil {
			utils.PrintError(w, 500, "Error processing data.")
			return
		}
	}
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}

func GetDomainFromCache(ctx context.Context, url, location string) ([]byte, error, int, string, string) {

	utils.SpectreLog.Debugln("Get domain from cache")
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init redis Cache: %v", errC)
	}

	lookupUrlFromCache, err := cache.GetString(constants.CACHE_PREFIX + url + "::" + location)
	if err == nil && lookupUrlFromCache != "" {
		fmt.Println("Lookup URL found in cache", lookupUrlFromCache)
		return LookupContentFromCache(ctx, lookupUrlFromCache)
	}

	fmt.Printf("Url is %s \n", url)
	fmt.Println(err)
	return nil, err, http.StatusNotFound, "Not found in cache", url
}

func LookupContentFromCache(ctx context.Context, url string) ([]byte, error, int, string, string) {
	cache, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		//return 0, errors.New("Unable to find cache in context."), http.StatusInternalServerError, "", "", false

	}

	var content []byte
	var status int
	var statusText string
	if results, _ := cache.GetString(constants.CACHE_PREFIX + url); results == "" {
		content, err, status := helper.WebPageTester{}.GetContent(url)
		return content, err, status, "", url
	} else {

		content = []byte(results)
		status = 200
		statusText = "Cached Response"
	}
	return content, nil, status, statusText, url
}

func UnMarshalDomainFromContent(content []byte) ([]models.DomainResult, error, string, int) {
	err, testingJsonResponse := UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugln("Unmarshalling normal response .... ")
		err, details := UnmarshalNormalResponse(content)
		if err != nil {
			return nil, err, "Error unmarshalling normal response", http.StatusInternalServerError
		}
		return ParseDomainInfo(details)
	} else {
		utils.SpectreLog.Debugln("Unmarshalling pending test response .... ")
		return nil, err, testingJsonResponse.Data.StatusText, testingJsonResponse.Data.StatusCode
	}
}

func ParseDomainInfo(details models.JsonResults1) ([]models.DomainResult, error, string, int) {
	if details.Data.SuccessfulFVRuns != 0 {
		results := make([]models.DomainResult, 0)
		domains := details.Data.Runs.One.FirstView.Domains
		totalLoadTime := 0
		for k, val := range domains {
			for _, val1 := range details.Data.Runs.One.FirstView.Requests {
				if val1.Host == k {
					foundDuplicate := false
					for index, val2 := range results {
						if val2.DomainName == val1.Host {
							time2, _ := strconv.Atoi(val1.LoadMs)
							results[index].LoadTime = val2.LoadTime + time2
							foundDuplicate = true
						}
					}

					if (!foundDuplicate) {
						var domain models.DomainResult
						domain.DomainName = val1.Host
						domain.LoadTime, _ = strconv.Atoi(val1.LoadMs)
						domain.PageSize = val.Bytes
						domain.Requests = val.Requests
						results = append(results, domain)
					}
					loadtime, _ := strconv.Atoi(val1.LoadMs)
					totalLoadTime = totalLoadTime + loadtime
				}
			}
		}
		//Calculate percentage TODO move to other calcs
		for index, val := range results {
			var result float64
			result = float64(val.LoadTime)/float64(totalLoadTime)
			results[index].Percent = RoundDown(result*100, 2)
		}

		return results, nil, details.StatusText, http.StatusOK
	} else {
		utils.SpectreLog.Debugln("Successful Runs not found!!")
		return nil, errors.New(details.Data.Runs.One.FirstView.ConsoleLog[0].Text),
		"Successful runs not found",
		http.StatusInternalServerError
	}
}

func GetLoadTimePerDomain(wpt helper.Tester, summaryUrl string) ([]models.DomainResult, error, string, int) {
	content, err, status := wpt.GetContent(summaryUrl)
	if err != nil {
		return nil, err, "Error getting content from summary url status", status
	}

	utils.SpectreLog.Debugln("Getting load time per domain from content: ", summaryUrl)
	return UnMarshalDomainFromContent(content)
}
