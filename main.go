package main

import (
	"fmt"
	"os"

	"golang.org/x/perf/v2/benchfmt"
)

func main() {
	w := benchfmt.NewWriter(os.Stdout)

	r := benchfmt.Result{
		FullName: []byte("TestBench"),
		Iters:    1,
		Values:   []benchfmt.Value{
			{
				Value: 42,
				Unit:  "ns",
			},
		},
	}
	if err := w.Write(&r); err != nil {
		fmt.Fprintf(os.Stderr, "error writing: %v\n", err)
	}
}
