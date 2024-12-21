package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// FileStatus represents the status of a file in the Git repository
type FileStatus struct {
	File   string `json:"file"`   // File path relative to repository root
	Status string `json:"status"` // Status code: "M" for modified, "A" for added, "D" for deleted and "?" for untracked
	Staged bool   `json:"staged"` // Whether the file is staged
}

// BranchInfo represents information about a Git branch
type BranchInfo struct {
	Name     string `json:"name"`
	IsRemote bool   `json:"isRemote"`
	IsHead   bool   `json:"isHead"`
}

// CommitInfo represents information about a Git commit
type CommitInfo struct {
	Hash        string    `json:"hash"`
	Message     string    `json:"message"`
	Author      string    `json:"author"`
	AuthorEmail string    `json:"authorEmail"`
	Date        time.Time `json:"date"`
	ParentHashes []string `json:"parentHashes"`
}

// CommitFilter contains options for filtering commits
type CommitFilter struct {
	Branch      string    `json:"branch"`      // Branch to list commits from
	StartHash   string    `json:"startHash"`   // Start listing from this commit
	Limit       int       `json:"limit"`       // Max number of commits to return
	Offset      int       `json:"offset"`      // Skip this many commits (numeric offset)
	OffsetHash  string    `json:"offsetHash"`  // Start listing after this commit hash (more efficient for pagination)
	Author      string    `json:"author"`      // Filter by author
	SearchQuery string    `json:"searchQuery"` // Search in commit messages
	StartDate   time.Time `json:"startDate"`   // Filter commits after this date
	EndDate     time.Time `json:"endDate"`     // Filter commits before this date
}

// GitService handles Git operations for projects
type GitService struct {
	// We might want to add a cache of repositories later
}

// NewGitService creates a new Git service instance
func NewGitService() *GitService {
	return &GitService{}
}

// IsGitRepository checks if the given directory is a Git repository
// Returns true if it is a Git repository, false if not
// Returns error if there was a problem checking (e.g., directory doesn't exist)
func (s *GitService) IsGitRepository(projectPath string) (bool, error) {
	// Ensure we have an absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return false, err
	}

	// Try to open the repository
	_, err = git.PlainOpen(absPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			// Not a Git repository, but not an error
			return false, nil
		}
		// Some other error occurred
		return false, err
	}

	return true, nil
}

// InitRepository initializes a new Git repository in the given directory
func (s *GitService) InitRepository(projectPath string) error {
	// First check if it's already a Git repository
	isRepo, err := s.IsGitRepository(projectPath)
	if err != nil {
		return fmt.Errorf("failed to check if directory is a Git repository: %w", err)
	}
	if isRepo {
		return errors.New("directory is already a Git repository")
	}

	// Ensure we have an absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Initialize the repository
	_, err = git.PlainInit(absPath, false) // false means not bare repository
	if err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	return nil
}

// GetStatus returns the current Git status of the repository
// Returns two slices: staged files and unstaged files
func (s *GitService) GetStatus(projectPath string) ([]FileStatus, error) {
	// Ensure we have an absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Open the repository
	repo, err := git.PlainOpen(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get the status
	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// Convert status to our format
	var files []FileStatus
	for file, fileStatus := range status {
		// Skip unmodified files
		if fileStatus.Staging == git.Unmodified && fileStatus.Worktree == git.Unmodified {
			continue
		}

		fs := FileStatus{
			File: file,
		}

		// For untracked files
		if fileStatus.Worktree == git.Untracked {
			fs.Staged = false
			fs.Status = string(git.Untracked)
		} else if fileStatus.Staging != git.Unmodified {
			// For staged changes
			fs.Staged = true
			fs.Status = string(fileStatus.Staging)
		} else {
			// For unstaged changes
			fs.Staged = false
			fs.Status = string(fileStatus.Worktree)
		}

		files = append(files, fs)
	}

	return files, nil
}

// getWorktree is a helper function that returns the worktree for a given project path
func (s *GitService) getWorktree(projectPath string) (*git.Worktree, error) {
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return worktree, nil
}

// StageFile adds a file to the staging area
func (s *GitService) StageFile(projectPath string, file string) error {
	worktree, err := s.getWorktree(projectPath)
	if err != nil {
		return err
	}

	_, err = worktree.Add(file)
	if err != nil {
		return fmt.Errorf("failed to stage file: %w", err)
	}

	return nil
}

func (s *GitService) UnstageFile(projectPath string, file string) error {
	worktree, err := s.getWorktree(projectPath)
	if err != nil {
		return err
	}

	_, err = worktree.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to unstage file: %w", err)
	}

	return nil
}

