package wlog

import (
	"fmt"
	"io"
	"time"
)

// ToText conver given entries into text format.
func ToText(w io.Writer, entries []*Entry) error {
	for _, e := range entries {
		if _, err := fmt.Fprintln(w, e.Day.Format(timeFormat)); err != nil {
			return fmt.Errorf("write entry info: %w", err)
		}
		for _, t := range e.Tasks {
			hours := t.Duration / time.Hour
			if _, err := fmt.Fprintf(w, "%dh %s\n", hours, t.Description); err != nil {
				return fmt.Errorf("write task info: %w", err)
			}
		}
		fmt.Fprint(w, "\n")
	}
	return nil
}
