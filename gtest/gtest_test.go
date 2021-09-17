package gtest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/perf/v2/benchfmt"
)

///^BM_.*\s+[0-9]+\s+.*\s+[0-9]+\s+.*\s+[0-9]+$/ {

func TestGTest(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		match bool
		want  benchfmt.Result
	}{
		{
			name:  "typical",
			line:  "BM_Stat/64/real_time       16770 ns        16593 ns        42186",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Stat/64/real_time"),
				Iters:    42186,
				Values: []benchfmt.Value{
					{
						Value: 16770,
						Unit:  "ns",
					},
					{
						Value: 16593,
						Unit:  "cpu-ns",
					},
				},
			},
		},
		{
			name:  "trailing whitespace",
			line:  "BM_Stat/64/real_time       16770 ns        16593 ns        42186  ",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Stat/64/real_time"),
				Iters:    42186,
				Values: []benchfmt.Value{
					{
						Value: 16770,
						Unit:  "ns",
					},
					{
						Value: 16593,
						Unit:  "cpu-ns",
					},
				},
			},
		},
		{
			name:  "no units",
			line:  "BM_LargeConsistent                 4.69           4.69   500000000  ",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("LargeConsistent"),
				Iters:    500000000,
				Values: []benchfmt.Value{
					{
						Value: 4.69,
						Unit:  "ns",
					},
					{
						Value: 4.69,
						Unit:  "cpu-ns",
					},
				},
			},
		},
		{
			name:  "user counters",
			line:  "BM_Read/1/real_time              3820 ns         2757 ns       184963 bytes_per_second=255.676k/s",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Read/1/real_time"),
				Iters:    184963,
				Values: []benchfmt.Value{
					{
						Value: 3820,
						Unit:  "ns",
					},
					{
						Value: 2757,
						Unit:  "cpu-ns",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, ok := Line(tc.line)
			if ok != tc.match {
				t.Errorf("Line(%q) ok got %v want %v", tc.line, ok, tc.match)
			}
			if diff := cmp.Diff(r, tc.want, cmpopts.IgnoreUnexported(benchfmt.Result{})); diff != "" {
				t.Errorf("Line(%q) mismatch (-want +got):\n%s", tc.line, diff)
			}
		})
	}
}
