package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"
)

type Builder interface {
	Build(args []string) error

	Clean()

	GetExecutable() string
}

type builder struct {
	executable string
	mutex      sync.Mutex
}

func NewBuilder() Builder {
	return &builder{
		executable: path.Join(
			os.TempDir(),
			strconv.FormatInt(time.Now().Unix(), 10),
		),
	}
}

func (builder *builder) GetExecutable() string {
	return builder.executable
}

func (builder *builder) Build(args []string) error {
	// Only allow one build at a time
	builder.mutex.Lock()
	defer builder.mutex.Unlock()

	commandArgs := []string{
		"build",
		"-o",
		builder.executable,
	}
	commandArgs = append(commandArgs, args...)

	// nolint:gosec
	cmd := exec.Command("go", commandArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: \n %s", err, output)
	}

	return nil
}

func (builder *builder) Clean() {
	os.Remove(builder.executable)
}
