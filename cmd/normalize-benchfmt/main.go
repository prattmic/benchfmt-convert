// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Binary normalize-benchfmt converts benchmark measurements from other
// benchmark formats to Go's [benchfmt], for use with tools like [benchstat].
//
// Supported input formats:
// - Linux [perf stat] text output
// - [Google Benchmark] C++ benchmarking framework
//
// normalize-benchfmt reads input from stdin or the file passed as the first
// argument. It writes each input line to stdout unmodified unless it matches a
// measurement line from one of the input formats, in which case it is
// converted to benchfmt and written to stdout.
//
// # Example use
//
// Suppose we have a Go benchmark that we would like to collect various
// hardware performance counters from using perf stat.
//
//	package bench
//
//	import "testing"
//
//	func BenchmarkMakeMap(b *testing.B) {
//		for b.Loop() {
//			_ = make(map[int]int)
//		}
//	}
//
// To use with perf stat, build a standalone test binary.
//
//	$ go test -c -o bench.test
//
// And run it under perf stat.
//
//	$ for i in $(seq 10); do
//	  perf stat ./bench.test \
//	    -test.run ^$ \
//	    -test.bench MakeMap \
//	    -test.benchtime 500000000x \
//	    2>&1 | normalize-benchfmt | tee -a out.benchfmt
//	done
//
// Finally, view a summary with [benchstat]:
//
//	$ benchstat out.benchfmt
//	goos: linux
//	goarch: amd64
//	pkg: example.com/bench
//	cpu: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
//	           │ out.benchfmt │
//	           │    sec/op    │
//	MakeMap-12    7.006n ± 2%
//
//	            │ out.benchfmt │
//	            │     sec      │
//	Task-clock     3.526 ±  2%
//	Wall-time      3.509 ±  2%
//	User-time      3.512 ±  1%
//	System-time   22.02m ± 45%
//	geomean       989.0m
//
//	                 │ out.benchfmt │
//	                 │     val      │
//	Context-switches    595.0 ±  4%
//	Cpu-migrations      16.50 ± 70%
//	Page-faults         434.5 ±  1%
//	Cycles             15.25G ±  1%
//	Instructions       41.46G ±  0%
//	Branches           4.699G ±  0%
//	Branch-misses      24.41M ±  0%
//	geomean            2.269M
//
// Some things to note here:
//   - benchstat requires multiple runs to determine statistical significance
//     when comparing benchmarks, so we run perf stat multiple times.
//   - perf stat measures the entire process execution, so we only run a single benchmark.
//   - Measurements will also include process and benchmark setup, so we select a
//     -benchtime that takes a few seconds to amortize those costs.
//   - Measurements are converted to benchfmt verbetum. e.g., Task-clock is "sec"
//     (entire process duration), not "sec/op" like Go benchmarks report.
//   - For apples-to-apples comparisons of unnormalized units like "sec", use a
//     fixed iteration count (-benchtime=1234x for 1234 iterations). By default,
//     Go tests will adjust benchmark iterations to achieve ~1s. Comparisons of
//     total counts of cycles, instructions, etc don't make sense if the before and
//     after cases ran a different number of iterations.
//
// [benchfmt]: https://golang.org/design/14313-benchmark-format
// [benchstat]: https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
// [perf stat]: https://man7.org/linux/man-pages/man1/perf-stat.1.html
// [Google Benchmark]: https://github.com/google/benchmark
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/prattmic/benchfmt-convert/gtest"
	"github.com/prattmic/benchfmt-convert/perf"
	"golang.org/x/perf/v2/benchfmt"
)

type parser func(string) (benchfmt.Result, bool)

var formats = []parser{
	gtest.Line,
	perf.Line,
}

func run() error {
	if len(os.Args) > 2 {
		return fmt.Errorf("usage: %s [input]", os.Args[0])
	}

	f := os.Stdin
	if len(os.Args) == 2 {
		var err error
		f, err = os.Open(os.Args[1])
		if err != nil {
			return err
		}
		defer f.Close()
	}

	w := benchfmt.NewWriter(os.Stdout)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		found := false
		for _, fn := range formats {
			r, ok := fn(line)
			if !ok {
				continue
			}
			found = true

			if err := w.Write(&r); err != nil {
				return fmt.Errorf("error writing %+v: %v", r, err)
			}
			break
		}

		// Include unmatched lines in output for context. benchstat
		// ignores them.
		if !found {
			fmt.Println(line)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
