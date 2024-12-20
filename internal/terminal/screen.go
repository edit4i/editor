package terminal

import "github.com/gdamore/tcell/v2"

// Screen represents a terminal screen
type Screen struct {
	screen   tcell.Screen
	content  [][]rune
	style    tcell.Style
	cursorX  int
	cursorY  int
	cols     int
	rows     int
	onChange func(*Event)
}

// newScreen creates a new terminal screen
func newScreen(onChange func(*Event)) (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	cols, rows := screen.Size()
	content := make([][]rune, rows)
	for i := range content {
		content[i] = make([]rune, cols)
	}

	return &Screen{
		screen:   screen,
		content:  content,
		style:    tcell.StyleDefault,
		cols:     cols,
		rows:     rows,
		onChange: onChange,
	}, nil
}

// write writes data to the screen
func (s *Screen) write(data []byte) {
	for _, b := range data {
		s.writeChar(rune(b))
	}
	s.update()
}

// writeChar writes a single character to the screen
func (s *Screen) writeChar(ch rune) {
	switch ch {
	case '\n':
		s.cursorY++
		s.cursorX = 0
	case '\r':
		s.cursorX = 0
	case '\b':
		if s.cursorX > 0 {
			s.cursorX--
		}
	default:
		if s.cursorX >= s.cols {
			s.cursorX = 0
			s.cursorY++
		}

		if s.cursorY >= s.rows {
			// Scroll content up
			copy(s.content, s.content[1:])
			s.content[s.rows-1] = make([]rune, s.cols)
			s.cursorY = s.rows - 1
		}

		if s.cursorY < len(s.content) && s.cursorX < len(s.content[s.cursorY]) {
			s.content[s.cursorY][s.cursorX] = ch
			s.cursorX++
		}
	}

	// Notify cursor position change
	s.onChange(&Event{
		Type:    EventCursor,
		CursorX: s.cursorX,
		CursorY: s.cursorY,
	})
}

// resize handles screen resize
func (s *Screen) resize(cols, rows int) {
	s.cols = cols
	s.rows = rows

	// Create new content buffer
	newContent := make([][]rune, rows)
	for i := range newContent {
		newContent[i] = make([]rune, cols)
	}

	// Copy old content
	for y := 0; y < min(len(s.content), rows); y++ {
		for x := 0; x < min(len(s.content[y]), cols); x++ {
			newContent[y][x] = s.content[y][x]
		}
	}

	s.content = newContent
	s.cursorX = min(s.cursorX, cols-1)
	s.cursorY = min(s.cursorY, rows-1)

	s.update()
}

// update updates the screen display
func (s *Screen) update() {
	s.screen.Clear()

	// Write content to screen
	for y := 0; y < s.rows && y < len(s.content); y++ {
		for x := 0; x < s.cols && x < len(s.content[y]); x++ {
			ch := s.content[y][x]
			if ch == 0 {
				ch = ' '
			}
			s.screen.SetContent(x, y, ch, nil, s.style)
		}
	}

	s.screen.ShowCursor(s.cursorX, s.cursorY)
	s.screen.Show()
}

// close closes the screen
func (s *Screen) close() {
	s.screen.Fini()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
