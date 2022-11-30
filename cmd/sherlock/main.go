package main

import (
	"github.com/davecb/Sherlock/pkg/sherlock" // nolint:gotype

	"flag"
	"fmt"
	"log"
	"os"
)

// usage reports how to use sherlock
func usage() { // nolint
	fmt.Fprint(os.Stderr, "Usage: sherlock --options log-file\n") // nolint:gas
	flag.PrintDefaults()
}

// main parses options and starts the work
func main() {
	var initFile, ruleset string
	var verbose, debug bool

	// run from the command-line
	flag.StringVar(&ruleset, "ruleset", "", "specify a ruleset")

	// run as a daemon, a proposed future feature
	//flag.StringVar(&initFile, "deamon", "", "specify daemon's .ini file")
	// try various combinations of tests
	//flag.StringVar(&add, "add", "", "specify a rule to add")
	//flag.StringVar(&subtract, "subtract", "", "specify a rule to not use")
	// add a rule to a ruleset, with a version and date
	//flag.StringVar(&commitVersion, "commit", "", "commit to a specified version")
	// requires --add "rule" and --ruleset "path"

	flag.BoolVar(&verbose, "v", false, "Turn verbose logging on")
	flag.BoolVar(&debug, "d", false, "Turn debug logging on")
	flag.Parse()

	log.SetFlags(0) // log.Lshortfile | log.Ldate | log.Ltime) // show file:line in logs

	err := testableMain(initFile, flag.Args(), sherlock.Config{ // nolint:gotype
		Verbose: verbose,
		Debug:   debug,
		Ruleset: ruleset,
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}

// testableMain, for use with BDD tests
func testableMain(initFile string, args []string, conf sherlock.Config) error { // nolint: gotype

	if conf.Ruleset == "" {
		// usage()
		return fmt.Errorf("you must provide a ruleset")
	}
	ruleset, err := sherlock.LoadRules(conf.Ruleset)
	if err != nil {
		return err
	}

	// Try running rules against one or more log files
	if len(args) < 1 || args[0] == "" {
		// usage()
		return fmt.Errorf("you must provide a log") // nolint
	}
	// apply a ruleset to logfiles, with and without specific rules
	for _, arg := range args {
		err := sherlock.Try(arg, conf, ruleset) // nolint: gotype
		if err != nil {
			return fmt.Errorf("failed to evaluate log %q using ruleset %q, %v",
				arg, conf.Ruleset, err)
		}
	}
	return nil
}