// DiscardChanges discards changes in an unstaged file, reverting it to the last commit
func (s *GitService) DiscardChanges(projectPath string, file string) error {
	// Open the repository
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get worktree to check file status
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Check if file is untracked
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fileStatus := status.File(file)
	if fileStatus.Staging == git.Untracked {
		fullPath := filepath.Join(projectPath, file)
		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("failed to delete untracked file: %w", err)
		}
		return nil
	}

	// Get HEAD commit
	ref, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf("failed to get commit: %w", err)
	}

	// Get the tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get tree: %w", err)
	}

	// Find the file entry in the tree to get both content and mode
	entry, err := tree.FindEntry(file)
	if err != nil {
		return fmt.Errorf("failed to find file in tree: %w", err)
	}

	// Get file object
	treeFile, err := tree.File(file)
	if err != nil {
		return fmt.Errorf("failed to get file from tree: %w", err)
	}

	// Get the contents
	contents, err := treeFile.Contents()
	if err != nil {
		return fmt.Errorf("failed to get file contents: %w", err)
	}

	// Write the contents back to the file
	fullPath := filepath.Join(projectPath, file)
	err = os.WriteFile(fullPath, []byte(contents), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Set the original permissions from git tree entry
	mode, modeErr := entry.Mode.ToOSFileMode()
	if modeErr != nil {
		return fmt.Errorf("failed to convert file mode: %w", modeErr)
	}
	if err := os.Chmod(fullPath, mode); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// Commit creates a new commit with the staged changes
func (s *GitService) Commit(projectPath string, message string) error {
	worktree, err := s.getWorktree(projectPath)
	if err != nil {
		return err
	}

	// Create the commit
	_, err = worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// ListBranches returns a list of all branches in the repository
func (s *GitService) ListBranches(projectPath string) ([]BranchInfo, error) {
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	branches := []BranchInfo{}

	// Get current branch to mark the HEAD
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}
	currentBranchName := head.Name().Short()

	// List local branches
	branchIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()
		branches = append(branches, BranchInfo{
			Name:     branchName,
			IsRemote: false,
			IsHead:   branchName == currentBranchName,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	// List remote branches
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to list remotes: %w", err)
	}

	for _, remote := range remotes {
		refs, err := remote.List(&git.ListOptions{})
		if err != nil {
			continue // Skip this remote if we can't list its refs
		}

		for _, ref := range refs {
			if ref.Name().IsBranch() {
				branchName := ref.Name().Short()
				branches = append(branches, BranchInfo{
					Name:     branchName,
					IsRemote: true,
					IsHead:   false,
				})
			}
		}
	}

	return branches, nil
}

// GetCurrentBranch returns the name of the current branch
func (s *GitService) GetCurrentBranch(projectPath string) (string, error) {
	repo, err := git.PlainOpen(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	return head.Name().Short(), nil
}

// ListCommits returns a list of commits based on the provided filters
func (s *GitService) ListCommits(projectPath string, filter CommitFilter) ([]CommitInfo, error) {
	// Ensure we have an absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Open the repository
	repo, err := git.PlainOpen(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the reference to start from (branch or commit)
	var startRef plumbing.Hash
	if filter.Branch != "" {
		// If branch is specified, use it
		ref, err := repo.Reference(plumbing.NewBranchReferenceName(filter.Branch), true)
		if err != nil {
			return nil, fmt.Errorf("failed to get branch reference: %w", err)
		}
		startRef = ref.Hash()
	} else if filter.StartHash != "" {
		// If start hash is specified, use it
		hash := plumbing.NewHash(filter.StartHash)
		if !hash.IsZero() {
			startRef = hash
		}
	} else {
		// Default to HEAD
		ref, err := repo.Head()
		if err != nil {
			return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
		}
		startRef = ref.Hash()
	}

	// Create log options
	logOptions := &git.LogOptions{
		From:  startRef,
		Order: git.LogOrderCommitterTime,
	}

	// Get the commit iterator
	commitIter, err := repo.Log(logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit iterator: %w", err)
	}
	defer commitIter.Close()

	var commits []CommitInfo
	var skipped int
	var foundOffsetHash bool = filter.OffsetHash == "" // If no offset hash specified, we start collecting immediately

	err = commitIter.ForEach(func(c *object.Commit) error {
		// Handle hash-based offset
		if !foundOffsetHash {
			if c.Hash.String() == filter.OffsetHash {
				foundOffsetHash = true
			}
			return nil
		}

		// Apply date filters
		if !filter.StartDate.IsZero() && c.Author.When.Before(filter.StartDate) {
			return nil
		}
		if !filter.EndDate.IsZero() && c.Author.When.After(filter.EndDate) {
			return nil
		}

		// Apply author filter
		if filter.Author != "" && !strings.Contains(c.Author.Name, filter.Author) && !strings.Contains(c.Author.Email, filter.Author) {
			return nil
		}

		// Apply message search
		if filter.SearchQuery != "" && !strings.Contains(strings.ToLower(c.Message), strings.ToLower(filter.SearchQuery)) {
			return nil
		}

		// Handle numeric offset (only if we're not using hash-based offset)
		if filter.OffsetHash == "" && skipped < filter.Offset {
			skipped++
			return nil
		}

		// Create parent hashes list
		parentHashes := make([]string, len(c.ParentHashes))
		for i, hash := range c.ParentHashes {
			parentHashes[i] = hash.String()
		}

		// Add commit to results
		commits = append(commits, CommitInfo{
			Hash:        c.Hash.String(),
			Message:     strings.TrimSpace(c.Message),
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
			Date:        c.Author.When,
			ParentHashes: parentHashes,
		})

		// Check if we've reached the limit
		if filter.Limit > 0 && len(commits) >= filter.Limit {
			return errors.New("stop iteration")
		}

		return nil
	})

	// If we stopped due to reaching the limit, don't treat it as an error
	if err != nil && err.Error() != "stop iteration" {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	// If we were looking for an offset hash but didn't find it
	if filter.OffsetHash != "" && !foundOffsetHash {
		return nil, fmt.Errorf("offset hash %s not found in commit history", filter.OffsetHash)
	}

	return commits, nil
}

// ListCommitsAfter is a convenience function to list commits after a specific commit
func (s *GitService) ListCommitsAfter(projectPath string, offsetHash string, limit int) ([]CommitInfo, error) {
	return s.ListCommits(projectPath, CommitFilter{
		OffsetHash: offsetHash,
		Limit:     limit,
	})
}

// ListCommitsByBranch is a convenience function to list commits from a specific branch
func (s *GitService) ListCommitsByBranch(projectPath string, branch string, limit int) ([]CommitInfo, error) {
	return s.ListCommits(projectPath, CommitFilter{
		Branch: branch,
		Limit:  limit,
	})
}

// ListCommitsByAuthor is a convenience function to list commits by a specific author
func (s *GitService) ListCommitsByAuthor(projectPath string, author string, limit int) ([]CommitInfo, error) {
	return s.ListCommits(projectPath, CommitFilter{
		Author: author,
		Limit:  limit,
	})
}

// SearchCommits is a convenience function to search commits by message
func (s *GitService) SearchCommits(projectPath string, query string, limit int) ([]CommitInfo, error) {
	return s.ListCommits(projectPath, CommitFilter{
		SearchQuery: query,
		Limit:      limit,
	})
}
