# gowatch

`gowatch` is a simple tool that builds a Go application, runs the resulting binary, and repeats this process if any go files change.  

You should only use the tool itself during development, and it will not restart the program if it exits abnormally.

## Installation

To install the tool, use the `go install` command:

```sh
go install github.com/HenriBeck/gowatch@latest
```

## Usage

To build and run the Go application in your current working directory, simply run:
```sh
gowatch
```

`gowatch` will then watch any Go files in your current working directory and rebuild and restart your application when any file changes.



In case your Go programm (`main.go` file) is not in your current working directory, you can pass the path to the directory to `gowatch`, similar to this:
```sh
gowatch ./api
```

This will build the application in the `api` folder but still watch all the files from your current working directory.

### Go Build Arguments

If your application relies on other `go build` arguments, you can pass them to `gowatch`, which will automatically be forwarded.  
For example, you can pass a tags list:

```sh
gowatch --tags api .
```

> Don't pass the `-o` option as this needs to be set from `gowatch` itself.
