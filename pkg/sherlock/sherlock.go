package sherlock

import (
	"log"
	"time"
	"encoding/csv"
	"os"
	"fmt"
	"io"
	"regexp"
	"strings"
	"bufio"
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
	PrintConfig(cfg)
	rules, err := load(cfg.Ruleset)
	if err != nil {
		return err
	}
	//add(cfg.Add, "", nil)
	//subtract(cfg.Subtract)
	evaluate(logFile, rules)
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

// Load a rule file
func load(ruleFile string) (rules, error) {
	var ruleset rules

	if ruleFile == "" {
		// FIXME. This is belt-and-suspenders: better as a precondition?
		// precondition(ruleFile == "", "Programmer error:  no load-test .csv file")
		return nil, fmt.Errorf("No load-test .csv file provided, halting.\n")
	}
	f, err := os.Open(ruleFile)
	if err != nil {
		return nil, fmt.Errorf("%s, halting.", err)
	}
	defer f.Close() // nolint

	r := csv.NewReader(f)
	r.Comma = ','
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
forloop:
	for line := 1; ; line++ {
		record, err := r.Read()
		switch {
		case err == io.EOF:
			break forloop
		case err != nil:
			return nil, fmt.Errorf("Fatal error mid-way in %s: %s\n", ruleFile, err)
		}
		if len(record) < 3 {
			log.Printf("ill-formed record %d (%q) ignored\n",
				line, record)
			continue
		}

		regex, err := createRegexp(record[0])
		if err != nil {
			log.Printf("ill-formed regexp %q ignored in line %d\n",
				record[0], line)
			continue
		}

		time, err := time.Parse(time.ANSIC, record[1])
		if err != nil {
			log.Printf("ill-formed time %q ignored in line %d\n",
				record[1], line)
			continue
		}
		ruleset = append(ruleset, rule{regex, time, record[2]})
	}
	printRuleset(ruleset)
	return ruleset, nil
}

// createRegexp removes RE characters, which may or may not be a bad idea...
// right now it's very desirable
func createRegexp(s string) (*regexp.Regexp, error) {
	// remove RE characters EXCEPT . and *, which I use
	for _, x := range []string{"[", "]", "+", "?", "\\", "{", "}" } {
		s = strings.Replace(s, x, ".",-1)
	}
	// prepend "(?i)" to make it case-insensitive
	return regexp.Compile(s)
}

func printRuleset(set rules) {
	log.Print("type []rule {\n")
	for _, r := range set {
		log.Printf("    { %q, %q, %q }\n", r.pat, r.date, r.vers)
	}
	log.Print("}\n")
}

// evaluate tries rule file, once
func evaluate(logFile string, ruleset rules) error {
	f, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("%s, halting.", err)
	}
	defer f.Close() // nolint
	scanner := bufio.NewScanner(f)
	
outerLoop:
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Printf("Bad line ignored, %v\n", err)
			continue
		}
		s := scanner.Text()
		for _, r := range ruleset {
			if r.pat.FindStringIndex(s) != nil {
				// we found it, skip this whole line
				continue outerLoop
			}
		}
		log.Printf("new stuff: %q\n",  s)
	}
	return nil
}

//// Add a rule to a rule file, but only in memory
//func add(rule, version string, today time.Time) {}
//
//// Subtract a rule from a rule file, only in memory
//func subtract(rule string) {}
//
//// Save a rule file to disk
//func save(ruleFile string) {}
