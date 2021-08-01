package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/husio/worklog/wlog"
)

func cmdPush(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("push", flag.ContinueOnError)
	urlFl := fl.String("url", "", "Worklog storage URL.")
	tokenFl := fl.String("token", "", "Worklog storage write token,")
	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("flag parse: %w", err)
	}
	if *urlFl == "" {
		return fmt.Errorf("\"url\" not provided")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var body bytes.Buffer
	if entries, err := wlog.Parse(input); err != nil {
		return fmt.Errorf("parse log: %s", err)
	} else if err := wlog.ToText(&body, entries); err != nil {
		return fmt.Errorf("format to text: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", *urlFl, &body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("write-token", *tokenFl)
	req.Header.Set("content-type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http PUT: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("forbidden: invalid push token")
	}
	if resp.StatusCode > 204 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1e5))
		return fmt.Errorf("http response: %d %s", resp.StatusCode, string(b))
	}

	return nil
}
