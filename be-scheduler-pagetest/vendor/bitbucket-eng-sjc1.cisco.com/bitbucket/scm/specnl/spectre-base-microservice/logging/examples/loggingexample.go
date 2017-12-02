package main

import (
	"os"
	"spectre-base-microservice/logging"

	"github.com/Sirupsen/logrus"
)

func main() {
	Logex := logging.NewSpectreLog(logrus.New())
	Logex.LoggingInit("Text", os.Stdout, "Info")
	Logex.GenerateLog("Error", "test_microservice", "1234 2345 Test Msg")
	Logex.GenerateLog("Info", "test_microservice", "1234 2345 Info_Test Msg")
	// Sample output : ERRO[0000] 2016-03-24 18:02:43 test_microservice:loggingexample.go:13 1234 2345 Test Msg
}
