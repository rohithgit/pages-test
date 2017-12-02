package handlers

import (
	"testing"
	log "github.com/Sirupsen/logrus"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"fmt"
)

//func initMdbSession() (Session) {
//	session := mdb.NewMdbSession()
//	global.Session = session
//	// make sure the mongodb session is closed when we end
//	defer persist.EndSession(session)
//	// initialize the mongodb session
//	persist.InitSession(session)
//	return session
//}


func TestGetPerfScoreDetailsTime(t *testing.T) {
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer persist.EndSession(session)
	// initialize the mongodb session
	persist.InitSession(session)
	time.Sleep(time.Second*10)
	resultsDB,  err := persist.NewTestResultDB()
	if err != nil {
		log.Errorf("Failed to init Test Result DB: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx1 := context.WithValue(context.TODO(), ctxDBKey, resultsDB)
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		log.Errorf("Failed to init Test Result ES: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx1 = context.WithValue(ctx1, CtxESKey, resultsES)
	cache, errC := helper.FetchRCacheFromContext(ctx)

	if errC != nil {
		log.Errorf("Failed to init Hazelcast Cache: %v", err)
		//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
		return
	}
	ctx1 = context.WithValue(ctx1, ctxCacheKey, cache)

	location := "Test"
	wpt := helper.WebPageTester{}
	scoreDetailResults, err1, status, statusText, lookupUrl := GetPerfScoreDetailsFromCache(ctx1, wpt, url, location)

	if err1 != nil  {
		log.Debug("Not found in Cache. Lookup in ES")
		utils.SpectreLog.Debugln("Not found in Cache. Getting details from ES")
		scoreDetailResults, err1, status, statusText, lookupUrl = GetPerfScoreDetailsFromES(ctx1, wpt, url, location )

		if err1 != nil || scoreDetailResults.Overallscore == ""  {
			testUrl := buildTestURL(url, location)
			lookupUrl, _, err, status = GetSummaryUrl(wpt, testUrl)
			statusText = "Live Test"
			if(err != nil) {
				log.Errorf("Failed to get summary URL %s for url ",  lookupUrl, testUrl)
				//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
				return
			}

			content, err, _ := wpt.GetContent(lookupUrl)
			if err == nil && content != nil {
				log.Errorf("Error getting content for lookup URL %s", lookupUrl)
				jsonResult, err, status, statusText := GetDetailsFromContent( content)
				utils.SpectreLog.Debugf("Error is %v, status %d and statusText %s", err, status, statusText)
				scoreDetailResults := getScoreDetails(jsonResult)
				if scoreDetailResults.Overallscore != "" {
					utils.SpectreLog.Debugf("Successfully retrieved score details")
				}
				SetLookupUrlInCache(ctx1, url, location, lookupUrl)
			} else {
				log.Errorf("Failed to response for Lookup URL %s ",  lookupUrl)
				//http.Error(w, "internal error - see logs", http.StatusInternalServerError)
				return
			}
		}

		utils.SpectreLog.Debugf("Got results from Cache")
	}

	utils.SpectreLog.Debugf("Got score results. status %s and statusText %s ", status, statusText)
}

func TestSetResultsInCache(t *testing.T) {
	url := "www.vta.org"
	lookupurl := "hello123"
	location := "Test"
	cache, _ := helper.FetchRCacheFromContext(ctx1)
	//cache
	SetLookupUrlInCache(ctx1, url, location, lookupurl)
	lookupurl,_ = cache.GetString("spectre-pageload::www.vta.org::Test")
	utils.SpectreLog.Debugln("Lookup url " +lookupurl)
	if lookupurl != "hello123" {
		t.Error("Expected hello123!")
	}
}

func TestGetPerfScoreDetailsWithUrl(t *testing.T) {
	scoreDetailResults, err, status, statusText, lookupUrl := GetPerfScoreDetails(ctx1, "www.facebook.com", "Test")
	utils.SpectreLog.Debugf("Error: %v\n", err)
	utils.SpectreLog.Debugf("Status: %d\n", status)
	utils.SpectreLog.Debugf("Status text: %s\n", statusText)
	utils.SpectreLog.Debugf("Lookup Url: %s\n", lookupUrl)

	if(status == http.StatusOK) {
		if scoreDetailResults.Overallscore == "" || scoreDetailResults.Overallscore == "0" {
			t.Error("Expected scoreDetailResult with non zero or non empty overallscore")
		}
	} else if status == 101 {
		if lookupUrl == "" {
			t.Error("Expected non empty lookup url for live test")
		}
	}
}

func TestErrorJsonResponse(t *testing.T) {
	wpt := helper.WebPageTester{}
	content, _, _ :=  wpt.GetContent("http://128.107.18.61/jsonResult.php?test=161101_FZ_7V")

	if _, err, status, statusText := GetDetailsFromContent(content); err == nil {
		fmt.Printf("Status is %d", status)
		fmt.Printf("Status text is %s", statusText)
	} else {
		fmt.Printf("Error is %v", err)
	}
}

func TestGetPerfScoreDetailsWithLookUrl(t *testing.T) {
	scoreDetailResults, err, status, statusText := GetPerfScoreDetailsWithLookUrl(ctx1, helper.WebPageTester{},"http://spectre-webpagetest.infra-01.cisco-appanalytics.com/jsonResult.php?test=160930_GA_Q")
	utils.SpectreLog.Debugf("Error: %v\n", err)
	utils.SpectreLog.Debugf("Status: %d\n", status)
	utils.SpectreLog.Debugf("Status text: %s\n", statusText)
	//utils.SpectreLog.Debugf("Lookup Url: %s\n", lookupUrl)

	if(status == http.StatusOK) {
		if scoreDetailResults.Overallscore == "" || scoreDetailResults.Overallscore == "0" {
			t.Error("Expected scoreDetailResult with non zero or non empty overallscore")
		}
	} else if status != 101 {
		t.Error("Expected status 200 or 101")
	}
}
