# benchfmt-convert

This repo contains tools for converting various formats to benchfmt for use with
benchstat.

Most useful is
[`github.com/prattmic/benchfmt-convert/cmd/normalize-benchfmt`](https://pkg.go.dev/github.com/prattmic/benchfmt-convert/cmd/normalize-benchfmt),
which will discover benchmark results in many different formats from an input
log file and output their benchfmt equivalents.

Supported formats:

 * `perf stat`
 * googletest benchmarks
