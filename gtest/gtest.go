// Package gtest converts googletest benchmark output to benchfmt.
package gtest

import (
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/perf/v2/benchfmt"
)

var (
	// Format:
	// BM_Stat/64/real_time       16770 ns        16593 ns        42186
	// BM_LargeConsistent                 4.69           4.69   500000000
	benchmarkRe = regexp.MustCompile(`^BM_(\S+)\s+([0-9\.]+)(?: ns)?\s+([0-9\.]+)(?: ns)?\s+([0-9]+)\s*$`)
)

func Line(s string) (benchfmt.Result, bool) {
	m := benchmarkRe.FindStringSubmatch(s)
	if len(m) == 0 {
		return benchfmt.Result{}, false
	}

	name := m[1]

	iters := m[4]
	i, err := strconv.ParseUint(iters, 10, 64)
	if err != nil {
		// regexp says this is a number, conversion really shouldn't fail.
		panic(fmt.Sprintf("failed to parse %q as uint64: %v", iters, err))
	}

	r := benchfmt.Result{
		FullName: []byte(name),
		Iters:    int(i),
	}

	wall := m[2]
	w, err := strconv.ParseFloat(wall, 64)
	if err != nil {
		// regexp says this is a number, conversion really shouldn't fail.
		panic(fmt.Sprintf("failed to parse %q as float: %v", wall, err))
	}
	r.Values = append(r.Values, benchfmt.Value{
		Value: w,
		Unit:  "ns",
	})

	cpu := m[3]
	c, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		// regexp says this is a number, conversion really shouldn't fail.
		panic(fmt.Sprintf("failed to parse %q as float: %v", cpu, err))
	}
	r.Values = append(r.Values, benchfmt.Value{
		Value: c,
		Unit:  "cpu-ns",
	})

	return r, true
}
