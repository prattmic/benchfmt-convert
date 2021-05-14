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
	unitlessRe = regexp.MustCompile(`^\s*([0-9,\.]+)\s+([a-z0-9_.-]+)\s+(#.*)?$`)
)

// Stat with unit. Format:
//               1.11 msec task-clock                #    0.001 CPUs utilized
// /^\s+[0-9,\.]+\s+[a-z]+\s+[a-z0-9_.-]+\s+(#.*)?$/ {

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
	if m := unitlessRe.FindStringSubmatch(s); len(m) > 0 {
		value := m[1]
		value = strings.ReplaceAll(value, ",", "") // strip commas
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			// regexp says this is a number, conversion really shouldn't fail.
			panic(fmt.Sprintf("failed to parse %q as float64: %v", value, err))
		}

		name := m[2]
		name = "Benchmark" + capitalize(name)

		r := benchfmt.Result{
			FullName: []byte(name),
			Iters:    1,
			Values:   []benchfmt.Value{
				{
					Value: v,
				},
			},
		}
		return r, true
	}

	return benchfmt.Result{}, false
}
