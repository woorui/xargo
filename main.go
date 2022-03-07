package main

import (
	"context"
	"flag"
	"os"
	"os/exec"

	"github.com/woorui/xargo/xargs"
)

var (
	maxprocs int
	number   int
	command  string
)

func init() {
	flag.IntVar(&maxprocs, "P", 3, "maxprocs")
	flag.IntVar(&number, "n", 3, "number")
	flag.StringVar(&command, "C", "echo", "command to exec")
}

func main() {
	flag.Parse()

	xargs.New(context.Background(), buildCmd, command, maxprocs, number).Execute(os.Stdin)
}

type cmder struct {
	command string
	args    []string
}

func buildCmd(command string, args ...string) xargs.Cmder {
	return &cmder{command: command, args: args}
}

func (c *cmder) Exec() error {
	ec := exec.Command(c.command, c.args...)
	ec.Stderr = os.Stderr
	ec.Stdout = os.Stdout
	return ec.Run()
}
