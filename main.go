package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"sync"
)

var (
	maxprocs int
	number   int
	command  string
)

func init() {
	flag.IntVar(&maxprocs, "n", 1, "maxprocs")
	flag.IntVar(&number, "P", 1, "number")
	flag.StringVar(&command, "C", "echo", "command to exec")
}

func main() {
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		os.Exit(1)
	}

	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanWords)

	cmderch := make(chan *cmder, maxprocs)

	go func() {
		var args []string
		for s.Scan() {
			args = append(args, s.Text())
			if len(args) >= number {
				cmderch <- cmd(command, args...)
				args = []string{}
			}
		}
		if len(args) > 0 {
			cmderch <- cmd(command, args...)
		}

		close(cmderch)
	}()

	var wg sync.WaitGroup
	for c := range cmderch {
		wg.Add(1)
		go func(c *cmder) {
			defer wg.Done()
			c.exec()
		}(c)
	}

	wg.Wait()
}

type cmder struct {
	command string
	args    []string
}

func cmd(command string, args ...string) *cmder { return &cmder{command: command, args: args} }

func (c *cmder) exec() {
	ec := exec.Command(c.command, c.args...)
	ec.Stderr = os.Stderr
	ec.Stdout = os.Stdout
	ec.Run()
}
