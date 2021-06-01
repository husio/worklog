package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/husio/worklog/wlog"
)

func cmdSummary(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("summary", flag.ContinueOnError)
	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("flag parse: %w", err)
	}

	entries, err := wlog.Parse(input)
	if err != nil {
		return fmt.Errorf("parse log: %s", err)
	}
	var total time.Duration
	for _, e := range entries {
		for _, t := range e.Tasks {
			total += t.Duration
		}
	}
	if _, err := fmt.Fprintln(output, "total", total); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, "days", float32(total/time.Hour/8)); err != nil {
		return err
	}
	return nil
}
