// Package controller contains code for one or more controllers
package controller

import (
	"net/http"
	"strconv"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/router"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/security"
)

// TopoController handles the REST calls for the Template service
func TopoController() error {
	r, err := router.NewTopoRouter()
	if err != nil {
		return err
	}

	// Bind to a port and pass our router in
	// and start listening for requests
	port := ":" + strconv.Itoa(global.Options.Port)
	utils.SpectreLog.Debugln("Starting server - listening on port " + port)
	if global.Options.Protocol == "https" {
		err = http.ListenAndServeTLS(port, global.Options.CertFile, global.Options.PkeyFile, security.Validator{r})
	} else {
		err = http.ListenAndServe(port, security.Validator{r})
	}
	if err != nil {
		utils.SpectreLog.Fatalln("Server Error:", err.Error())
		// log.Fatalf("Error: %s", err.Error())
	}

	return nil
}
