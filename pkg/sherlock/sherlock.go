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

// Run runs lots of consulting detectives in parallel, the daemon case
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
	log.Panic("detective is not implemented yet")
}


// Try running the rules on a single log file
func Try(logFile string, cfg Config) error {
	//PrintConfig(cfg)
	//load(cfg.Ruleset)
	//add(cfg.Add, "", nil)
	//subtract(cfg.Subtract)
	//evaluate(logFile)
	return nil
}

// Commit will update a rule file, triggering a daemon refresh
func Commit(cfg Config) error {
	//PrintConfig(cfg)
	//load(cfg.Ruleset)
	//add(cfg.Add, cfg.Version, time.Now())
	//save(cfg.Ruleset)   // May change daemon
	return nil
}

func PrintConfig(conf Config) {
		log.Print("type Config struct {\n")
		log.Printf("    Verbose  bool = %v\n", conf.Verbose)
		log.Printf("    Debug    bool = %v\n", conf.Debug)
		log.Printf("    Ruleset  string = %q\n", conf.Ruleset)
		log.Printf("    Add      string = %q\n", conf.Add)
		log.Printf("    Subtract string = %q\n", conf.Subtract)
		log.Printf("    Version  string = %q\n", conf.Version)
		log.Print("}\n")
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
