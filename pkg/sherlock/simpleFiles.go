package sherlock

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

/*
 * Structs
 */

// Config provides the options for the various operations
type Config struct {
	Verbose bool
	Debug   bool
	Stop    bool
	Ruleset string
}

// PrintConfig displays a config struct's contents
func PrintConfig(conf Config) {
	log.Print("type Config struct {\n")
	log.Printf("    Verbose  bool = %v\n", conf.Verbose)
	log.Printf("    Debug    bool = %v\n", conf.Debug)
	log.Printf("    Stop     bool = %v\n", conf.Stop)
	log.Printf("    Ruleset  string = %q\n", conf.Ruleset)
	log.Print("}\n")
}

type rules []rule
type rule struct {
	pat  *regexp.Regexp
	date string
	vers string
}

// PrintRuleset does just that
func PrintRuleset(ruleset rules) { // nolint
	log.Print("type []rule {\n")
	log.Print("    // pat, date, vers\n")
	for _, r := range ruleset {
		log.Printf("    { %q, %q, %q }\n", r.pat, r.date, r.vers)
	}
	log.Print("}\n")
}

/*
 * Top-level functions
 */

// LoadRules loads a rule file. FIXME, called many times, hoist call
func LoadRules(ruleFile string) (rules, error) { // nolint: gocyclo
	var ruleset rules
	var record []string
	var version string
	var regex *regexp.Regexp
	var date string

	// precondition(ruleFile == "", "Programmer error:  no load-test .csv file")
	f, err := os.Open(ruleFile)
	if err != nil {
		wd, _ := os.Getwd()
		return nil, fmt.Errorf("%s, in directory %s, halting", err, wd)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	r := csv.NewReader(f)
	r.Comma = ','
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
	r.LazyQuotes = true

forloop:
	for line := 1; ; line++ {
		record, err = r.Read()
		switch err {
		case io.EOF:
			// eof is not an error
			break forloop
		case nil:
			// the desired case, fall through
		default:
			// everything else is an error
			return nil, fmt.Errorf("fatal error at line %d "+
				"in %q: %v", line, ruleFile, err)
		}

		switch len(record) {
		case 0:
			// skip blank lines quietly
			continue
		case 3:
			// the desired case,
			regex, date, version, err = compileRule(record[0], record[1], record[2])
			if err != nil {
				log.Printf("Can't compile an RE from %q, line %d of %q, ignored",
					record, line, ruleFile)
				continue
			}
		default: // for not, just accept others
			regex, date, version, err = compileRule(record[0], "", "")
			if err != nil {
				log.Printf("Can't compile an RE from %q, line %d of %q, ignored",
					record, line, ruleFile)
				continue
			}
		}
		ruleset = append(ruleset, rule{regex, date, version})
	}
	return ruleset, nil
}

// Try running the rules on a single log file
func Try(logFile string, cfg Config, ruleset rules) error {
	if cfg.Verbose {
		PrintConfig(cfg)
		PrintRuleset(ruleset)
	}
	return evaluate(logFile, ruleset, cfg.Verbose, cfg.Stop)
}

// evaluate tries a rule file, on a single input
// Note that we loop across individual REs, rather that concatenating
// them and trying to match that. The latter is ~124 times slower.
func evaluate(logFile string, ruleset rules, verbose, stop bool) error {
	var lines, hits int

	f, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("%s, halting", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	scanner := bufio.NewScanner(f)
outerLoop:
	for lines = 0; scanner.Scan(); lines++ {
		if err = scanner.Err(); err != nil {
			log.Printf("Bad line ignored, %v\n", err)
			continue
		}
		s := scanner.Text()
		if s == "" {
			// empty lines
			continue
		}
		if verbose {
			log.Printf("line=%q\n", s)

		}

		for _, r := range ruleset {
			//if verbose {
			//	log.Printf("trying %q\n", r.pat.String())
			//}
			if r.pat.FindStringIndex(s) != nil {
				// we found it in the rules, continue to the next line
				if verbose {
					log.Printf("found it\n")
				}
				continue outerLoop
			}
		}
		// postcondition: it wasn't in the ruleset if we got here

		hits++
		fmt.Printf("%s\n", s) // to stdout
		if verbose {
			fmt.Printf("this is new stuff... %s\n", s)
		}
		if stop {
			// Stop and return an error
			os.Exit(1)
		}

	}
	log.Printf("%d hits in %d lines of log\n", hits, lines) // to stderr
	return nil
}

// compileRule compiles a RE and optional date and version
func compileRule(rule, timestamp, version string) (*regexp.Regexp, string, string, error) {
	var regex *regexp.Regexp
	var date string
	var err error

	// later, see if we should convert all dates, or runs of digits or hex, into ".*"

	// Compile the pattern into a regexp.
	// prepend "(?i)" if you need to make it case-insensitive
	regex, err = regexp.Compile(rule)
	if err != nil {
		return nil, "", "", fmt.Errorf("ill-formed regexp %q",
			rule)
	}

	// Parse the time, as an ANSI C date/time
	if timestamp == "" {
		// if empty, it's now.
		timestamp = time.Now().Format(time.RFC3339)
	}
	if version == "" {
		// if it's empty, force it to zero.  FIXME hoist
		version = "0.0"
	}
	return regex, date, version, nil

}

/*
 * Daemon operations, TBD
 */

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

// add a rule to a rule file, but only in memory.
func add(ruleset rules, newRule, today, version string) (rules, error) {
	regex, date, version, err := compileRule(newRule, today, version)
	if err != nil {
		return nil, err
	}
	return append(ruleset, rule{regex, date, version}), nil
}

// Subtract a rule from a rule file, only in memory
func subtract(ruleset rules, removeRule, today, version string) (rules, error) {
	b := make(rules, 0)
	for _, x := range ruleset {
		if !x.pat.MatchString(removeRule) {
			b = append(b, x)
		}
	}
	return b, nil
}

// // Save a rule file to disk
// func save(ruleFile string) {}
