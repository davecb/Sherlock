package main

import (
	"github.com/davecb/Sherlock/pkg/sherlock"
	"time"
)

func main() {
	// run as a --daemon, or
	// try a --ruleset on just one --log
	// try --with a new rule and a ruleset, on just one log
	// try --without a rule on just one log
	// --add a rule to a ruleset, with a version and date
	daemonic()
}

func daemonic() {
	var configFile, ruleFile, logFile, rule, version string
	var today time.Time

	sherlock.Run(configFile)  // deamon case

	sherlock.Load(ruleFile)   // interactive or batch case
	sherlock.Add(rule, version, today)
	sherlock.Subtract(rule)
	sherlock.Try(logFile)
	sherlock.Save(ruleFile)   // May change daemon
}

