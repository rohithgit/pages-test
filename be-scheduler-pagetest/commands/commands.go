// Package commands contains the code to process the service flags and commands.
package commands

import (
	"flag"
	"os"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/version"
)

// InitCommands parses the command line commands and flags
// updates the global options
func InitCommands() error {

	printUsage := false
	printOptions := false
	printVersion := false

	// read flags and commands from the command line and update the
	// global options
	flgset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flgset.BoolVar(&printUsage, "h", false, "")
	flgset.BoolVar(&printUsage, "help", false, "Print help")
	flgset.BoolVar(&printOptions, "options", false, "Print options")
	flgset.BoolVar(&printVersion, "version", false, "Print version")
	flgset.BoolVar(&global.Options.Mock, "mock", global.Options.Mock, "Mock client response")
	flgset.BoolVar(&global.Options.Debug, "debug", global.Options.Debug, "Output debug log data")
	flgset.BoolVar(&global.Options.Verbose, "v", global.Options.Verbose, "Verbose output")
	flgset.IntVar(&global.Options.Port, "p", global.Options.Port, "Server `Port`")
	flgset.IntVar(&global.Options.Port, "port", global.Options.Port, "Server `Port`")
	flgset.Parse(os.Args[1:])

	// if requested, print usage/help and exit
	if printUsage {
		utils.SpectreLog.Debugln("\n%s [options]\n\n  options\n  -------\n\n", os.Args[0])
		flgset.PrintDefaults()
		utils.SpectreLog.Debugln("")
		os.Exit(0)
	}

	// if requested, print options and exit
	if printOptions {
		utils.PrintFields("options", global.Options)
		os.Exit(0)
	}

	// if requested, print version and exit
	if printVersion {
		utils.SpectreLog.Debugln("%s\n", version.Version)
		os.Exit(0)
	}

	return nil
}
