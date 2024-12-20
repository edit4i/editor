package terminal

import (
	"fmt"
	"sync"
)

type Terminal struct {
	screen  *Screen
	done    chan struct{}
	events  chan *Event
	mu      sync.Mutex
	onEvent func(*Event)
}

func NewTerminal(opts TerminalOptions, onEvent func(*Event)) (*Terminal, error) {
	screen, err := newScreen(func(ev *Event) {
		if onEvent != nil {
			onEvent(ev)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	t := &Terminal{
		screen:  screen,
		done:    make(chan struct{}),
		events:  make(chan *Event, 100),
		onEvent: onEvent,
	}

	// Set initial size
	if opts.Cols > 0 && opts.Rows > 0 {
		screen.resize(opts.Cols, opts.Rows)
	}

	return t, nil
}

func (t *Terminal) Start() error {
	go t.handleEvents()
	return nil
}

func (t *Terminal) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	close(t.done)
	t.screen.close()
}

func (t *Terminal) Write(data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.screen.write(data)
	if t.onEvent != nil {
		t.onEvent(&Event{
			Type: EventData,
			Data: data,
		})
	}
	return nil
}

func (t *Terminal) HandleInput(data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.screen.write(data)
	return nil
}

func (t *Terminal) Resize(cols, rows int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.screen.resize(cols, rows)
	if t.onEvent != nil {
		t.onEvent(&Event{
			Type: EventResize,
			Cols: cols,
			Rows: rows,
		})
	}
	return nil
}

func (t *Terminal) handleEvents() {
	for {
		select {
		case <-t.done:
			return
		case ev := <-t.events:
			t.mu.Lock()
			switch ev.Type {
			case EventData:
				t.screen.write(ev.Data)
			case EventResize:
				t.screen.resize(ev.Cols, ev.Rows)
			}
			t.mu.Unlock()

			// Notify event handler
			if t.onEvent != nil {
				t.onEvent(ev)
			}
		}
	}
}
