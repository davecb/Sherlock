package main

import (
	"github.com/davecb/Sherlock/pkg/sherlock"   // nolint:gotype

	"flag"
	"fmt"
	"os"
	"log"
)

// usage reports how to use sherlock
func usage() {  // nolint
	fmt.Fprint(os.Stderr, "Usage: sherlock --deamon ini-file, or\n" +
		"       sherlock --options log-file\n")
	// FIXME "--commit" later
	flag.PrintDefaults()
}

// main parses options and starts the work
func main() {
	var initFile, add, subtract, ruleset, commitVersion string
	var verbose, debug bool

	// run as a daemon, or
	flag.StringVar(&initFile, "deamon", "", "specify daemon's .ini file")

	// add a rule to a ruleset, with a version and date
	flag.StringVar(&commitVersion, "commit", "", "commit to a specified version")
	// requires --add "rule" and --ruleset "path"

	// try various combinations of tests
	flag.StringVar(&add, "add", "", "specify a rule to add")
	flag.StringVar(&subtract, "subtract", "", "specify a rule to not use")
	flag.StringVar(&ruleset, "ruleset", "", "specify a ruleset")

	flag.BoolVar(&verbose, "verbose", false, "Turn verbose logging on")
	flag.BoolVar(&debug, "debug", false, "Turn debug logging on")
	flag.Parse()

	log.SetFlags(0) // log.Lshortfile | log.Ldate | log.Ltime) // show file:line in logs

	if debug {
		log.Printf("%d flags\n", flag.NFlag())
		log.Printf("%d args\n", flag.NArg())
		log.Printf("args=%v\n", flag.Args())
	}
	err := testableMain(initFile, flag.Args(), sherlock.Config {  // nolint:gotype
		Verbose:    verbose,
		Debug:      debug,
		Ruleset:	ruleset,
		Add:		add,
		Subtract:	subtract,
		Version:    commitVersion,
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}

// testableMain, for use with BDD tests
func testableMain(initFile string, args []string, conf sherlock.Config) error {  // nolint: gotype
	
	// see if we're to be a runDaemon
	if initFile != "" {
		if len(args) < 1 || args[0] == "" {
			// usage()
			return fmt.Errorf( "you must provide a .ini file") // nolint
		}
		runDaemon(flag.Arg(1) ) // nolint: errcheck
		// never exits normally
	}
	// all subsequent uses require a ruleset
	if conf.Ruleset == "" {
		// usage()
		return fmt.Errorf("you must provide a ruleset")
	}


	// see if we're to commit a change to the runDaemon
	if conf.Version != "" {
		return commit(conf)
	}

	// Otherwise try running rules against one or more log files
	if len(args) < 1 || args[0] == "" {
		// usage()
		return fmt.Errorf("you must provide a log") // nolint
	}
	// apply a ruleset to logfiles, with and without specific rules
	for _, arg := range args {
		err := try(arg, conf)
		if err != nil {
			return err
		}
	}
	return nil
}


// runDaemon runs a runDaemon from an ini file
// syntactially checks config file
func runDaemon(iniFile string) error {
	if iniFile == "" {
		return fmt.Errorf("you must provide a .ini file")
	}
	sherlock.Run(iniFile) // nolint
	// should never exit
	return fmt.Errorf("sherlock.Run exited unexpectedly, halting")
}

// commit tries to update a config file, thus updating any daemons
// syntactially checks add and version
func commit(conf sherlock.Config) error { // nolint: gotype
	if conf.Add == "" && conf.Subtract == "" {
		return fmt.Errorf("you must provide a rule to add or subtract")
	}
	if conf.Version == "" {
		// belt and suspenders: this is currently unreachable
		return fmt.Errorf("you must provide a commit version-string file")
	}
	err := sherlock.Commit(conf) // nolint: gotype
	if err != nil {
		return fmt.Errorf("update of ruleset %q with rule %q failed, %v",
			conf.Ruleset, conf.Add, err)
	}
	return nil
}

// try to evaluate one file
func try(logFile string, conf sherlock.Config ) error { // nolint: gotype
	err := sherlock.Try(logFile, conf) // nolint: gotype
	if err != nil {
		return fmt.Errorf("failed to evaluate log %q using ruleset %q, %v",
			logFile, conf.Ruleset, err)
	}
	return nil
}

