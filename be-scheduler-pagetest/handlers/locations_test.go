package handlers

import (
	"testing"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
)

func TestGetLocations(t *testing.T) {
	result, err, status := GetLocations(FakeTester{}, "http://128.107.18.61/getLocations.php?f=json")
	utils.SpectreLog.Debugln("Result Status : %d \n", result.StatusCode)
	utils.SpectreLog.Debugln("Status : %d \n", status)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
}

func TestWPTGetLocations(t *testing.T) {
	result, err, status := GetLocations(FakeTester{}, "http://www.webpagetest.org/getLocations.php?f=json")
	utils.SpectreLog.Debugln(len(result.Test))
	utils.SpectreLog.Debugln("Status : %d \n", status)

	if err != nil {
		utils.SpectreLog.Debugln("Error for get details: %s ", err)
	}
}
