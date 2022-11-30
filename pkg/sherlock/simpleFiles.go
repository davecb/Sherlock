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
	Verbose  bool
	Debug    bool
	Ruleset  string
	Add      string
	Subtract string
	Version  string
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

type rules []rule
type rule struct {
	pat  *regexp.Regexp
	date time.Time
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
	var date time.Time
	var warned = false

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
		case 1:
			// warn once
			if !warned {
				warned = true
				log.Printf("Lines with only the pattern field "+
					" encountered in %q, missing dates and versions "+
					"ignored\n", ruleFile)
			}
			regex, date, version, err = compileRule(record[0], "", "")
			if err != nil {
				log.Printf("Can't compile an RE from %q, line %d of %q, ignored",
					record, line, ruleFile)
				continue
			}
		case 3:
			// the desired case,
			regex, date, version, err = compileRule(record[0], record[1], record[2])
			if err != nil {
				log.Printf("Can't compile an RE from %q, line %d of %q, ignored",
					record, line, ruleFile)
				continue
			}
		default:
			// ignore UFOs, loudly
			log.Printf("Ill-formed record %d (%q) ignored in %q\n",
				line, record, ruleFile)
			continue
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
	return evaluate(logFile, ruleset, cfg.Verbose)
}

// evaluate tries a rule file, on a single input
// Note that we loop across individual REs, rather that concatenating
// them and trying to match that. The latter is ~124 times slower.
func evaluate(logFile string, ruleset rules, verbose bool) error {
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
	for scanner.Scan() {
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
			if r.pat.FindStringIndex(s) != nil {
				// we found it, skip this whole line
				if verbose {
					log.Printf("found it\n")
				}
				continue outerLoop
			}
		}
		if verbose {
			fmt.Printf("this is new stuff: %q\n", s)
		}
	}
	return nil
}

// compileRule compiles a RE and optional date and version  FIXME used twice or more, hoist
func compileRule(rule, today, version string) (*regexp.Regexp, time.Time, string, error) {
	var regex *regexp.Regexp
	var date time.Time
	var err error

	// protect any special characters
	rule = regexp.QuoteMeta(rule)

	// later, convert all dates, or runs of digits or hex, into ".*"

	// Compile the pattern into a regexp.
	// prepend "(?i)" if you need to make it case-insensitive
	regex, err = regexp.Compile(rule)
	if err != nil {
		return nil, time.Time{}, "", fmt.Errorf("ill-formed regexp %q",
			rule)
	}

	// Parse the time, as an ANSI C date/time
	if today == "" {
		// if empty, it's now.
		date = time.Now()
	} else {
		date, err = time.Parse(time.ANSIC, today)
		if err != nil {
			log.Printf("Ill-formed time %q in rule { %q, %q, %q }, ignored\n",
				today, rule, today, version)
			date = time.Now()
		}
	}
	if version == "" {
		// if it's empty, force it to zero.  FIXME hoist
		version = "0.0"
	}
	return regex, date, version, nil

}

/*
 * Daemon operations
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
