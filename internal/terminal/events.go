package terminal

import "github.com/gdamore/tcell/v2"

// handleTcellEvent converts tcell events to terminal events
func handleTcellEvent(ev tcell.Event) *Event {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return handleKeyEvent(ev)
	case *tcell.EventResize:
		width, height := ev.Size()
		return &Event{
			Type: EventResize,
			Cols: width,
			Rows: height,
		}
	}
	return nil
}

// handleKeyEvent converts tcell key events to terminal events
func handleKeyEvent(ev *tcell.EventKey) *Event {
	switch ev.Key() {
	case tcell.KeyEnter:
		return &Event{
			Type: EventData,
			Data: []byte{'\r', '\n'},
		}
	case tcell.KeyTab:
		return &Event{
			Type: EventData,
			Data: []byte{'\t'},
		}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return &Event{
			Type: EventData,
			Data: []byte{'\b'},
		}
	case tcell.KeyCtrlC:
		return &Event{
			Type: EventData,
			Data: []byte("^C\r\n"),
		}
	case tcell.KeyRune:
		return &Event{
			Type: EventData,
			Data: []byte(string(ev.Rune())),
		}
	}
	return nil
}
