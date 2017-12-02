package scheduler

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"github.com/jasonlvhit/gocron"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"golang.org/x/net/context"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/handlers"
	"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"fmt"
	"strings"
	"strconv"
)


//Get list of URLs from DB
//scheduler task to start test for each URL
// Get test results from the test run
func SubmitBatch(wpt helper.Tester, urls []string , testParams string , serverUrl string) (map[string]string, error) {
	utils.SpectreLog.Debug("Entering Scheduler.SubmitBatch()")
	var testIds map[string]string
	testIds = make(map[string]string)

	for _, url := range urls {
		testUrl :=  global.Options.WptServer + "/runtest.php?"+ "url=" + url + "&ignoreSSL=1&f=json&fvonly=1"
		utils.SpectreLog.Println("Getting lookup URL for " + testUrl)
		jsonUrl, testId, err, status := handlers.GetSummaryUrl(wpt, testUrl)

		if err == nil {
			fmt.Println("Returned lookupUrl:" + jsonUrl + " for testUrl:" + testUrl)
			fmt.Println("Status for get lookupUrl:", status)

			testIds[testId] = url
		} else {
			utils.SpectreLog.Println("Error while getting lookup url", err, testUrl)
		}
		//time.Sleep(time.Second * 30)
	}
	utils.SpectreLog.Debug("Leaving Scheduler.SubmitBatch()")
	return testIds, nil
}

