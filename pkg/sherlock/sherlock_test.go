package sherlock

import (
	"testing"
)


// nolint
func BenchmarkSherlockRE(b *testing.B) {
	for i :=0; i < b.N; i++ {
		Try("./ceph.log", Config{// nolint:gotype
			Verbose: false,
			Debug: false,
			Ruleset: "./ceph.rules",
			Add: "",
			Subtract: "",
			Version: "",
		})
	}
}

// With the loop-based comparison, this took
// pkg: github.com/davecb/Sherlock/pkg/sherlock
// 2000	    968,182 ns/op
// PASS
// ok  	github.com/davecb/Sherlock/pkg/sherlock	2.052s

// With concatenated REs it took
// pkg: github.com/davecb/Sherlock/pkg/sherlock
// 10	 117,029,945 ns/op
// PASS
// ok  	github.com/davecb/Sherlock/pkg/sherlock	1.308s

// Heavily concatenated REs with .*s can take very long times,
// in this case ~120 times longer. Less with " *", but still large.
