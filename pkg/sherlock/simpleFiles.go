package sherlock

import (
	"time"
	"regexp"
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

type rules []rule
type rule struct {
	pat		*regexp.Regexp
	date 	time.Time
	vers 	string
}


// Try running the rules on a single log file
func Try(logFile string, cfg Config) error {
	if cfg.Verbose {
		PrintConfig(cfg)
	}
	ruleset, err := LoadRules(cfg.Ruleset)
	if err != nil {
		return err
	}
	if cfg.Add != "" {
		ruleset, err = add(ruleset, cfg.Add,  "", "")
		if err != nil {
			return err
		}
	}
	if cfg.Subtract != "" {
		ruleset, err = subtract(ruleset, cfg.Subtract,  "", "")
		if err != nil {
			return err
		}
	}
	if cfg.Verbose {
		PrintRuleset(ruleset)
	}
	return evaluate(logFile, ruleset)
}

// Commit will update a rule file, triggering a daemon refresh
func Commit(cfg Config) error {
	if cfg.Verbose {
		PrintConfig(cfg)
	}
	// load(cfg.Ruleset)
	// add(cfg.Add, cfg.Version, time.Now())
	// save(cfg.Ruleset)   // May change daemon
	return nil
}
