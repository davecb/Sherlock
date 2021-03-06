package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	
	"github.com/davecb/Sherlock/pkg/sherlock"
)

const (
	r = "../scripts/ceph-log-filters/ceph.rules"
	a = "rule to add"
	s = "rule to subtract"
	v = "42"
	i = "test.ini"
	l = "../scripts/ceph-log-filters/ceph.log"
)


// test the happy paths
func TestSuccess(t *testing.T) {
	Convey("Given a panic testing framework", t, func() {

		var tests = []struct {
			descr    string
			ruleset  string
			add      string
			subtract string
			version  string
			initfile string
			log      []string
		}{
			{"just a ruleset and log", r, "", "", "", "", []string{l}},
			{"a ruleset, add and log", r, a, "", "", "", []string{l}},
			{"a ruleset, subtract and log", r, "", s, "", "", []string{l}},
			{"a ruleset, add and version to commit", r, a, "", v, "", []string{""}},
			// daemon next
		}
		for _, tx := range tests {
			Convey(tx.descr, func() {
				err := testableMain(tx.initfile, tx.log[:], sherlock.Config{// nolint:gotype
					Verbose: false,
					Debug: false,
					Ruleset: tx.ruleset,
					Add: tx.add,
					Subtract: tx.subtract,
					Version: tx.version,
				})
				So(err, ShouldEqual, nil)

			})
		}
	})
}

// test the falure paths, which have been set up to panic for goConvey tests
func TestFailures(t *testing.T) {
	Convey("Given a panic testing framework", t, func() {

		 var tests = []struct {
		 	match string
		 	ruleset string
		 	add string
		 	subtract string
		 	version string
		 	initfile string
		 	log []string
		 }{
			 {"You must provide a .ini file", "", "",
				 "", "", i, []string{""}},
			 {"You must provide a ruleset", "", "",
			 	"", "", "", []string{""}},
			 {"You must provide a log", r, "",
			 	"", "", "", []string{""}},
			 {"You must provide a rule to add or subtract", r, "",
				 "", v, "", []string{""}},
		 }
	     for _, tx := range tests {
			 Convey(tx.match, func() {
				 err := testableMain(tx.initfile, tx.log[:], sherlock.Config{ // nolint:gotype
					 Verbose: true,
					 Debug: true,
					 Ruleset: tx.ruleset,
					 Add: tx.add,
					 Subtract: tx.subtract,
					 Version: tx.version,
				 })
				 So( err.Error(), ShouldEqual, tx.match)
			 })
		 }
  	})
}
