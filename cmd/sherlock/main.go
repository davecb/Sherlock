package main

import (
	"github.com/davecb/Sherlock/pkg/sherlock"

	"flag"
	"fmt"
	"os"
	"log"
)
var iniFile, add, subtract, ruleset, commitVersion string
var verbose, debug bool

// usage reports how to use sherlock
func usage() {
	//nolint
	fmt.Fprint(os.Stderr, "Usage: sherlock --daemonic ini-file, or\n" +
		"       sherlock --options log-file\n")
	// reserve --commit for later
	flag.PrintDefaults()
}

// main parses options and starts the work
func main() {
	// run as a --daemon, or
	flag.StringVar(&iniFile, "daemonic", "", "specify daemon's .ini file")

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

	log.SetFlags(0) //log.Lshortfile | log.Ldate | log.Ltime) // show file:line in logs

	if debug {
		log.Printf("%d flags\n", flag.NFlag())
		log.Printf("%d args\n", flag.NArg())
		log.Printf("args=%v\n", flag.Args())
	}
	
	// see if we're to be a daemon
	if iniFile != "" {
		if flag.NArg() < 1 {
			usage()
			log.Fatal( "You must supply a config file\n") //nolint
		}
		daemonic(flag.Arg(1) )
		// never exits normally
	}
	// all subsequent uses require a ruleset
	if ruleset == "" {
		usage()
		log.Fatalf("You must provide a ruleset")
	}


	// see if we're to commit a change to the daemon
	if commitVersion != "" {
		commit()
		os.Exit(0)
	}

	// Otherwise try running rules against one or more log files
	if flag.NArg() < 1 {
		usage()
		log.Fatal("You must supply a log file\n") //nolint
	}
	// apply a ruleset to logfiles, with and without specific rules
	for _, arg := range flag.Args() {
		try(arg)
	}
	os.Exit(0)
}


// daemonic runs a daemon from an ini file
// syntactially checks config file
func daemonic(iniFile string) {
	if iniFile == "" {
		log.Fatal("You must supply a .ini file\n")
	}
	sherlock.Run(iniFile) // nolint
	// should never exit
	log.Fatal("sherlock.Run exited unexpectedly, halting\n")
}

// commit tries to update a config file, thus updating any daemons
// syntactially checks add and version
func commit() {
	if add == "" {
		log.Fatal("You must supply a rule to add\n")
	}
	if commitVersion == "" {
		log.Fatal("You must supply a commit version-string file\n")
	}
	err := sherlock.Commit(sherlock.Config {
		Verbose:    verbose,
		Debug:      debug,
		Ruleset:	ruleset,
		Add:		add,
		Subtract:	subtract,
		Version:    commitVersion,
	})
	if err != nil {
		log.Fatalf("Update of ruleset %q with rule %q failed, %v\n",
			ruleset, add, err)
	}
}

// try to evaluate one file
func try(arg string) {
	err := sherlock.Try(arg, sherlock.Config {
		Verbose:    verbose,
		Debug:      debug,
		Ruleset:	ruleset,
		Add:		add,
		Subtract:	subtract,
	})
	if err != nil {
		log.Fatalf("Failed to evaluate log %q using ruleset %q, %v\n",
			arg, ruleset, err)
	}
}

