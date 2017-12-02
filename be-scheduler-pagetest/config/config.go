// Package config contains the code to process environmental variables, config files, etc.
package config

import (
	"github.com/goanywhere/env" // see: https://github.com/goanywhere/env
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"strings"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/helper"
)

// read config variables from environmental variables or '.env' file and
// populate globals.Options with the values
func SetupConfig() {
	// Load '.env' from /usr/bin directory.
	env.Load(".env")
	global.Options.Environment = strings.ToLower(env.String("ENVIRONMENT"))
	// Get secret key from environmental variable or .env file
	// init SecretKey field in global options
	global.Options.Protocol = env.String("PROTOCOL")
	global.Options.SecretKey = env.String("SECRETKEY")
	global.Options.Port = env.Int("GATEWAY_PORT")
	global.Options.CertFile = env.String("CERTFILE")
	global.Options.PkeyFile = env.String("PKEYFILE")
	global.Options.WptServer = env.String("WPT_SERVER")
	utils.SpectreLog.Debugln("Port number %d \n", env.Int("PORT"))
	utils.SpectreLog.Debugln("WPT server %s \n", global.Options.WptServer)
	global.Options.Mongo = env.String("SVCNAME_MONGO_SRVR")

	global.Options.RedisServers = env.String("REDIS_URL")
	global.Options.SchedulerEnabled = env.Bool("SCHEDULER_ENABLED")
	global.Options.SchedulerInt = env.Int("SCHEDULER_INTERVAL")
	global.Options.SchedulerSleep = env.Int("SCHEDULER_SLEEP_INTERVAL")
	global.Options.ValidateToken = env.Bool("VALIDATE_TOKEN")
	global.Options.EnableSvcDiscovery = env.Bool("ENABLE_SERVICE_DISCOVERY")
	global.Options.ClientId = env.String("CLIENT_ID")
	global.Options.ClientSecret = env.String("CLIENT_SECRET")
	global.Options.Elastic = env.String("ELASTIC")
	global.Options.EsUser = env.String("ES_USER")
	global.Options.EsPassword = env.String("ES_PASSWORD")
	global.Options.IsMongoSecure = env.Bool("IS_MONGO_SECURE")
	utils.SpectreLog.Debugln("ES " + env.String("ELASTIC"))
	helper.InitGradeMap()
	helper.InitCategoryMap()
}
