package service

import (
	"fmt"
	"log"
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
	log.Println("[TerminalService] Creating new terminal service")
	return &TerminalService{
		terminals: make(map[string]*terminal.Terminal),
		onEvent:   onEvent,
	}
}

// CreateTerminal creates a new terminal instance with the specified shell
func (s *TerminalService) CreateTerminal(id string, shell string) error {
	log.Printf("[TerminalService] Creating terminal: id=%s, shell=%s", id, shell)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if terminal already exists
	if _, exists := s.terminals[id]; exists {
		log.Printf("[TerminalService] Terminal %s already exists", id)
		return fmt.Errorf("terminal with id %s already exists", id)
	}

	// Create event handler for this terminal
	terminalEventHandler := func(event *terminal.Event) {
		log.Printf("[TerminalService] Event from terminal %s: type=%d", id, event.Type)
		if s.onEvent != nil {
			s.onEvent(id, event)
		}
	}

	// Create new terminal
	term, err := terminal.NewTerminal(id, terminal.TerminalOptions{
		Shell: shell,
		Cols:  80,
		Rows:  24,
	}, terminalEventHandler)
	if err != nil {
		log.Printf("[TerminalService] Failed to create terminal: %v", err)
		return fmt.Errorf("failed to create terminal: %w", err)
	}

	// Send welcome message
	log.Printf("[TerminalService] Sending welcome message to terminal %s", id)
	term.Write([]byte(fmt.Sprintf("Welcome to Edit4i Terminal\r\nShell: %s\r\n", shell)))

	// Start the terminal
	log.Printf("[TerminalService] Starting terminal %s", id)
	if err := term.Start(); err != nil {
		log.Printf("[TerminalService] Failed to start terminal: %v", err)
		return fmt.Errorf("failed to start terminal: %w", err)
	}

	s.terminals[id] = term
	log.Printf("[TerminalService] Terminal %s created successfully", id)
	return nil
}

// DestroyTerminal stops and removes a terminal instance
func (s *TerminalService) DestroyTerminal(id string) error {
	log.Printf("[TerminalService] Destroying terminal %s", id)
	s.mu.Lock()
	defer s.mu.Unlock()

	term, exists := s.terminals[id]
	if !exists {
		log.Printf("[TerminalService] Terminal %s not found", id)
		return fmt.Errorf("terminal with id %s not found", id)
	}

	term.Stop(id)
	delete(s.terminals, id)

	// Notify that terminal has exited
	if s.onEvent != nil {
		s.onEvent(id, &terminal.Event{
			Type: terminal.EventExit,
		})
	}

	log.Printf("[TerminalService] Terminal %s destroyed", id)
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
	log.Printf("[TerminalService] Resizing terminal %s to %dx%d", id, cols, rows)
	term, err := s.GetTerminal(id)
	if err != nil {
		log.Printf("[TerminalService] Error resizing terminal: %v", err)
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
	log.Printf("[TerminalService] Input received for terminal %s: %v", id, data)
	term, err := s.GetTerminal(id)
	if err != nil {
		log.Printf("[TerminalService] Error handling input: %v", err)
		return err
	}

	return term.HandleInput(data)
}
