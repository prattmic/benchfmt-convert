module github.com/prattmic/benchfmt-convert

go 1.17

replace golang.org/x/perf/v2 => github.com/aclements/go-perf-v2/v2 v2.0.0-20201114230402-4cc84854ceef

require (
	github.com/google/go-cmp v0.5.5
	golang.org/x/perf/v2 v2.0.0-00010101000000-000000000000
)

require golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