func GetUrlsFromDB(tenantId string) ([]string) {
	utils.SpectreLog.Debug("Entering Scheduler.GetUrlsFromDB()")
	urlsDB,  err := persist.NewUrlDB()
	var urls []string
	if err == nil {
		utils.SpectreLog.Println("Getting results from DB")
		urls, err = urlsDB.GetAllUrls()
		utils.SpectreLog.Println(err)

	} else {
		utils.SpectreLog.Println("Error fetching db from context")
		urls =[]string{"http://www.msn.com", "http://www.google.com"}
	}
	uniqueUrls := removeDuplicates(urls)
	utils.SpectreLog.Debug("Leaving Scheduler.GetUrlsFromDB()")
	return uniqueUrls
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			utils.SpectreLog.Println("URLs for which results are to be collected", elements[v])
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func StartScheduler() {
	if(global.Options.SchedulerEnabled) {
		time.Sleep(time.Second*10)

		utils.SpectreLog.Println("Starting scheduler")
		s := gocron.NewScheduler()
		utils.SpectreLog.Println("Scheduler Interval in minutes: ", global.Options.SchedulerInt)
		s.Every(uint64(global.Options.SchedulerInt)).Minutes().Do(GetAllTestResults)
		<-s.Start()
	} else {
		utils.SpectreLog.Println("Scheduler not enabled")
	}
}

//func StopScheduler() {
//	<- s.Start()
//}

func GetAllTestResults() {
	utils.SpectreLog.Debug("Scheduler triggered")
	if(global.IsSchedulerBusy) {
		utils.SpectreLog.Warn("Scheduler is busy... Next batch will be started after scheduler completes current batch")
		return
	}

	global.IsSchedulerBusy = true

	//TODO move this to start scheduler
	wpt := helper.WebPageTester{}
	utils.SpectreLog.Println("Starting Batch Collection.......")
	utils.SpectreLog.Debug("Starting Batch Collection.......")
	//Get List of URLs from DB
	urls := GetUrlsFromDB("customer1")
	//Submit Batch
	testIds, err := SubmitBatch(wpt, urls, "", "")
	if(err != nil) {
		utils.SpectreLog.Debugln("Error submitting batch for %d urls \n", len(urls))
		return
	}

	//check result after 5minutes
	// Timer does not work in gocron scheduler. Have to use sleep


	checkBatchResults(wpt, testIds)
}

func initializeSystemUrls() ([]string) {
	var systemUrls []string
	systemUrls = make([]string, 0)
	systemUrls = append(systemUrls, "www.cisco.com")
	systemUrls = append(systemUrls, "www.msn.com")
	systemUrls = append(systemUrls, "www.google.com")
	systemUrls = append(systemUrls, "www.facebook.com")
	return systemUrls
}

func insertSystemUrls(urls []string)  {
	urlDb,  err := persist.NewUrlDB()

	if err == nil {
		for _, url := range urls {
			utils.SpectreLog.Debugf("Inserting system Url in Db: %s ", url)
			modelToStore := new(models.SystemUrl)
			modelToStore.Url = url
			id, err1 := urlDb.InsertUrl(modelToStore)

			if (err1 != nil) {
				utils.SpectreLog.Warnf("Url %s not inserted due to error:%s ", url, err1)
			} else {
				utils.SpectreLog.Debugf("Successfully inserted url: %s with id ", url, id)
			}
		}
	} else {
		utils.SpectreLog.Errorf("Error getting handle to DB %s", err)
	}
}

func triggerTimer(wpt helper.Tester, testIds map[string]string) {
	timeChan := time.NewTimer(time.Minute*2).C
	for {
		select {
		case <-timeChan:
			utils.SpectreLog.Debug("Timer expired. Checking test results from WPT....")
			checkBatchResults(wpt, testIds)
		}
	}
}

func checkBatchResults(wpt helper.Tester, testIds map[string]string)  {
	utils.SpectreLog.Debug("Entering Scheduler.checkBatchResults()")
	utils.SpectreLog.Debug("Checking Results.......")
	//Get Batch Results

	for {
		resultMap, pendingTests, err := GetBatchResults(wpt, testIds, "")
		if (err != nil) {
			utils.SpectreLog.Warnf("Error getting batch results. Please check connection to WebPageTest server: %s\n", err)
			return
		}

		//Save Results to DB
		utils.SpectreLog.Infof("Retrieved results for the batch %d urls. Saving results retrieved.\n", len(resultMap))

		if (len(resultMap) > 0) {
			SaveResultMapToES(resultMap)
		}

		if(len(pendingTests) > 0 ) {
			utils.SpectreLog.Infof("Found %d pending tests", len(pendingTests))
			testIds = pendingTests
		} else {
			utils.SpectreLog.Info("No pending test results found. Enabling scheduler for next batch.....")
			global.IsSchedulerBusy = false
			break
		}
		time.Sleep(time.Minute)
	}
	utils.SpectreLog.Debug("Leaving Scheduler.checkBatchResults()")
}

func GetBatchResults(wpt helper.Tester, testIds map[string]string,  serverUrl string) (map[string]models.JsonResults1, map[string]string, error) {
	utils.SpectreLog.Debug("Entering Scheduler.GetBatchResults()")
	resultMap := map[string]models.JsonResults1{}
	pendingTestIds := map[string]string{}
	for testId, _ := range testIds {
		jsonUrl := global.Options.WptServer + "/jsonResult.php?test=" + testId
		details, err, statusCode, statusText := GetResultsFromDetails(wpt, jsonUrl)

		if(err == nil) {
			if statusCode == 200 {
				utils.SpectreLog.Printf("Got test results: Status code: %d and status: %s for lookupUrl %s for url %s\n", statusCode, statusText, jsonUrl, testIds[testId])
				resultMap[testId] = details
			} else if statusCode < 200 {
				//TODO add testId to pending testIds
				utils.SpectreLog.Infof("Pending test for testId: %s", testId)
				utils.SpectreLog.Infof("Status code: %d and Text: %s for lookupUrl: %s\n", statusCode, statusText, jsonUrl)

				if(statusCode == 101 && strings.Contains(statusText, "Waiting behind")) {
					checkQueueLenandWarn(statusText)
				}
				pendingTestIds[testId] = testIds[testId]
			} else {
				utils.SpectreLog.Warnf("Error retrieving resuls from WPT for testId %s", testId)
				utils.SpectreLog.Warnf("Status code: %d and status: %s for lookupUrl: %s\n", statusCode, statusText, jsonUrl)
			}
		} else {
			utils.SpectreLog.Warnf("Failed to get results for testId: %s due to error %s", testId, err)
		}
	}

	utils.SpectreLog.Debugf("Collected results for %d urls", len(resultMap))
	utils.SpectreLog.Debug("Leaving Scheduler.GetBatchResults()")
	return resultMap, pendingTestIds, nil
}

func checkQueueLenandWarn(statusText string) {
	lastBin := strings.LastIndex(statusText, "behind ")
	firstIndex := strings.Index(statusText, "other")

	queuelenStr := statusText[lastBin + len("behind") + 1:firstIndex - 1]
	utils.SpectreLog.Debugf("Web Page test queue len: %s", queuelenStr)
	queuelen, err := strconv.Atoi(queuelenStr)
	if (err == nil) {
		if (queuelen > 30) {
			utils.SpectreLog.Warn("!!!!Number of Pending tests in queue for Spectre Web Endpoints exceeded threshold!!!!")
		}
	}

}

func GetResultsFromDetails(wpt helper.Tester, lookupurl string) (models.JsonResults1, error, int, string) {
	utils.SpectreLog.Debug("Entering Scheduler.GetResultsFromDetails()")
	content, err, status := wpt.GetContent(lookupurl)
	if err != nil {
		return models.JsonResults1{}, err, status, "Error Getting Content from details Url"
	}

	err, testingJsonResponse := handlers.UnmarshalTestingResponse(content)
	if err != nil {
		utils.SpectreLog.Debugf("Unmarshalling normal response for lookupurl:%s", lookupurl)
		utils.SpectreLog.Println("Unmarshalling normal response .... ")
		err, details := handlers.UnmarshalNormalResponse(content)
		utils.SpectreLog.Debug("Leaving Scheduler.GetResultsFromDetails()")
		return details, err, details.StatusCode, details.StatusText
	} else {
		utils.SpectreLog.Println("Unmarshalling pending test response for lookupurl:%s", lookupurl)
		utils.SpectreLog.Debug("Leaving Scheduler.GetResultsFromDetails()")
		return models.JsonResults1{}, err, testingJsonResponse.StatusCode, testingJsonResponse.StatusText
	}
}

func SaveResultMapToES(resultMap map[string]models.JsonResults1) {
	utils.SpectreLog.Debug("Entering Scheduler.SaveResultMapToES()")
	utils.SpectreLog.Debugln("Saving results to ES for %d urls \n", len(resultMap))
	resultsES,  err := persist.NewTestResultES()
	if err != nil {
		utils.SpectreLog.Errorf("Failed to init Test Result ES: %v", err)
		utils.SpectreLog.Debug("Leaving Scheduler.SaveResultMapToES()")
		return
	}
	ctx := context.WithValue(context.TODO(), handlers.CtxESKey, resultsES)
	//if ok {
	for testId, result := range resultMap {
		utils.SpectreLog.Debugf("Saving Result for Test ID from map: %s\n ", testId)
		if result.Data.SuccessfulFVRuns > 0 {
			_, err := handlers.SaveToES(ctx, global.Options.WptServer + "/jsonResult.php?test=" + testId, result, nil, nil, nil)
			if err != nil {
				utils.SpectreLog.Warnf("Failed to fetch ES from context: %s", err)
			}
		} else {
			utils.SpectreLog.Warnf("No successful runs found for testId %s", testId) //TODO save unsuccessful runs as well
		}
	}
	utils.SpectreLog.Debug("Leaving Scheduler.SaveResultMapToES()")
	return
}
