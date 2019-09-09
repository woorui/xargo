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

func readArgs(maxArgs int, argsCh chan []string) {
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
	close(argsCh)
}

func buildCmd(cmdsCh chan CmdWithArgs, name string, args []string) {
	cwa := CmdWithArgs{name: name, args: args}
	// fmt.Println(cwa)
	cmdsCh <- cwa
}

func buildGroupCmd(n int, maxProcs int, cmdsCh chan CmdWithArgs, cmdsGroupCh chan []CmdWithArgs) {
	cmds := []CmdWithArgs{}
	for i := 0; i < n; i++ {
		cwa := <-cmdsCh
		if len(cmds) == maxProcs {
			cmdsGroupCh <- cmds
			cmds = []CmdWithArgs{}

		}
		cmds = append(cmds, cwa)
	}
	if len(cmds) != 0 {
		cmdsGroupCh <- cmds
	}
	close(cmdsGroupCh)
}

func execGroupCmd(wg *sync.WaitGroup, cmdsGroup []CmdWithArgs) {
	for _, c := range cmdsGroup {
		wg.Add(1)
		go func(c CmdWithArgs) {
			cmd := exec.Command(c.name, c.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			wg.Done()
		}(c)
	}
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

	argsCh := make(chan []string)

	// Read From stdin
	go readArgs(*maxArgs, argsCh)

	cmdsCh := make(chan CmdWithArgs, 2)
	cmdsGroupCh := make(chan []CmdWithArgs)

	// Build name and args used by command
	n := 0
	for args := range argsCh {
		n++
		go buildCmd(cmdsCh, *bin, args)
	}

	// Group the command for execute asynchronously
	go buildGroupCmd(n, *maxProcs, cmdsCh, cmdsGroupCh)

	// execute the command Grouped
	var wg sync.WaitGroup
	for cg := range cmdsGroupCh {
		execGroupCmd(&wg, cg)
	}
	wg.Wait()
}
