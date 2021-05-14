package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/prattmic/benchfmt-convert/perf"
	"golang.org/x/perf/v2/benchfmt"
)

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

		r, ok := perf.Line(line)
		if ok {
			if err := w.Write(&r); err != nil {
				return fmt.Errorf("error writing %+v: %v", r, err)
			}
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
