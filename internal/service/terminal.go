package service

import (
	"fmt"
	"sync"

	"github.com/edit4i/editor/internal/terminal"
)

// TerminalService manages multiple terminal instances
type TerminalService struct {
	terminals map[string]*terminal.Terminal
	mu        sync.RWMutex
	onEvent   func(id string, event *terminal.Event)
}

// NewTerminalService creates a new terminal service
func NewTerminalService(onEvent func(id string, event *terminal.Event)) *TerminalService {
	return &TerminalService{
		terminals: make(map[string]*terminal.Terminal),
		onEvent:   onEvent,
	}
}

// CreateTerminal creates a new terminal instance with the specified shell
func (s *TerminalService) CreateTerminal(id string, shell string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if terminal already exists
	if _, exists := s.terminals[id]; exists {
		return fmt.Errorf("terminal with id %s already exists", id)
	}

	// Create event handler for this terminal
	terminalEventHandler := func(event *terminal.Event) {
		if s.onEvent != nil {
			s.onEvent(id, event)
		}
	}

	// Create new terminal
	term, err := terminal.NewTerminal(terminal.TerminalOptions{
		Shell: shell,
		Cols:  80,
		Rows:  24,
	}, terminalEventHandler)
	if err != nil {
		return fmt.Errorf("failed to create terminal: %w", err)
	}

	// Start the terminal
	if err := term.Start(); err != nil {
		return fmt.Errorf("failed to start terminal: %w", err)
	}

	s.terminals[id] = term
	return nil
}

// DestroyTerminal stops and removes a terminal instance
func (s *TerminalService) DestroyTerminal(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	term, exists := s.terminals[id]
	if !exists {
		return fmt.Errorf("terminal with id %s not found", id)
	}

	term.Stop()
	delete(s.terminals, id)

	// Notify that terminal has exited
	if s.onEvent != nil {
		s.onEvent(id, &terminal.Event{
			Type: terminal.EventExit,
		})
	}

	return nil
}

// GetTerminal returns a terminal instance by ID
func (s *TerminalService) GetTerminal(id string) (*terminal.Terminal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	term, exists := s.terminals[id]
	if !exists {
		return nil, fmt.Errorf("terminal with id %s not found", id)
	}

	return term, nil
}

// ResizeTerminal resizes the terminal window
func (s *TerminalService) ResizeTerminal(id string, cols, rows int) error {
	term, err := s.GetTerminal(id)
	if err != nil {
		return err
	}

	return term.Resize(cols, rows)
}

// WriteToTerminal writes data to the terminal
func (s *TerminalService) WriteToTerminal(id string, data []byte) error {
	term, err := s.GetTerminal(id)
	if err != nil {
		return err
	}

	return term.Write(data)
}

// HandleInput handles input from the frontend
func (s *TerminalService) HandleInput(id string, data []byte) error {
	term, err := s.GetTerminal(id)
	if err != nil {
		return err
	}

	return term.HandleInput(data)
}
