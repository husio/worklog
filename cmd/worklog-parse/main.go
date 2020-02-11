package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

func main() {
	monthFl := flag.String("month", "", "Filter to only entries from a given month")
	summaryFl := flag.Bool("summary", false, "Display summary only")
	flag.Parse()

	var filterMonth time.Month = 0 // Default to not set.
	if len(*monthFl) != 0 {
		m, ok := months[*monthFl]
		if !ok {
			panic("invalid month value")
		}
		filterMonth = m
	}

	fd, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	rd := bufio.NewReader(fd)

	var (
		entries []*Entry
		current *Entry = &Entry{}
	)

consume:
	for {
		line, err := rd.ReadString('\n')
		switch err {
		case nil:
			// all good
		case io.EOF:
			break consume
		default:
			panic(err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if t, err := time.Parse(timeFormat, line); err == nil {
			t.Month()
			current = &Entry{Day: t}
			if filterMonth == 0 || t.Month() == filterMonth {
				entries = append(entries, current)
			}
			continue
		}

		if potentialDuration, offset := firstWord(line); len(potentialDuration) != 0 {
			if d, err := time.ParseDuration(potentialDuration); err == nil {
				current.Duration += d
				line = line[offset:]
			}
		}

		if len(current.Description) != 0 {
			current.Description += "\n"
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "; ", ".\n")
		current.Description += line
	}

	if *summaryFl {
		var total time.Duration
		for _, e := range entries {
			total += e.Duration
		}
		fmt.Println("total", total)
		fmt.Println("days", float32(total/time.Hour/8))
		return
	}

	b, _ := json.MarshalIndent(entries, "", "  ")
	os.Stdout.Write(b)
}

const timeFormat = "# 2 Jan 2006 Monday"

type Entry struct {
	Day         time.Time
	Duration    time.Duration
	Description string
}

func firstWord(line string) (string, int) {
	for i, c := range line {
		if unicode.IsSpace(c) {
			return line[:i], i
		}
	}
	return line, len(line)
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
