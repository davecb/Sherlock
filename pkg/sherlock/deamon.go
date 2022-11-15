package sherlock

import (
	"fmt"
	ini "gopkg.in/ini.v1"
	"log"
	"os"
	"gopkg.in/fsnotify.v1"
	"io"
	"time"
	"bufio"
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

	sectionNames := cfg.SectionStrings()
	for _, name := range (sectionNames) {
		if name == "default" {
			continue
		}
		if !cfg.Section(name).HasKey("input") {
			return fmt.Errorf("missing input=<value> line for [%s] in %q, halting",
				name, iniFile)
		}
		if !cfg.Section(name).HasKey("rules") {
			return fmt.Errorf("missing rules=<value> line dor [%s] in %q, halting",
				name, iniFile)
		}
		input := cfg.Section(name).Key("input").String()
		rules := cfg.Section(name).Key("rules").String()
		go detective(input, rules)
	}

	<- stopChan
	return fmt.Errorf("deamon is not fully written yet")
}

// run one detective until told to stop
func detective(logFile, ruleFile string) {
	var err error

	log.Printf("daemon(%s, %s)\n", logFile, ruleFile)
	rules,  err := LoadRules(ruleFile)
	if err != nil {
		log.Fatalf("failed to load %v, %v, halting\n", ruleFile, err)
	}

	f, err := os.Open(logFile)
	if err != nil {
		log.Fatalf("%s, halting", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	r := bufio.NewReader(f)

	watcher, err := createWatcher(err, logFile)
	if  err != nil {
		log.Fatalf("%s", err)
	}
	defer func() {
		if err = watcher.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fred(r, watcher, rules, logFile)

	//time.Sleep(1 * time.Second)
	//stopChan <- true

}

// fred loops madly...
func fred(r *bufio.Reader, watcher *fsnotify.Watcher, rules rules, logFile string) {
	for {
		record, err := r.ReadString('\n')
		switch {
		case err == io.EOF:
			// just keep reading, even if we truncate...
			if watcher == nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				//log.Print("waiting for fsnotify\n")
				if err = waitForChange(watcher); err != nil {
					log.Fatalf("Fatal error waiting for fsnotify on %s, %v\n", logFile, err)
				}
			}
			continue
		case err != nil:
			log.Fatalf("Fatal error mid-way in %s: %s\n", logFile, err)
		}
		log.Printf("%s\n", record)
		// reload on change
	}
}
func createWatcher(err interface{}, logFile string) (*fsnotify.Watcher, error) {
	var watcher *fsnotify.Watcher

	watcher, _ = fsnotify.NewWatcher()
	if err != nil {
		return nil, nil	// to fall back to polling, we set watcher to nil
	}
	err = watcher.Add(logFile)
	if err != nil {
		return nil, fmt.Errorf("Fatal error addding %s to fsnotify: %s", logFile, err)
	}
	return watcher, nil
}

// waitForChange waits for the tail of a file to be written to
// cargo courtesy Satyajit Ranjeev, http://satran.in/2017/11/15/Implementing_tails_follow_in_go.html
func waitForChange(w *fsnotify.Watcher) error {
	for {
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}
