package handlers

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"encoding/json"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
)

func SuccessfulRuns(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering SuccessfulRuns()")
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}

	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}

	resultsES,  err := persist.NewTestResultES()
	rawurl := r.URL.Query().Get("url")

	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	interval := r.URL.Query().Get("interval")
	location := r.URL.Query().Get("location")
	count, err := resultsES.GetSuccessfulRuns(url, location, interval)

	if(err != nil) {
		utils.PrintError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(count); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(200)
}

func SlowestRun(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}
	if !security.ValidateRequest(w, r, constants.SCOPE_READ) {
		return
	}
	resultsES,  err := persist.NewTestResultES()
	rawurl := r.URL.Query().Get("url")
	url, parseerr := CanonicalizeUrl(rawurl)
	if parseerr != nil {
		utils.PrintError(w, http.StatusBadRequest, parseerr.Error())
		return
	}
	interval := r.URL.Query().Get("interval")
	location := r.URL.Query().Get("location")
	run, err := resultsES.GetSlowestRun(url, location, interval)

	if(err != nil) {
		utils.PrintError(w, 500, err.Error())
		return
	}
	//datetime := run.Runtime.Format("January 2, 06 3:04:05")
	slowestRun := models.SlowestRun{Loadtime: run.Loadtime, Runtime: run.Runtime.Unix(), Location: run.Location}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(slowestRun); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
	w.WriteHeader(200)
}
