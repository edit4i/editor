package terminal

// TerminalOptions contains options for creating a new terminal
type TerminalOptions struct {
    Shell string
    Cols  int
    Rows  int
    Cwd   string // Working directory for the terminal
}
