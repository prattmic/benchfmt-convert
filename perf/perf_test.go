package perf

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/perf/v2/benchfmt"
)

func TestPerf(t *testing.T) {
	tests := []struct{
		name  string
		line  string
		match bool
		want  benchfmt.Result
	}{
		{
			name:  "unitless",
			line:  "          1,291,018      cycles                    #    1.167 GHz",
			match: true,
			want:  benchfmt.Result{
				FullName: []byte("BenchmarkCycles"),
				Iters:    1,
				Values:   []benchfmt.Value{
					{
						Value: 1291018,
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

func TestCapitalize(t *testing.T) {
	tests := []struct{
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "one",
			input: "a",
			want:  "A",
		},
		{
			name:  "multi",
			input: "abcd",
			want:  "Abcd",
		},
		{
			name:  "unicode",
			input: "ā",
			want:  "Ā",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := capitalize(tc.input)
			if got != tc.want {
				t.Errorf("capitalize(%q) got %q want %q", tc.input, got, tc.want)
			}
		})
	}
}
