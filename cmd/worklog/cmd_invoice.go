package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/husio/worklog/wlog"
)

var (
	//go:embed cmd_invoice.html
	rawTmpl string
	tmpl    = template.Must(template.New("").Parse(rawTmpl))
)

type TemplateContext struct {
	Debtor          string
	InvoiceNumber   string
	InvoiceDate     string
	ToCompany       string
	ToAddress       string
	ToCo            string
	ToVATID         string
	FromName        string
	FromAddress     string
	FromCountry     string
	FromVATID       string
	FromEmail       string
	PaymentName     string
	PaymentIBAN     string
	PaymentBIC      string
	PaymentBankName string
	ItemDescription string
	ItemHours       int
	ItemRate        int
	ItemTotal       int
	NoTax           bool
	SignatureBase64 string
}

func cmdInvoice(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("invoice", flag.ContinueOnError)
	confFl := fl.String("c", "config.txt", "Path to the configuration file.")
	outFl := fl.String("o", "", "Output file. Stdout if not given.")
	exConfFl := fl.Bool("g", false, "Generate an example configuration file.")
	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("flag parse: %w", err)
	}

	if *exConfFl {
		fmt.Fprint(output, `
Debtor            =

ToCompany         =
ToAddress         =
ToCo              =
ToVATID           =

FromName          =
FromAddress       =
FromCountry       =
FromVATID         =
FromEmail         =

PaymentName       =
PaymentIBAN       =
PaymentBIC        =
PaymentBankName   =

ItemRate          = 100
ItemDescription   = Software development.

# Below entries are generated from the worklog if not provided.
ItemHours         =
InvoiceNumber     =
InvoiceDate       =

NoTax             = false

# base64 encoded PNG image.
SignatureBase64   =
		`)
		return nil
	}

	var tctx TemplateContext

	fd, err := os.Open(*confFl)
	if err != nil {
		return fmt.Errorf("cannot open %q: %w", *confFl, err)
	}
	defer fd.Close()

	if err := populateFromConfig(&tctx, fd); err != nil {
		return fmt.Errorf("cannot read configuration: %w", err)
	}

	entries, err := wlog.Parse(input)
	if err != nil {
		return fmt.Errorf("parse log: %s", err)
	}
	if err := populateFromLog(&tctx, entries); err != nil {
		return fmt.Errorf("cannot interpred log: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, tctx); err != nil {
		return fmt.Errorf("cannot render template: %w", err)
	}
	if *outFl == "" {
		if _, err = b.WriteTo(output); err != nil {
			return fmt.Errorf("cannot write to stdout: %w", err)
		}
	} else {
		if err := ioutil.WriteFile(*outFl, b.Bytes(), 0644); err != nil {
			return fmt.Errorf("cannot write to %q: %w", *outFl, err)
		}
	}
	return nil
}

func populateFromLog(c *TemplateContext, entries []*wlog.Entry) error {

	if c.ItemHours == 0 {
		var total time.Duration
		for _, e := range entries {
			total += e.TotalDuration()
		}
		c.ItemHours = int(total / time.Hour)
	}

	c.ItemTotal = c.ItemHours * c.ItemRate

	last := entries[len(entries)-1]
	if c.InvoiceDate == "" {
		c.InvoiceDate = last.Day.Format("2006-01-02")
	}
	if c.InvoiceNumber == "" {
		c.InvoiceNumber = last.Day.Format("2006-01-") + "-01"
	}

	first := entries[0]
	c.ItemDescription += fmt.Sprintf("<br><em>(%s - %s)</em>",
		first.Day.Format("02.01.2006"),
		last.Day.Format("02.01.2006"),
	)

	return nil
}

func populateFromConfig(s interface{}, r io.Reader) error {
	v := reflect.ValueOf(s).Elem()

	rd := bufio.NewReader(r)
	for {
		line, err := rd.ReadString('\n')
		switch err {
		case nil:
			// All good.
		case io.EOF:
			return nil
		default:
			return fmt.Errorf("read string: %w", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Ignore comments.
		if line[0] == '#' {
			continue
		}

		chunks := strings.SplitN(line, "=", 2)
		name := strings.TrimSpace(chunks[0])

		value := strings.TrimSpace(chunks[1])
		value = strings.ReplaceAll(value, "\\n", "\n")

		field := v.FieldByName(name)
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("cannot set %q field value", name)
		}

		switch field.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32:
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("value of %q is not a valid number: %w", name, err)
			}
			field.SetInt(n)
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("%q field value is not a valid boolean: %w", name, err)
			}
			field.SetBool(v)
		}

	}
}
