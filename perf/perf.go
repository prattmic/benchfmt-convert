package perf

import (
	//"regexp"

	"golang.org/x/perf/v2/benchfmt"
)

// Unitless stat. Format:
//          1,291,018      cycles                    #    1.167 GHz
// /^\s+[0-9,\.]+\s+[a-z0-9_.-]+\s+(#.*)?$/ {

// Stat with unit. Format:
//               1.11 msec task-clock                #    0.001 CPUs utilized
// /^\s+[0-9,\.]+\s+[a-z]+\s+[a-z0-9_.-]+\s+(#.*)?$/ {

func Line(s string) (benchfmt.Result, bool) {
	return benchfmt.Result{}, false
}
