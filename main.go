package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"os/exec"
	"sync"
)

// CmdWithArgs contain cmd name and args
type CmdWithArgs struct {
	name string
	args []string
}

func isDeLimited(b byte) bool {
	return b == 10 || b == 32 || b == 9
}

func readArgs(maxArgs int) (argsCh chan []string) {
	defer close(argsCh)
	reader := bufio.NewReader(os.Stdin)
	var args []string
	arg := ""
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			argsCh <- args
			break
		}

		if !isDeLimited(b) {
			arg = arg + string(b)
		} else {
			args = append(args, string(arg))
			arg = ""
			if len(args) == maxArgs {
				argsCh <- args
				args = []string{}
			}
		}
	}
	return argsCh
}

func buildCmd(bin string, maxProcs int, argsCh chan []string) chan CmdWithArgs {
	cmdsCh := make(chan CmdWithArgs, maxProcs)
	defer close(cmdsCh)
	go func() {
		for args := range argsCh {
			cwa := CmdWithArgs{name: bin, args: args}
			cmdsCh <- cwa
		}
	}()
	return cmdsCh
}

func execCmd(cmdsCh chan CmdWithArgs) {
	var wg sync.WaitGroup
	for c := range cmdsCh {
		wg.Add(1)
		go func(c CmdWithArgs) {
			cmd := exec.Command(c.name, c.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			wg.Done()
		}(c)
	}
	wg.Wait()
}

func main() {
	// Parse flag
	maxArgs := flag.Int("n", 1, "max-args, default 0, It mean no limit")
	maxProcs := flag.Int("P", 1, "max-procs, default 0, It mean no limit")
	bin := flag.String("bin", "echo", "command to exec, default echo")
	flag.Parse()

	// Check if stdin is from a terminal
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		os.Exit(1)
	}

	// Read From stdin
	argsCh := readArgs(*maxArgs)

	// Build name and args used by command
	cmdsCh := buildCmd(*bin, *maxProcs, argsCh)

	// execute the command Grouped
	execCmd(cmdsCh)
}
