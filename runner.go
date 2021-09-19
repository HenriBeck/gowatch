package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Runner interface {
	Run(args []string)

	Stop()
}

type runner struct {
	builder Builder

	command *exec.Cmd
}

func NewRunner(builder Builder) Runner {
	return &runner{
		builder: builder,
	}
}

func (runner *runner) Stop() {
	if runner.command == nil || runner.command.Process == nil {
		return
	}

	done := make(chan bool)
	go func() {
		runner.command.Wait()
		close(done)
	}()

	// Trying a "soft" kill first
	if err := runner.command.Process.Signal(os.Interrupt); err != nil {
		log.Println("Failed to interrupt previous process", err)
	}

	// Wait for our process to die before we return or hard kill after 3 sec
	select {
	case <-time.After(3 * time.Second):
		if err := runner.command.Process.Kill(); err != nil {
			log.Println("failed to kill: ", err)
		}

	case <-done:
		break
	}

	runner.command = nil
}

func (runner *runner) Run(args []string) {
	if runner.command != nil {
		runner.Stop()
	}

	err := runner.builder.Build(args)
	if err != nil {
		fmt.Println("Error while building application")
	}

	runner.command = exec.Command(
		runner.builder.GetExecutable(),
	)
	runner.command.Stderr = os.Stderr
	runner.command.Stdout = os.Stdout

	err = runner.command.Start()
	if err != nil {
		fmt.Println("Error while starting the application")
	}
}
