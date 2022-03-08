package xargs

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type testcmder struct {
	command string
	args    []string
}

func buildTestCmd(command string, args ...string) Cmder {
	return &testcmder{command: command, args: args}
}

func (c *testcmder) Exec() error {
	ec := exec.Command(c.command, c.args...)
	return ec.Run()
}

func Test_Xargs(t *testing.T) {
	ctx := context.Background()

	reader := strings.NewReader("aa bb cc dd cc")

	New(ctx, buildTestCmd, "echo", 2, 2).Execute(reader)
}

type errcmder struct {
	command string
	args    []string
}

func buildErrCmd(command string, args ...string) Cmder {
	return &errcmder{command: command, args: args}
}

func (c *errcmder) Exec() error {
	return errors.New("This is an error")
}

func Test_Err_Xargs(t *testing.T) {
	ctx := context.Background()

	reader := strings.NewReader("aa bb cc dd cc")

	err := New(ctx, buildErrCmd, "echo", 2, 2).Execute(reader)
	if err == nil {
		t.Fail()
	}
}

type slowcmder struct {
	command string
	args    []string
}

func buildSlowCmd(command string, args ...string) Cmder {
	return &slowcmder{command: command, args: args}
}

func (c *slowcmder) Exec() error {
	time.Sleep(time.Second * 3)
	ec := exec.Command(c.command, c.args...)
	return ec.Run()
}

func Test_Slow_Xargs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := strings.NewReader("aa bb cc dd cc")

	err := New(ctx, buildSlowCmd, "echo", 2, 2).Execute(reader)
	if err == nil {
		t.Fail()
	}
}

func Test_Cancel_Xargs(t *testing.T) {
	ctx := context.Background()

	reader := strings.NewReader("aa bb cc dd cc")

	xargs := New(ctx, buildTestCmd, "echo", 2, 2)
	xargs.Cancel()
	err := xargs.Execute(reader)
	if err == nil {
		t.Fail()
	}
}
