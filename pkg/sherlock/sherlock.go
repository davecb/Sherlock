package sherlock

import (
	"time"
)

// Run runs lots of consulting detectives in parallel, the daemonic case
func Run(configFile string) {
	// load initially
	// for each section
	//   create a detective
		 detective()
	// loop reading worker output
}

// run one detective until told to stop
func detective() {
	// load specific file
	// loop on select
	//    do work
	//    wait for changes in file
	// 	      reload on change
}


// other operations, allowing one to try out new rules or search
// without specific rules

// Load a rule file
func Load(ruleFile string) {}

// Try a rule file, once
func Try(logFile string) {}

// Add a rule to a rule file, but only in memory
func Add(rule, version string, today time.Time) {}

// Subtract a rule from a rule file, only in memory
func Subtract(rule string) {}

// Save a rule file to disk
func Save(ruleFile string) {}
