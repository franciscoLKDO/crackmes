package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
)

type HandlerFunc func(reader io.Reader, writer io.WriteCloser) error

type Context struct {
	Session *Session
	Input   []byte
}

type Session struct {
	Stdin  io.WriteCloser
	Stdout io.Reader
}

type route struct {
	match   []byte
	handler HandlerFunc
}

type Engine struct {
	routes []route
	ptmx   *os.File
}

func New(c *exec.Cmd) (Engine, error) {
	ptmx, err := pty.Start(c)
	if err != nil {
		return Engine{}, err
	}

	return Engine{
		ptmx: ptmx,
	}, nil
}

func (e *Engine) Handle(match string, h HandlerFunc) {
	e.routes = append(e.routes, route{
		match:   []byte(match),
		handler: h,
	})
}

func (e Engine) Dispatch(ctx context.Context, data []byte) error {
	for _, r := range e.routes {
		if bytes.Contains(data, r.match) {
			return r.handler(e.ptmx, e.ptmx)
		}
	}
	return nil
}

func (e Engine) Serve() error {
	defer func() { _ = e.ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, e.ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }()
	go func() { e.ptmx.ReadFrom(os.Stdin) }()

	buf := make([]byte, 4096)
	var line []byte
	for {
		n, err := e.ptmx.Read(buf)
		if err != nil {
			return err
		}
		line = append(line, buf[:n]...)
		fmt.Fprintf(os.Stdout, "%s", buf[:n])

		if err := e.Dispatch(context.Background(), line); err != nil {
			return err
		}

		if bytes.Contains(line, []byte("\n")) {
			line = []byte{}
		}
	}
}
