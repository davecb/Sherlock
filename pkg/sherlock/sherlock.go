package sherlock

import (
	"log"
)

// Config provides the options for Try and Commit operations
type Config struct {
	Verbose  bool
	Debug    bool
	Ruleset  string
	Add      string
	Subtract string
	Version  string
}

// Run runs lots of consulting detectives in parallel, the daemonic case
// It does not use Config, but a .ini file instead.
func Run(configFile string) error {
	// load initially
	// for each section
	//   create a detective
		 detective()
	// loop reading worker output
	return nil
}

// run one detective until told to stop
func detective() {
	// load specific file
	// loop on select
	//    do work
	//    wait for changes in file
	// 	      reload on change
}


// Try running the rules on a single log file
func Try(logFile string, cfg Config) error {
	//load(cfg.Ruleset)
	//add(cfg.Add, "", nil)
	//subtract(cfg.Subtract)
	//evaluate(logFile)
	log.Fatal("Try is not implemented yet")
	return nil
}

// Commit will update a rule file, triggering a daemon refresh
func Commit(cfg Config) error {
	//load(cfg.Ruleset)
	//add(cfg.Add, cfg.Version, time.Now())
	//save(cfg.Ruleset)   // May change daemon
	log.Fatal("Commit is not implemented yet")
	return nil
}


// other operations, allowing one to try out new rules or search
// without specific rules

//// Load a rule file
//func load(ruleFile string) {}
//
//// evaluate tries rule file, once
//func evaluate(logFile string) {}
//
//// Add a rule to a rule file, but only in memory
//func add(rule, version string, today time.Time) {}
//
//// Subtract a rule from a rule file, only in memory
//func subtract(rule string) {}
//
//// Save a rule file to disk
//func save(ruleFile string) {}
