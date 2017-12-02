package utils

import (
	"os"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"

	"golang.org/x/net/context"
	//log "github.com/omidnikta/logrus"
)

var (
	SpectreLog                *log.Entry
	TRACKING_ID_FIELD_FOR_LOG = "trackingId"
	CUSTOMER_ID_FIELD_FOR_LOG = "customerId"
)

func init() {
	logger := log.New()
	if os.Getenv("DEBUG") == "true" {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.InfoLevel
	}
	SpectreLog = logger.WithField("Service Name", "Spectre Pagetest Microservice")
}

func GetLoggerWithTrackingIdFromContext(ctx context.Context, log *log.Entry) *log.Entry {
	val := ctx.Value(TRACKING_ID_FIELD_FOR_LOG)
	if val != nil {
		log = log.WithField(TRACKING_ID_FIELD_FOR_LOG, val)
	}
	cus := ctx.Value(CUSTOMER_ID_FIELD_FOR_LOG)
	if cus != nil {
		log = log.WithField(CUSTOMER_ID_FIELD_FOR_LOG, cus)
	}
	return log
}
