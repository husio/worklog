package main

import (
	"fmt"
	"io"
	"os"

	"github.com/husio/worklog/wlog"
)

func cmdFmt(input io.Reader, output io.Writer, args []string) error {
	entries, err := wlog.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse log: %s\n", err)
		os.Exit(1)
	}

	if err := wlog.ToText(output, entries); err != nil {
		return fmt.Errorf("format to text: %w", err)
	}
	return nil
}
