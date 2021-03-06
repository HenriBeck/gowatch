package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Options struct {
	Interval time.Duration

	Path string
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	WatchForChanges(&Options{
		Interval: 5 * time.Second,
		Path:     wd,
	})
}

func WatchForChanges(options *Options) {
	lastModificationTime := time.Now()
	ticker := time.NewTicker(options.Interval)

	builder := NewBuilder()
	defer builder.Clean()

	runner := NewRunner(builder)

	// Initially start the runner
	runner.Run(os.Args[1:])
	defer runner.Stop()

	cancel := make(chan os.Signal, 2)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			hasFileChanges := scanChanges(options.Path, lastModificationTime)
			if hasFileChanges {
				lastModificationTime = time.Now()
				runner.Run(os.Args[1:])
			}

		case <-cancel:
			// Print an empty live because a SIGTERM doesn't include one
			// and it would write our logs in the next line
			fmt.Println()
			return
		}
	}
}

func scanChanges(
	watchPath string,
	time time.Time,
) bool {
	var hasFileChanges = errors.New("file changes")

	err := filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Skip any hidden directory
			if filepath.Base(path)[0] == '.' {
				return filepath.SkipDir
			}

			return nil
		}

		// Skip any files which are not go files
		if filepath.Ext(path) != ".go" {
			return nil
		}

		if info.ModTime().After(time) {
			return hasFileChanges
		}

		return nil
	})

	return err == hasFileChanges
}
