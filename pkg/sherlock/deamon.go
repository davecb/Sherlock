package sherlock

import (
	"fmt"
	ini "gopkg.in/ini.v1"
	"log"
	"time"
)

var stopChan = make(chan bool)

// Run runs lots of consulting detectives in parallel, the daemon case
// It uses an .ini file to direct it.
func Run(iniFile string) error {
	var err error


	cfg, err := ini.InsensitiveLoad(iniFile)
	if err != nil {
		return err
	}

	names := cfg.SectionStrings()
	for _, name := range (names) {
		if name == "default" {
			continue
		}
		if !cfg.Section(name).HasKey("input") { // FIXME?
			return fmt.Errorf("missing input=<value> in %q, halting", name)
		}
		if !cfg.Section(name).HasKey("rules") {
			return fmt.Errorf("missing rules=<value> in %q, halting", name)
		}
		input := cfg.Section(name).Key("input").String()
		rules := cfg.Section(name).Key("rules").String()
		go detective(input, rules)
	}

	<- stopChan
	return fmt.Errorf("deamon is not fully written yet")


	// load initially
	// for each section
	//   create a detective
	
	// loop reading worker output
	//return nil
}

// run one detective until told to stop
func detective(input, rules string) {
	log.Printf("daemon(%s, %s)\n", input, rules)
	time.Sleep(1 * time.Second)
	stopChan <- true
	// load specific file
	// loop on select
	//    do work
	//    wait for changes in file
	// 	      reload on change
}

