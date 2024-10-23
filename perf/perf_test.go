package perf

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/perf/v2/benchfmt"
)

func TestPerf(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		match bool
		want  benchfmt.Result
	}{
		{
			name:  "unitless",
			line:  "          1,291,018      cycles                    #    1.167 GHz",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Cycles"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 1291018,
						Unit:  "val",
					},
				},
			},
		},
		{
			name:  "unitless-suffix",
			line:  "          1,291,018      cycles:u                    #    1.167 GHz",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Cycles:u"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 1291018,
						Unit:  "val",
					},
				},
			},
		},
		{
			name:  "unitless-capital",
			line:  "     2,651,008,697      L1-dcache-loads",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("L1-dcache-loads"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 2651008697,
						Unit:  "val",
					},
				},
			},
		},
		{
			name:  "unit",
			line:  "               1.11 msec task-clock                #    0.001 CPUs utilized",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Task-clock"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 0.00111,
						Unit:  "sec",
					},
				},
			},
		},
		{
			name:  "unit-suffix",
			line:  "               1.11 msec task-clock:u                #    0.001 CPUs utilized",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Task-clock:u"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 0.00111,
						Unit:  "sec",
					},
				},
			},
		},
		{
			name:  "wall",
			line:  "      1656.917143299 seconds time elapsed",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("Wall-time"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 1656.917143299,
						Unit:  "sec",
					},
				},
			},
		},
		{
			name:  "user",
			line:  "      1656.917143299 seconds user",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("User-time"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 1656.917143299,
						Unit:  "sec",
					},
				},
			},
		},
		{
			name:  "system",
			line:  "      1656.917143299 seconds sys",
			match: true,
			want: benchfmt.Result{
				FullName: []byte("System-time"),
				Iters:    1,
				Values: []benchfmt.Value{
					{
						Value: 1656.917143299,
						Unit:  "sec",
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
	tests := []struct {
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
