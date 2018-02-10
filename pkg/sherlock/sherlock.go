package sherlock

import (
	"log"
	"time"
	"encoding/csv"
	"os"
	"fmt"
	"io"
	"regexp"
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
	if cfg.Verbose {
		PrintConfig(cfg)
	}
	ruleset, err := load(cfg.Ruleset, cfg.Verbose)
	if err != nil {
		return err
	}
	//add(cfg.Add, "", nil)
	//subtract(cfg.Subtract)
	return evaluate(logFile, ruleset)
}

// Commit will update a rule file, triggering a daemon refresh
func Commit(cfg Config) error {
	if cfg.Verbose {
		PrintConfig(cfg)
	}
	//load(cfg.Ruleset, cfg.Verbose)
	//add(cfg.Add, cfg.Version, time.Now())
	//save(cfg.Ruleset)   // May change daemon
	return nil
}

// PrintConfig displays a config struct's contents
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

// Load a rule file.
func load(ruleFile string, verbose bool) (rules, error) {  // nolint: gocyclo
	var ruleset rules
	var warned = false

	// precondition(ruleFile == "", "Programmer error:  no load-test .csv file")
	f, err := os.Open(ruleFile)
	if err != nil {
		return nil, fmt.Errorf("%s, halting", err)
	}
	defer f.Close() // nolint

	r := csv.NewReader(f)
	r.Comma = ','
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
	
forloop:
	for line := 1; ; line++ {
		record, err := r.Read()
		switch err {
		case io.EOF:
			break forloop
		case nil:
			// the desired case, fall through
		default:
			return nil, fmt.Errorf("fatal error at line %d " +
				"in %q: %v", line, ruleFile, err)
		}


		switch len(record) {
		case 0:
			// skip blank lines quietly
			continue
		case 1:
			// warn once
			if !warned {
				warned = true
				log.Print("Lines with only the pattern field " +
					" encountered in %q, missing dates and versions " +
					"ignored\n", ruleFile)
			}
		case 3:
				// the desired case, fall out of the switch
		default:
			// ignore UFOs
			log.Printf("Ill-formed record %d (%q) ignored in %q\n",
				line, record, ruleFile)
			continue
		}

		// Compile the pattern into a regexp.
		// prepend "(?i)" to make it case-insensitive
		regex, err := regexp.Compile(record[0])
		if err != nil {
			log.Printf("Ill-formed regexp %q in line %d of %q, skipped\n",
				record[0], line, ruleFile)
			continue
		}

		// Parse the time, as an ANSI C date/time
		date, err := time.Parse(time.ANSIC, record[1])
		if err != nil {
			log.Printf("Ill-formed time %q in line %d of %q, ignored\n",
				record[1], line, ruleFile)
			date = time.Now()
		}
		ruleset = append(ruleset, rule{regex, date, record[2]})
	}
	if verbose {
		printRuleset(ruleset)
	}
	return ruleset, nil
}


// printRuleset does just that
func printRuleset(ruleset rules) {  // nolint
	log.Print("type []rule {\n")
	log.Print("    // pat, date, vers\n")
	for _, r := range ruleset {
		log.Printf("    { %q, %q, %q }\n", r.pat, r.date, r.vers)
	}
	log.Print("}\n")
}

// evaluate tries a rule file, once.
// Note that we loop across individual REs, rather that concatenating
// them and trying to match that. The latter is ~124 times slower.
func evaluate(logFile string, ruleset rules) error {
	f, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("%s, halting", err)
	}
	defer f.Close() // nolint
	scanner := bufio.NewScanner(f)
outerLoop:
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
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
