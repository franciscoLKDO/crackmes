package engine

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
)

type HandlerFunc func(reader io.Reader, writer io.WriteCloser) error

type route struct {
	match   []byte
	handler HandlerFunc
}

type Engine struct {
	mu      sync.Mutex
	session *Session
	routes  []route
	stdin   *bufio.Reader
}

func NewEngine(s *Session) Engine {
	return Engine{
		session: s,
		stdin:   bufio.NewReader(os.Stdin),
	}
}

func (e *Engine) Register(match string, h HandlerFunc) {
	e.routes = append(e.routes, route{
		match:   []byte(match),
		handler: h,
	})
}

func (e *Engine) Dispatch(data []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, r := range e.routes {
		if bytes.Contains(data, r.match) {
			return r.handler(e.session.ptmx, e.session.ptmx)
		}
	}
	return nil
}

func (e *Engine) Run() {
	buf := make([]byte, 4096)
	for {
		n, _ := e.session.Read(buf)
		if n == 0 {
			continue
		}
		data := buf[:n]

		e.Dispatch(data)
		os.Stdout.Write(data)
	}
}

// runInput reads from user, processes through handler, sends to binary
func (e *Engine) runInput() {
	for {
		// Read from stdin
		line, _ := ReadLineFromReader(e.stdin)
		if len(line) == 0 {
			continue
		}

		// Send (transformed or original) to binary
		e.session.Write(line)
	}
}

// ReadLineFromReader reads a line from a reader
// Used for interactive input in handlers
func ReadLineFromReader(r io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		return append(scanner.Bytes(), '\n'), nil
	}
	return nil, scanner.Err()
}

func (e *Engine) Start() {
	go e.Run()
	go e.runInput()
}
