package main

import "bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log" // spectre

var logger = log.New() // spectre

func init() {
	logger.Formatter = new(log.JSONFormatter) // spectre
	logger.Formatter = new(log.TextFormatter) // default // spectre
	logger.Level = log.DebugLevel             // spectre
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			logger.WithFields(log.Fields{ // spectre
				"omg":    true,
				"err":    err,
				"number": 100,
			}).Fatal("The ice breaks!")
		}
	}()

	logger.WithFields(log.Fields{ // spectre
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	logger.WithFields(log.Fields{ // spectre
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logger.WithFields(log.Fields{ // spectre
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logger.WithFields(log.Fields{ // spectre
		"temperature": -4,
	}).Debug("Temperature changes")

	logger.WithFields(log.Fields{ // spectre
		"animal": "orca",
		"size":   9009,
	}).Panic("It's over 9000!")
}
