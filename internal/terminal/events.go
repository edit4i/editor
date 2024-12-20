package terminal

// EventType represents the type of terminal event
type EventType int

const (
	EventData EventType = iota
	EventResize
	EventCursor
	EventExit
)

// Event represents a terminal event
type Event struct {
	Type    EventType
	Data    []byte
	Cols    int
	Rows    int
	CursorX int
	CursorY int
}
