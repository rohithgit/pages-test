package handlers

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"net/http"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"encoding/json"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"reflect"
	"errors"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"fmt"
)

// Locations API returns location of  testers. Json URL is used to lookup detailed results in LookupPageLoad API
func Locations(w http.ResponseWriter, r *http.Request) {
	utils.SpectreLog.Debug("Entering locations()")
	if (r.Method != http.MethodGet) {
		utils.SpectreLog.Errorln("Unsupported Method: %v", r.Method)
		http.Error(w, "Unsupported Method:", http.StatusInternalServerError)
		return
	}	
	locationsUrl :=  global.Options.WptServer + "/getLocations.php?f=json"
	utils.SpectreLog.Debugln("Token is Valid")
	result, err, status := GetLocations(helper.WebPageTester{}, locationsUrl)

	if err != nil {
		utils.PrintError(w, status, "Internal Server Error: Cannot get locations of testing agents")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		utils.PrintError(w, 500, "Error processing data.")
		return
	}
}

func GetLocations(wpt helper.Tester, url string) (models.LocationsResponse, error, int) {
	content, err, status := wpt.GetContent(url)

	if err != nil {
		return models.LocationsResponse{}, err, status
	}

	err, result := UnmarshalLocationsResult(content)
	return result, err, status
}

func UnmarshalLocationsResult(content []byte) (error, models.LocationsResponse) {
	var details models.LocationsResult1
	if err := json.Unmarshal(content, &details); err != nil {
		return err, models.LocationsResponse{}
	}

	var response models.LocationsResponse
	response.StatusCode = details.StatusCode
	for k, _ := range details.Data {
		var result TestAgentStruct
		if err := json.Unmarshal(*details.Data[k], &result); err != nil {
			utils.SpectreLog.Debugln(err)
		}

		var testAgent models.TestAgentInfo
		testAgent.TestBrowser = result.Browser
		testAgent.TestLabel = result.Label
		testAgent.TestLocation =  result.Location
		response.Test = append(response.Test, testAgent)

		utils.SpectreLog.Debugln("---------Test Agent: ------ ", k)
		utils.SpectreLog.Debugln(testAgent)
		utils.SpectreLog.Debugln("--------------------------- ")
		utils.SpectreLog.Debugln()
	}

	utils.SpectreLog.Debugln("Number of test Agent", len(response.Test))
	return nil, response
}

type TestAgentStruct struct {
	Label         string
	Location      string
	Browser       string
	RelayServer   interface{}
	RelayLocation interface{}
	LabelShort    string
	Default       bool
	PendingTests  struct {
			      P1           int
			      P2           int
			      P3           int
			      P4           int
			      P5           int
			      P6           int
			      P7           int
			      P8           int
			      P9           int
			      Total        int
			      HighPriority int
			      LowPriority  int
			      Testing      int
			      Idle         int
		      }
}


func (s *TestAgentStruct) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}