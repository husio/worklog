package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/husio/worklog/wlog"
)

func cmdFilter(input io.Reader, output io.Writer, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: fmt <month-name>")
	}

	m, ok := months[args[0]]
	if !ok {
		return fmt.Errorf("invalid month: %q", args[0])
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
