package main

import (
	"context"
	"fmt"
	"os"

	"github.com/edit4i/editor/internal/db"
	"github.com/edit4i/editor/internal/service"
	"github.com/edit4i/editor/internal/terminal"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx             context.Context
	projects        *service.ProjectsService
	files           *service.FileService
	config          *service.ConfigService
	terminalService *service.TerminalService
	git             *service.GitService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	dbConn, err := db.InitDB(db.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Initialize services
	a.projects = service.NewProjectsService(dbConn)
	a.files = service.NewFileService()
	a.git = service.NewGitService()

	config, err := service.NewConfigService()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize ConfigService: %v", err))
	}
	a.config = config

	// Initialize terminal service with event handler
	a.terminalService = service.NewTerminalService(func(id string, event *terminal.Event) {
		// Emit terminal events to frontend
		runtime.EventsEmit(a.ctx, fmt.Sprintf("terminal:%s", id), event)
	})
}

// GetRecentProjects returns the list of recent projects
func (a *App) GetRecentProjects() ([]db.Project, error) {
	return a.projects.GetRecentProjects(4)
}

// OpenProjectFolder opens a folder selection dialog
func (a *App) OpenProjectFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/"
	}

	options := runtime.OpenDialogOptions{
		Title:            "Open Project Folder",
		DefaultDirectory: homeDir,
	}

	path, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", fmt.Errorf("error opening directory dialog: %v", err)
	}

	return path, nil
}

// AddProject adds a new project or updates existing one
func (a *App) AddProject(name, path string) (*db.Project, error) {
	return a.projects.AddProject(name, path)
}

// GetProjectFiles returns the file tree for a project
func (a *App) GetProjectFiles(projectPath string) (*service.FileNode, error) {
	return a.files.GetProjectFiles(projectPath)
}

// LoadDirectoryContents loads the contents of a specific directory
func (a *App) LoadDirectoryContents(dirPath string) (*service.FileNode, error) {
	return a.files.LoadDirectoryContents(dirPath)
}

// GetFileContent returns the content of a file
func (a *App) GetFileContent(path string) (string, error) {
	return a.files.GetFileContent(path)
}

// SaveFile saves content to a file
func (a *App) SaveFile(path, content string) error {
	return a.files.SaveFile(path, content)
}

// SearchFiles performs a fuzzy search on files in a directory
func (a *App) SearchFiles(dirPath, query string) ([]*service.FileNode, error) {
	// Create a new context that will be cancelled when a new search starts
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()

	return a.files.SearchFiles(ctx, dirPath, query)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetEditorConfig returns the current editor configuration
func (a *App) GetEditorConfig() (*service.EditorConfig, error) {
	return a.config.GetConfig(), nil
}

// OpenConfigFile opens the config file in the editor
func (a *App) OpenConfigFile() (string, error) {
	path := a.config.OpenConfigFile()
	return path, nil
}

// CreateFile creates a new empty file
func (a *App) CreateFile(path string) error {
	return a.files.CreateFile(path)
}

// CreateDirectory creates a new directory
func (a *App) CreateDirectory(path string) error {
	return a.files.CreateDirectory(path)
}

// RenameFile renames a file or directory
func (a *App) RenameFile(oldPath, newPath string) error {
	return a.files.RenameFile(oldPath, newPath)
}

// DeleteFile deletes a file or directory
func (a *App) DeleteFile(path string) error {
	return a.files.DeleteFile(path)
}

// CreateTerminal creates a new terminal instance
func (a *App) CreateTerminal(id string, shell string, cwd string) error {
	return a.terminalService.CreateTerminal(id, shell, cwd)
}

// DestroyTerminal destroys a terminal instance
func (a *App) DestroyTerminal(id string) error {
	return a.terminalService.DestroyTerminal(id)
}

// ResizeTerminal resizes a terminal instance
func (a *App) ResizeTerminal(id string, cols int, rows int) error {
	return a.terminalService.ResizeTerminal(id, cols, rows)
}

// HandleInput handles terminal input from the frontend
func (a *App) HandleInput(id string, data []byte) error {
	return a.terminalService.HandleInput(id, data)
}

// GetAvailableShells returns a list of available shells, with the default shell as the first item
func (a *App) GetAvailableShells() ([]string, error) {
	return a.terminalService.GetAvailableShells()
}

// IsGitRepository checks if the given directory is a Git repository
func (a *App) IsGitRepository(projectPath string) (bool, error) {
	return a.git.IsGitRepository(projectPath)
}

// InitGitRepository initializes a new Git repository in the given directory
func (a *App) InitGitRepository(projectPath string) error {
	return a.git.InitRepository(projectPath)
}

// GetGitStatus returns the current Git status of the repository
func (a *App) GetGitStatus(projectPath string) ([]service.FileStatus, error) {
	return a.git.GetStatus(projectPath)
}

// StageFile adds a file to the staging area
func (a *App) StageFile(projectPath string, file string) error {
	return a.git.StageFile(projectPath, file)
}

// UnstageFile removes a file from the staging area
func (a *App) UnstageFile(projectPath string, file string) error {
	return a.git.UnstageFile(projectPath, file)
}

// DiscardChanges discards changes in a file, reverting it to the last commit
func (a *App) DiscardChanges(projectPath string, file string) error {
	return a.git.DiscardChanges(projectPath, file)
}

// Commit creates a new commit with the staged changes
func (a *App) Commit(projectPath string, message string) error {
	return a.git.Commit(projectPath, message)
}

// ListBranches returns a list of all branches in the repository
func (a *App) ListBranches(projectPath string) ([]service.BranchInfo, error) {
	return a.git.ListBranches(projectPath)
}

// GetCurrentBranch returns the name of the current branch
func (a *App) GetCurrentBranch(projectPath string) (string, error) {
	return a.git.GetCurrentBranch(projectPath)
}

// ListCommits returns a list of commits based on the provided filters
func (a *App) ListCommits(projectPath string, filter service.CommitFilter) ([]service.CommitInfo, error) {
	return a.git.ListCommits(projectPath, filter)
}

// ListCommitsAfter returns commits after a specific commit hash
func (a *App) ListCommitsAfter(projectPath string, offsetHash string, limit int) ([]service.CommitInfo, error) {
	return a.git.ListCommitsAfter(projectPath, offsetHash, limit)
}

// ListCommitsByBranch returns commits from a specific branch
func (a *App) ListCommitsByBranch(projectPath string, branch string, limit int) ([]service.CommitInfo, error) {
	return a.git.ListCommitsByBranch(projectPath, branch, limit)
}

// ListCommitsByAuthor returns commits by a specific author
func (a *App) ListCommitsByAuthor(projectPath string, author string, limit int) ([]service.CommitInfo, error) {
	return a.git.ListCommitsByAuthor(projectPath, author, limit)
}

// SearchCommits searches for commits by message
func (a *App) SearchCommits(projectPath string, query string, limit int) ([]service.CommitInfo, error) {
	return a.git.SearchCommits(projectPath, query, limit)
}

// GetHeadCommit returns the head commit of the repository
func (a *App) GetHeadCommit(projectPath string) (*service.CommitInfo, error) {
	return a.git.GetHeadCommit(projectPath)
}
