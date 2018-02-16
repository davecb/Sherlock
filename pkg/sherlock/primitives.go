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

// load a rule file.
func load(ruleFile string) (rules, error) {  // nolint: gocyclo
	var ruleset rules
	var record []string
	var version string
	var regex *regexp.Regexp
	var date time.Time
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
		record, err = r.Read()
		switch err {
		case io.EOF:
			// eof is not an error
			break forloop
		case nil:
			// the desired case, fall through
		default:
			// everything else is an error
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
				log.Printf("Lines with only the pattern field " +
					" encountered in %q, missing dates and versions " +
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


// PrintRuleset does just that
func PrintRuleset(ruleset rules) {  // nolint
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
		if s == "" {
			// empty lines
			continue
		}
		for _, r := range ruleset {
			if r.pat.FindStringIndex(s) != nil {
				// we found it, skip this whole line
				continue outerLoop
			}
		}
		fmt.Printf("new stuff: %q\n",  s)
	}
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

// compileRule compiles a RE and optional date and version  FIXME use twice
func compileRule(rule, today, version string) (*regexp.Regexp, time.Time, string, error){
	var regex *regexp.Regexp
	var date time.Time
	var err error
	
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
	return regex, date, version, nil


}

 // Subtract a rule from a rule file, only in memory
func subtract(ruleset rules, removeRule, today, version string) (rules, error) {
	b := make(rules,0)
	for _, x := range ruleset {
		if !x.pat.MatchString(removeRule) {
			b = append(b, x)
		}
	}
	return b, nil
 }

// // Save a rule file to disk
// func save(ruleFile string) {}
