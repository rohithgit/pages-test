// Package core contains code used to run the template service. Instead of
// having the code to run the service in main.go, the code can be moved into
// this 'core' package and then called from main.go. This keeps main.go clean and simple.
package core

import (
	//"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/controller"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/persist"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	//"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/scheduler"
	//"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/scheduler"
	//"fmt"
	//"runtime"
	//"sync"
	// "fmt"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/scheduler"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	//"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/controller"
	"sync"
	"runtime"
)

// Realmain starts the server and processes the REST calls
func Realmain() error {
	session := mdb.NewMdbSession()
	global.Session = session
	// make sure the mongodb session is closed when we end
	defer persist.EndSession(session)
	// initialize the mongodb session
	persist.InitSession(session)
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	wg.Add(2)

	utils.SpectreLog.Debugln("Starting Go Routines")
	go func() {
		defer wg.Done()

		if err := controller.TopoController(); err != nil {
			utils.SpectreLog.Errorln("Cannot start topocontroller", err)
			return
		}
	}()

	go func() {
		scheduler.StartScheduler()
	}()
	utils.SpectreLog.Debugln("Waiting To Finish")
	wg.Wait()
	utils.SpectreLog.Debugln("\nTerminating Program")
	return nil
}
