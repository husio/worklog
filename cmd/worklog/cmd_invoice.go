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

ItemDescription   = Software development: 01.01.2020 – 30.01.2020
ItemHours         = 123
ItemRate          = 100

InvoiceNumber     = 2020-01-1
InvoiceDate       = 2020-01-30

NoTax             = false

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

	if err := populateFrom(&tctx, fd); err != nil {
		return fmt.Errorf("cannot read configuration: %w", err)
	}

	tctx.ItemTotal = tctx.ItemHours * tctx.ItemRate

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

func populateFrom(s interface{}, r io.Reader) error {
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
