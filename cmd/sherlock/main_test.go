package main

import (
	"testing"
	"os"
)

func TestMainProgram(t *testing.T) {
	os.Args = []string{"sherlock",
		"--debug", 
		"--add", "zero",
		"--ruleset", "../scripts/ceph-log-filters/ceph.rules",
		"../scripts/ceph-log-filters/ceph.log"}
	main()
}
