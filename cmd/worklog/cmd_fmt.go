package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/husio/worklog/wlog"
)

func cmdFmt(input io.Reader, output io.Writer, args []string) error {
	var format string
	switch len(args) {
	case 0:
		format = "txt"
	case 1:
		format = args[0]
	default:
		return fmt.Errorf("usage: fmt [<format>]")
	}

	entries, err := wlog.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse log: %s\n", err)
		os.Exit(1)
	}

	switch format {
	case "text", "txt":
		if err := wlog.ToText(output, entries); err != nil {
			return fmt.Errorf("format to text: %w", err)
		}
		return nil
	case "json":
		b, err := json.MarshalIndent(entries, "", "\t")
		if err != nil {
			return fmt.Errorf("serialize: %w", err)
		}
		if _, err := output.Write(b); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	case "csv":
		wr := csv.NewWriter(output)
		defer wr.Flush()

		if err := wr.Write([]string{"day", "hours", "description"}); err != nil {
			return fmt.Errorf("write header: %w", err)
		}
		for _, e := range entries {
			for _, t := range e.Tasks {
				err := wr.Write([]string{
					e.Day.Format("2/01/2006"),
					fmt.Sprint(t.Duration.Hours()),
					t.Description,
				})
				if err != nil {
					return fmt.Errorf("write entry: %w", err)
				}
			}
		}
		return nil
	default:
		return errors.New("valid formats are text, json, csv")
	}
}
