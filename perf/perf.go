// Package perf converts `perf stat` output to benchfmt.
package perf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/perf/v2/benchfmt"
)

var (
	// Unitless stat. Format:
	//          1,291,018      cycles                    #    1.167 GHz
	unitlessRe = regexp.MustCompile(`^\s*([0-9,\.]+)\s+([a-zA-Z0-9_\.:-]+)(\s+(\(.*%\))?)?(\s+(#.*)?)?$`)

	// Stat with unit. Format:
	//               1.11 msec task-clock                (55.55%)  #    0.001 CPUs utilized
	unitRe = regexp.MustCompile(`^\s*([0-9,\.]+)\s+([a-zA-Z]+)\s+([a-zA-Z0-9_\.:-]+)(\s+(\(.*%\))?)?(\s+(#.*)?)?$`)

	// Total elapsed wall time. Format:
	//      1656.917143299 seconds time elapsed
	wallRe = regexp.MustCompile(`^\s*([0-9\.]+)\s+seconds time elapsed$`)

	// Total elapsed user time. Format:
	//      1656.917143299 seconds user
	userRe = regexp.MustCompile(`^\s*([0-9\.]+)\s+seconds user$`)

	// Total elapsed system time. Format:
	//      1656.917143299 seconds sys
	sysRe = regexp.MustCompile(`^\s*([0-9\.]+)\s+seconds sys$`)
)

// capitalize capitalizes the first character in s.
func capitalize(s string) string {
	first := ""
	for i, r := range s {
		if i == 0 {
			first = string(unicode.ToUpper(r))
			continue
		}
		// Second iteration gives us the start of the second character.
		return first + s[i:]
	}
	return first
}

func Line(s string) (benchfmt.Result, bool) {
	var name, value, unit string

	if m := wallRe.FindStringSubmatch(s); len(m) > 0 {
		value = m[1]
		unit = "sec"
		name = "wall-time"
	} else if m := userRe.FindStringSubmatch(s); len(m) > 0 {
		value = m[1]
		unit = "sec"
		name = "user-time"
	} else if m := sysRe.FindStringSubmatch(s); len(m) > 0 {
		value = m[1]
		unit = "sec"
		name = "system-time"
	} else if m := unitlessRe.FindStringSubmatch(s); len(m) > 0 {
		value = m[1]
		unit = "val"
		name = m[2]
	} else if m := unitRe.FindStringSubmatch(s); len(m) > 0 {
		value = m[1]
		unit = m[2]
		name = m[3]
	}

	if value == "" {
		return benchfmt.Result{}, false
	}

	value = strings.ReplaceAll(value, ",", "") // strip commas
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		// regexp says this is a number, conversion really shouldn't fail.
		panic(fmt.Sprintf("failed to parse %q as float64: %v", value, err))
	}

	// perf outputs clock events in milliseconds. benchstat doesn't
	// understand "msec" (only "sec" and "ns"), so convert.
	if unit == "msec" {
		v /= 1000
		unit = "sec"
	}

	r := benchfmt.Result{
		FullName: []byte(capitalize(name)),
		Iters:    1,
		Values: []benchfmt.Value{
			{
				Value: v,
				Unit:  unit,
			},
		},
	}
	return r, true
}
