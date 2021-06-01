package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/husio/worklog/wlog"
)

func cmdFilter(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("filter", flag.ContinueOnError)
	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("flag parse: %w", err)
	}

	if len(fl.Args()) != 1 {
		return errors.New("usage: fmt <month-name>")
	}

	m, ok := months[fl.Args()[0]]
	if !ok {
		return fmt.Errorf("invalid month: %q", fl.Args()[0])
	}

	entries, err := wlog.Parse(input)
	if err != nil {
		return fmt.Errorf("parse log: %s", err)
	}

	all := entries
	entries = nil
	for _, e := range all {
		if e.Day.Month() == m {
			entries = append(entries, e)
		}
	}

	if err := wlog.ToText(output, entries); err != nil {
		return fmt.Errorf("format to text: %w", err)
	}
	return nil
}

var months = map[string]time.Month{
	"January":  time.January,
	"Jan":      time.January,
	"February": time.February,
	"Feb":      time.February,
	"March":    time.March,
	"Mar":      time.March,
	"April":    time.April,
	"Apr":      time.April,
	"May":      time.May,
	"June":     time.June,
	"Jun":      time.June,
	"July":     time.July,
	"Jul":      time.July,
	"August":   time.August,
	"Aug":      time.August,
	"Septembe": time.September,
	"Sep":      time.September,
	"October":  time.October,
	"Oct":      time.October,
	"November": time.November,
	"Nov":      time.November,
	"December": time.December,
	"Dec":      time.December,
}
