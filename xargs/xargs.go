package xargs

import (
	"bufio"
	"context"
	"io"
	"sync"
)

type Cmder interface {
	Exec() error
}

type BuildCmd func(command string, args ...string) Cmder

type Xargs struct {
	ctx      context.Context
	cancel   context.CancelFunc
	maxprocs int
	number   int
	command  string
	errch    chan error
	buildCmd BuildCmd
	worker   chan struct{}
	wg       *sync.WaitGroup
}

func New(ctx context.Context, buildCmd BuildCmd, command string, maxprocs, number int) *Xargs {
	ctx, cancel := context.WithCancel(ctx)

	worker := make(chan struct{}, maxprocs)

	for i := 0; i < maxprocs; i++ {
		worker <- struct{}{}
	}

	return &Xargs{
		command:  command,
		number:   number,
		ctx:      ctx,
		cancel:   cancel,
		errch:    make(chan error, 1),
		buildCmd: buildCmd,
		worker:   worker,
		wg:       &sync.WaitGroup{},
	}
}

func (a *Xargs) Execute(reader io.Reader) error {
	defer a.cancel()

	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanWords)

	var args []string
	for s.Scan() {
		args = append(args, s.Text())
		if len(args) >= a.number {
			a.work(args)
			args = []string{}
		}
	}
	if len(args) > 0 {
		a.work(args)
	}
	a.wg.Wait()

	select {
	case err := <-a.errch:
		return err
	default:
		return nil
	}
}

func (a *Xargs) work(g []string) {
	a.wg.Add(1)
	go func() {
		select {
		case <-a.worker:
			defer func() {
				a.worker <- struct{}{}
				a.wg.Done()
			}()
			if err := a.buildCmd(a.command, g...).Exec(); err != nil {
				a.errch <- err
				return
			}
		case <-a.ctx.Done():
			a.errch <- a.ctx.Err()
			return
		}
	}()
}

func (a *Xargs) Cancel() { a.cancel() }
