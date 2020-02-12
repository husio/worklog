package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/husio/worklog/wlog"
)

func cmdOpen(input io.Reader, output io.Writer, args []string) error {
	// cmdOpen is a special command becuse it ignores provided IO

	if err := ensureTodaysHeader(); err != nil {
		return fmt.Errorf("ensure header: %w", err)
	}

	// Apparently there is no perfect way to detect what is the preferred
	// terminal emulator.
	cmd := exec.Command("xterm", "-e",
		"vim",
		"+normal Gzzo",
		"+startinsert",
		worklogPath(),
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}

func ensureTodaysHeader() error {
	header := time.Now().Format(wlog.TimeFormat)
	wpath := worklogPath()
	fd, err := os.OpenFile(wpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("open %q file: %w", wpath, err)
	}
	defer fd.Close()

	rd := bufio.NewReader(fd)
	for {
		switch line, err := rd.ReadString('\n'); {
		case errors.Is(err, nil):
			if line[:len(line)-1] == header {
				// Today's header found.
				return nil
			}
		case errors.Is(err, io.EOF):
			if _, err := fmt.Fprintf(fd, "\n%s\n", header); err != nil {
				return fmt.Errorf("write header: %w", err)
			}
			return nil
		default:
			return fmt.Errorf("read worklog: %w", err)
		}
	}
}
