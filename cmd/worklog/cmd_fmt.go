package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/husio/worklog/wlog"
)

func cmdFmt(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("fmt", flag.ContinueOnError)
	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("flag parse: %w", err)
	}

	var format string
	switch len(fl.Args()) {
	case 0:
		format = "txt"
	case 1:
		format = fl.Args()[0]
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
	case "html":
		var extended []*wlog.Entry
		if len(entries) > 0 {
			var i int
			for t := entries[0].Day; !t.After(entries[len(entries)-1].Day); t = t.Add(time.Hour * 24) {
				if i < len(entries) && entries[i].Day.Equal(t) {
					extended = append(extended, entries[i])
					i++
				} else {
					extended = append(extended, &wlog.Entry{Day: t})
				}
			}
		}
		sort.Slice(extended, func(i, j int) bool {
			return extended[i].Day.After(extended[j].Day)
		})

		var byMonth [][]*wlog.Entry
		if len(extended) > 0 {
			current := []*wlog.Entry{extended[0]}
			for _, e := range extended[1:] {
				if e.Day.Month() == current[0].Day.Month() {
					current = append(current, e)
				} else {
					byMonth = append(byMonth, current)
					current = []*wlog.Entry{e}
				}
			}
			byMonth = append(byMonth, current)
		}
		context := struct {
			Entries [][]*wlog.Entry
		}{
			Entries: byMonth,
		}
		var b bytes.Buffer
		if err := fmtTmpl.Execute(&b, context); err != nil {
			return fmt.Errorf("render template: %w", err)
		}
		if _, err := b.WriteTo(output); err != nil {
			return fmt.Errorf("write to output: %w", err)
		}
		return nil
	default:
		return errors.New("valid formats are text, json, csv")
	}
}

//go:embed cmd_fmt.html
var htmlFmtTemplate string

var fmtTmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"narrowhours": func(d time.Duration) string {
		hours := d / time.Hour
		if hours == 0 {
			return ""
		}
		return fmt.Sprintf("%dh", hours)
	},
	"hoursduration": func(entries []*wlog.Entry) uint {
		var total time.Duration
		for _, e := range entries {
			total += e.TotalDuration()
		}
		return uint(total.Hours())
	},
}).Parse(htmlFmtTemplate))
