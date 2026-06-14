package engine

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Session struct {
	ptmx *os.File
	cmd  *exec.Cmd
	in   chan []byte
}

func NewSession(c *exec.Cmd, out io.Writer) (Session, error) {
	ptmx, err := pty.Start(c)
	if err != nil {
		return Session{}, err
	}

	s := Session{
		ptmx: ptmx,
		cmd:  c,
		in:   make(chan []byte, 100),
	}
	return s, nil
}

func (s Session) Read(p []byte) (int, error) {
	return s.ptmx.Read(p)
}

func (s Session) Write(p []byte) (int, error) {
	return s.ptmx.Write(p)
}

func (s *Session) Wait() error {
	return s.cmd.Wait()
}

func (s *Session) Close() error {
	close(s.in)
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}

	return s.ptmx.Close()
}

func (s *Session) Resize() error {
	return pty.InheritSize(os.Stdin, s.ptmx)
}
