package handlers

import (
	"net/http"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// Page Load API returns summary URL. Json URL is used to lookup detailed results in LookupAvailability API
func Availability(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering AvailabilityHandler()")
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
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result DB: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(context.TODO(), CtxESKey, resultsES)
	_, errC := helper.FetchRCacheFromContext(ctx)
	if errC != nil {
		utils.SpectreLog.Errorf("Failed to init Redis Cache: %v", err)
		http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	// ctx = context.WithValue(ctx, ctxCacheKey, cache)
	rawurl := r.URL.Query().Get("url")
	//userId := r.URL.Query().Get("userId")
	//customerId := r.URL.Query().Get("customerId")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	interval := r.URL.Query().Get("interval")
	location := r.URL.Query().Get("location")
	if location == "" {
		location = constants.DEFAULT_LOCATION
	}

	totalRuns, pageLoadTime, availability, err1, status, statusText := GetAvailability(helper.WebPageTester{}, ctx, url, location, interval)

	if(err1 != nil) {
		utils.PrintError(w, status, err1.Error())
		return
	}
	results := new(models.Result)
	results.TestStatusCode = status
	results.StatusText = statusText
	results.PageLoadTime = pageLoadTime
	results.Availability = availability
	results.Runs = totalRuns
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(status)
}

func GetAvailability(wpt helper.Tester, ctx context.Context, url, location, interval string) (int, int, float64, error, int, string) {
	es, ok := FetchESFromContext(ctx)

	if !ok {
		return 0, 0, 0, errors.New("ES not found in context."), 500, "ES not found in context."
	}
	results, err := es.GetUrl(url, location, interval)
	if err != nil && results == nil {
		return 0, 0, 0, errors.New("internal error - see logs."), 500, "internal error - see logs."
	}
	totalRuns := len(results)
	if err != nil && totalRuns == 0 {
		//New test/location combination not in range
		/*pageLoadTime, err1, status, statusText, _ := GetPageLoadFromCache(wpt, ctx, url, location, interval)
		if pageLoadTime > 0 {
			return pageLoadTime, 1, nil, status, statusText
		}*/
		utils.SpectreLog.Error("Could not retrieve results from ElasticSearch.")
		return 0, 0, 0, errors.New("Error retrieving data."),  500,"Error retrieving data."

	}
	successfulRuns := 0
	totalPageLoad := 0
	for _, result := range results {
		if result.SuccessfulRuns > 0 {
			successfulRuns++
			totalPageLoad += result.Pageload.Loadtime
		}
	}
	if successfulRuns > 0 {
		return totalRuns, totalPageLoad / successfulRuns, float64(successfulRuns)/ float64(totalRuns), nil, 200, "Success"
	}
	return totalRuns, 0, 0, nil, 400, "No successful runs found."
}
