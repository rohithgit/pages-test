package persist

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"

	"testing"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"os"
	"time"
)

func TestMain(m *testing.M) {
	//data.CreateTransactions()
	global.Options.Environment = "test"
	global.Options.WptServer = "128.107.18.61"
	global.Options.TopoUrl = "spectre-topocrudsvc.service.consul:31030"
	global.Options.TopoPath = "/api/v1/topo/apps"
	global.Options.SchedulerEnabled = true
	global.Options.SchedulerInt = 10
	global.Options.Elastic = "127.0.0.1:9200"


	os.Exit(m.Run())
}
type dbKey string

func TestUrlInsert(t *testing.T) {
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer EndSession(session)
	// initialize the mongodb session
	InitSession(session)
	//urls := []string{"www.cisco.com", "www.msn.com"}
	time.Sleep(time.Second*5)

	urlDb,  err := NewUrlDB()
	urlToStore := new(models.SystemUrl)
	urlToStore.Url = "www.cisco.com"
	if( err == nil ) {
		id, err1 := urlDb.InsertUrl(urlToStore)

		if(err1 != nil) {
			utils.SpectreLog.Println("Error is ", err1)
		} else {
			utils.SpectreLog.Println("Successfully inserted with id ", id)
		}
	} else {
		utils.SpectreLog.Println("Error is ", err)
	}

	urls, err := urlDb.GetAllUrls()

	if err == nil {
		for _, url := range urls {
			utils.SpectreLog.Println("Url found ", url)
		}
	}

	//Remove Url
	err = urlDb.RemoveUrl(urlToStore.Url)

	if( err == nil) {
		utils.SpectreLog.Println("Deleted url entry:", urlToStore.Url)
	}
}

func TestPopulateDB(t *testing.T) {
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer EndSession(session)
	// initialize the mongodb session
	InitSession(session)
	//urls := []string{"www.cisco.com", "www.msn.com"}
	time.Sleep(time.Second*5)

	urlDb,  err := NewUrlDB()
	urlToStore := new(models.SystemUrl)
	urlToStore.Url = "www.cisco.com"
	if( err == nil ) {
		id, err1 := urlDb.InsertUrl(urlToStore)

		if(err1 != nil) {
			utils.SpectreLog.Println("Error is ", err1)
		} else {
			utils.SpectreLog.Println("Successfully inserted with id ", id)
		}
	} else {
		utils.SpectreLog.Println("Error is ", err)
	}

	urlToStore.Url = "www.msn.com"

	if( err == nil ) {
		id, err1 := urlDb.InsertUrl(urlToStore)

		if(err1 != nil) {
			utils.SpectreLog.Println("Error is ", err1)
		} else {
			utils.SpectreLog.Println("Successfully inserted with id ", id)
		}
	} else {
		utils.SpectreLog.Println("Error is ", err)
	}
	urls, err := urlDb.GetAllUrls()

	if err == nil {
		for _, url := range urls {
			utils.SpectreLog.Println("Url found ", url)
		}
	}
}
