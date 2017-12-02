/*
The topo service does discovery of public cloud elements
*/
package main

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/commands"
	_ "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/config"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/core"
	_ "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/version"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/config"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
)

// main entrypoint for the topo service
func main() {

	// Only log the warning severity or above.
	config.SetupConfig()
	// read command line flags and commands
	err := commands.InitCommands()
	if err == nil {
		// process REST calls
		err = core.Realmain()
	}

	if err != nil {
		utils.SpectreLog.Error(err.Error())
	}
	utils.SpectreLog.Println("Completed successfully!")


}
