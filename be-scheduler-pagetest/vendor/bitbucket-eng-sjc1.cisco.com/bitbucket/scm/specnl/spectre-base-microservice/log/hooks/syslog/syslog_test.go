package log_syslog // spectre

import (
	"log/syslog"
	"testing"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log" // spectre
)

func TestLocalhostAddAndPrint(t *testing.T) {
	logger := log.New() // spectre
	hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")

	if err != nil {
		t.Errorf("Unable to connect to local syslog.")
	}

	logger.Hooks.Add(hook) // spectre

	for _, level := range hook.Levels() {
		if len(logger.Hooks[level]) != 1 { // spectre
			t.Errorf("SyslogHook was not added. The length of log.Hooks[%v]: %v", level, len(logger.Hooks[level])) // spectre
		}
	}

	logger.Info("Congratulations!") // spectre
}
