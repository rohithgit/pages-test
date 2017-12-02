// Package global contains global variables used in the other packages within the service.
package global

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
)

// Global variables
var (
	Options *models.Options
	Testing = false
	Session mdb.ISession
	IsSchedulerBusy bool
)

// init the global options
func init() {
	options, err := models.InitOptions()
	if err != nil {
		utils.SpectreLog.Fatal("Options init errored: ", err.Error())
	}
	Options = options
}
