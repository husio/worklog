package wlog

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"
)

// Parse given worklog.
func Parse(r io.Reader) ([]*Entry, error) {
	var entries []*Entry

	rd := bufio.NewReader(r)

	currentTask := &Task{}
	currentEntry := &Entry{}
	for {
		line, err := rd.ReadString('\n')
		switch err {
		case nil:
			// all good
		case io.EOF:
			if line == "" {
				return entries, nil
			}
		default:
			return nil, fmt.Errorf("read line: %w", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if t, err := time.Parse(TimeFormat, line); err == nil {
			if len(currentTask.Description) > 0 && len(currentEntry.Tasks) == 0 {
				currentEntry.Tasks = append(currentEntry.Tasks, currentTask)
			}
			currentTask = &Task{}
			currentEntry = &Entry{Day: t}
			entries = append(entries, currentEntry)
			continue
		}

		if potentialDuration, offset := firstWord(line); len(potentialDuration) != 0 {
			if d, err := time.ParseDuration(potentialDuration); err == nil {
				if len(currentTask.Description) > 0 && len(currentEntry.Tasks) == 0 {
					currentEntry.Tasks = append(currentEntry.Tasks, currentTask)
				}
				currentTask = &Task{Duration: d}
				currentEntry.Tasks = append(currentEntry.Tasks, currentTask)
				line = line[offset:]
			}
		}

		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "; ", ".\n")
		if len(currentTask.Description) != 0 {
			currentTask.Description += "\n"
		}
		currentTask.Description += line
	}

}

func firstWord(line string) (string, int) {
	for i, c := range line {
		if unicode.IsSpace(c) {
			return line[:i], i
		}
	}
	return line, len(line)
}

var TimeFormat = env("WORKLOG_HEADER", "# 2 Jan 2006 Monday")

type Entry struct {
	Day   time.Time
	Tasks []*Task
}

func (e *Entry) TotalDuration() time.Duration {
	var total time.Duration
	for _, t := range e.Tasks {
		total += t.Duration
	}
	return total
}

type Task struct {
	Duration    time.Duration
	Description string
}
