package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Runner interface {
	Run(args []string)

	Stop()
}

type runner struct {
	builder Builder

	command *exec.Cmd

	prefix string
}

func NewRunner(builder Builder) Runner {
	return &runner{
		builder: builder,
		prefix:  lipgloss.NewStyle().Foreground(lipgloss.Color("#48d597")).Render("[gowatch]"),
	}
}

func (runner *runner) Stop() {
	if runner.command == nil || runner.command.Process == nil {
		return
	}

	fmt.Printf("%s Stopping previous process\n", runner.prefix)

	done := make(chan bool)
	go func() {
		runner.command.Wait()
		close(done)
	}()

	// Trying a "soft" kill first
	if err := runner.command.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("%s Failed to interrupt previous process: %v\n", runner.prefix, err)
	}

	// Wait for our process to die before we return or hard kill after 3 sec
	select {
	case <-time.After(3 * time.Second):
		if err := runner.command.Process.Kill(); err != nil {
			fmt.Printf("%s Failed to kill previous process: %v\n", runner.prefix, err)
		}

	case <-done:
		break
	}

	runner.command = nil
}

func (runner *runner) Run(args []string) {
	// First we stop the running command if it exists
	if runner.command != nil {
		runner.Stop()
	}

	fmt.Printf("%s Building...\n", runner.prefix)

	err := runner.builder.Build(args)
	if err != nil {
		fmt.Printf("%s Error while building application:\n", runner.prefix)
		fmt.Println(err)
		return
	}

	fmt.Printf("%s Build finished\n", runner.prefix)

	runner.command = exec.Command(
		runner.builder.GetExecutable(),
	)
	runner.command.Stderr = os.Stderr
	runner.command.Stdout = os.Stdout

	err = runner.command.Start()
	if err != nil {
		fmt.Printf("%s Process failed to start: %v\n", runner.prefix, err)
	}
}
