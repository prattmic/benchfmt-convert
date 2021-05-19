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
