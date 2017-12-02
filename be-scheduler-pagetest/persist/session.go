package persist

import (
	"time"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
)

func InitSession(session mdb.ISession) {
	dInfo := &mdb.DialInfo{
		Addrs:    []string{global.Options.Mongo},
		Timeout:  5 * time.Second,
		FailFast: true,
		Database: constants.DB_NAME,
		// Source:        "admin",
		Username:      constants.DB_USERNAME,
		Password:      constants.DB_PASSWORD,
		RetryInterval: 10 * time.Second,
		DoResolveAddr: global.Options.EnableSvcDiscovery,
		TLSConnection: global.Options.IsMongoSecure,
	}
	err := session.InitSession(dInfo)
	if err != nil {
		utils.SpectreLog.Println("mongodb initialization errored:", err.Error())
		utils.SpectreLog.Fatal("mongodb initialization errored:", err.Error())
	}
}

// EndSession ends the mongodb session
func EndSession(session mdb.ISession) {
	if session != nil {
		session.EndSession()
	}
}
