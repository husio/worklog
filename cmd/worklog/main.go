package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [<flags>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nAvailable commands are:\n\t%s\n", strings.Join(availableCmds(), "\n\t"))
		fmt.Fprintf(os.Stderr, "Run '%s <command> -help' to learn more about each command.\n", os.Args[0])
		os.Exit(2)
	}
	run, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "\nAvailable commands are:\n\t%s\n", strings.Join(availableCmds(), "\n\t"))
		os.Exit(2)
	}

	input, err := logReader(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()

	// Skip first two arguments. Second argument is the command name that
	// we just consumed.
	if err := run(input, os.Stdout, os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var commands = map[string]func(input io.Reader, output io.Writer, args []string) error{
	"filter":  cmdFilter,
	"fmt":     cmdFmt,
	"summary": cmdSummary,
	"add":     cmdAdd,
}

func availableCmds() []string {
	available := make([]string, 0, len(commands))
	for name := range commands {
		available = append(available, name)
	}
	sort.Strings(available)
	return available
}

func logReader(r io.ReadCloser) (io.ReadCloser, error) {
	if s, ok := r.(interface{ Stat() (os.FileInfo, error) }); ok {
		if info, err := s.Stat(); err == nil {
			if isPipe := (info.Mode() & os.ModeCharDevice) == 0; isPipe {
				return r, nil
			}
		}
	}
	path, ok := os.LookupEnv("WORKLOG")
	if !ok {
		path = filepath.Join(os.Getenv("HOME"), "/worklog.txt")
	}
	fd, err := os.Open(path)
	return fd, err
}
