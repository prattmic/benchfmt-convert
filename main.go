package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/perf/v2/benchfmt"
)

func main() {
	r := csv.NewReader(os.Stdin)

	var header []string
	var results []benchfmt.Result
	for i := 0; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error reading: %v\n", err)
			return
		}

		if i == 0 {
			header = record
			continue
		}

		r := benchfmt.Result{
			FullName: []byte(record[0]),
			Iters:    1,
			Values:   make([]benchfmt.Value, 0),
		}
		// Skip label and count.
		for j := 2; j < len(record); j++ {
			val, err := strconv.ParseFloat(record[j], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", record[j], err)
				return
			}

			r.Values = append(r.Values, benchfmt.Value{
				Value: val,
				Unit:  header[j],
			})
		}

		results = append(results, r)
	}

	w := benchfmt.NewWriter(os.Stdout)

	for _, r := range results {
		if err := w.Write(&r); err != nil {
			fmt.Fprintf(os.Stderr, "error writing: %v\n", err)
		}
	}
}
