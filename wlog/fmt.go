package wlog

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ToText conver given entries into text format.
func ToText(w io.Writer, entries []*Entry) error {
	for _, e := range entries {
		// Ignore empty days.
		if e.TotalDuration() == 0 {
			continue
		}
		if _, err := fmt.Fprintln(w, e.Day.Format(TimeFormat)); err != nil {
			return fmt.Errorf("write entry info: %w", err)
		}
		for _, t := range e.Tasks {
			hours := t.Duration / time.Hour
			lines := strings.Split(t.Description, "\n")
			if _, err := fmt.Fprintf(w, "%dh %s\n", hours, strings.TrimSpace(lines[0])); err != nil {
				return fmt.Errorf("write task info: %w", err)
			}
			indent := strings.Repeat(" ", 3)
			if hours > 9 {
				indent = strings.Repeat(" ", 4)
			}
			for _, line := range lines[1:] {
				if _, err := io.WriteString(w, indent+strings.TrimSpace(line)+"\n"); err != nil {
					return fmt.Errorf("write task info: %w", err)
				}
			}
		}
		fmt.Fprint(w, "\n")
	}
	return nil
}
