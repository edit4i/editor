package terminal

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// TerminalManager manages all terminal instances
type TerminalManager struct {
	terminals sync.Map // map[string]*Terminal
}

var (
	manager = &TerminalManager{}
)

// GetTerminal gets a terminal by ID
func GetTerminal(id string) *Terminal {
	if t, ok := manager.terminals.Load(id); ok {
		return t.(*Terminal)
	}
	return nil
}

// NewTerminal creates a new terminal instance
func NewTerminal(id string, opts TerminalOptions, onEvent func(*Event)) (*Terminal, error) {
	t := &Terminal{
		done:    make(chan struct{}),
		onEvent: onEvent,
		shell:   opts.Shell,
	}

	// Store the terminal
	manager.terminals.Store(id, t)

	return t, nil
}

// Terminal represents a terminal instance
type Terminal struct {
	done    chan struct{}
	mu      sync.Mutex
	onEvent func(*Event)
	shell   string
	cmd     *exec.Cmd
	pty     *os.File
}

// Start starts the terminal
func (t *Terminal) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Start the shell process with PTY
	t.cmd = exec.Command(t.shell)
	t.cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Start with PTY
	var err error
	t.pty, err = pty.Start(t.cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}

	// Start reading from pty in a goroutine
	go func() {
		buffer := make([]byte, 4096)
		for {
			select {
			case <-t.done:
				return
			default:
				n, err := t.pty.Read(buffer)
				if err != nil {
					if err != io.EOF {
						log.Printf("[Terminal] Error reading from pty: %v", err)
					}
					return
				}
				if n > 0 {
					// Send data to frontend
					if t.onEvent != nil {
						t.onEvent(&Event{
							Type: EventData,
							Data: buffer[:n],
						})
					}
				}
			}
		}
	}()

	return nil
}

// Stop stops the terminal and removes it from the manager
func (t *Terminal) Stop(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}

	if t.pty != nil {
		t.pty.Close()
	}

	close(t.done)

	// Remove from manager
	manager.terminals.Delete(id)
}

// Write writes data directly to the terminal
func (t *Terminal) Write(data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.pty != nil {
		if _, err := t.pty.Write(data); err != nil {
			return fmt.Errorf("failed to write to pty: %w", err)
		}
	}

	return nil
}

// HandleInput handles input from the frontend
func (t *Terminal) HandleInput(data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Write input to pty
	if t.pty != nil {
		if _, err := t.pty.Write(data); err != nil {
			return fmt.Errorf("failed to write to pty: %w", err)
		}
	}

	return nil
}

// Resize resizes the terminal
func (t *Terminal) Resize(cols, rows int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Resize the pty
	if t.pty != nil {
		if err := pty.Setsize(t.pty, &pty.Winsize{
			Rows: uint16(rows),
			Cols: uint16(cols),
		}); err != nil {
			return fmt.Errorf("failed to resize pty: %w", err)
		}

		// Notify about resize
		if t.onEvent != nil {
			t.onEvent(&Event{
				Type: EventResize,
				Cols: cols,
				Rows: rows,
			})
		}
	}

	return nil
}
