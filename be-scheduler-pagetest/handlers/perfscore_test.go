package handlers

import (
	"testing"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
)

func TestGetPerfScoreFromCache(t *testing.T) {
	pageLoad, err, status, statusText, lookupUrl := GetPerfScoreFromCache(FakeTester{}, ctx, "https://ui-portalserver.sjdev2.cloud.cisco.com", "")
	utils.SpectreLog.Debugln("Status : %d \n", status)
	utils.SpectreLog.Debugln("StatusText : %s \n", statusText)
	utils.SpectreLog.Debugln("Lookup Url : %s \n", lookupUrl)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}

	utils.SpectreLog.Debugln(" PerfScore : %d ", pageLoad)
}

func TestLookupPerfScoreFromCache(t *testing.T) {
	content, err, status := FakeTester{}.GetContent("http://128.107.18.61/jsonResult.php?test=160629_BY_4")
	perfscore, err, status, statusText := GetPerfScoreFromDetails(ctx,
		"http://128.107.18.61/jsonResult.php?test=160629_H7_E", content, status, true)
	utils.SpectreLog.Debugln("Performance Score : %d \n", perfscore)
	utils.SpectreLog.Debugln("Status code : %d \n", status)
	utils.SpectreLog.Debugln("Status text : %s \n", statusText)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
}
