package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitConfig represents Git configuration
type GitConfig struct {
	Name  string
	Email string
}

// IsGitRepository checks if a directory is a Git repository
func IsGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	return DirExists(gitDir)
}

// InitGitRepository initializes a new Git repository in the specified directory
func InitGitRepository(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

// AddGitIgnore creates a .gitignore file with the specified content
func AddGitIgnore(projectPath string, content string) error {
	gitignorePath := filepath.Join(projectPath, ".gitignore")

	if err := WriteFile(gitignorePath, content); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// GetGitConfig retrieves the global Git configuration
func GetGitConfig() (*GitConfig, error) {
	name, err := getGitConfigValue("user.name")
	if err != nil {
		return nil, fmt.Errorf("failed to get git user.name: %w", err)
	}

	email, err := getGitConfigValue("user.email")
	if err != nil {
		return nil, fmt.Errorf("failed to get git user.email: %w", err)
	}

	return &GitConfig{
		Name:  name,
		Email: email,
	}, nil
}

// getGitConfigValue gets a specific Git configuration value
func getGitConfigValue(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", key)
	output, err := cmd.Output()
	if err != nil {
		// If the config value doesn't exist, return empty string
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// SetGitConfig sets Git configuration values
func SetGitConfig(name, email string) error {
	if name != "" {
		cmd := exec.Command("git", "config", "--global", "user.name", name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.name: %w", err)
		}
	}

	if email != "" {
		cmd := exec.Command("git", "config", "--global", "user.email", email)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set git user.email: %w", err)
		}
	}

	return nil
}

// GitAdd adds files to the Git staging area
func GitAdd(projectPath string, files ...string) error {
	if len(files) == 0 {
		files = []string{"."}
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	return nil
}

// GitCommit creates a Git commit with the specified message
func GitCommit(projectPath, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create git commit: %w", err)
	}

	return nil
}

// GitStatus returns the status of the Git repository
func GitStatus(projectPath string) (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}

	return string(output), nil
}

// GitBranch returns the current branch name
func GitBranch(projectPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GitRemoteAdd adds a remote repository
func GitRemoteAdd(projectPath, name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add git remote: %w", err)
	}

	return nil
}

// GitRemoteList lists all remote repositories
func GitRemoteList(projectPath string) ([]string, error) {
	cmd := exec.Command("git", "remote")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list git remotes: %w", err)
	}

	remotes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(remotes) == 1 && remotes[0] == "" {
		return []string{}, nil
	}

	return remotes, nil
}

// GitPush pushes changes to the remote repository
func GitPush(projectPath, remote, branch string) error {
	cmd := exec.Command("git", "push", remote, branch)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to git remote: %w", err)
	}

	return nil
}

// GitPull pulls changes from the remote repository
func GitPull(projectPath, remote, branch string) error {
	cmd := exec.Command("git", "pull", remote, branch)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull from git remote: %w", err)
	}

	return nil
}

// GitClone clones a repository to the specified path
func GitClone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone git repository: %w", err)
	}

	return nil
}

// GitLog returns the commit history
func GitLog(projectPath string, maxCount int) (string, error) {
	args := []string{"log", "--oneline"}
	if maxCount > 0 {
		args = append(args, fmt.Sprintf("-%d", maxCount))
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git log: %w", err)
	}

	return string(output), nil
}

// GitDiff returns the diff of changes
func GitDiff(projectPath string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	return string(output), nil
}

// GitTag creates a Git tag
func GitTag(projectPath, tag, message string) error {
	args := []string{"tag"}
	if message != "" {
		args = append(args, "-a", tag, "-m", message)
	} else {
		args = append(args, tag)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create git tag: %w", err)
	}

	return nil
}

// GitListTags lists all Git tags
func GitListTags(projectPath string) ([]string, error) {
	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list git tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 1 && tags[0] == "" {
		return []string{}, nil
	}

	return tags, nil
}

// IsGitInstalled checks if Git is installed and available
func IsGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	return cmd.Run() == nil
}

// GetGitVersion returns the installed Git version
func GetGitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git version: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateGitRepository creates a new Git repository with initial commit
func CreateGitRepository(projectPath string, gitignoreContent string) error {
	// Initialize repository
	if err := InitGitRepository(projectPath); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Add .gitignore if content provided
	if gitignoreContent != "" {
		if err := AddGitIgnore(projectPath, gitignoreContent); err != nil {
			return fmt.Errorf("failed to add .gitignore: %w", err)
		}
	}

	// Add all files
	if err := GitAdd(projectPath); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Create initial commit
	if err := GitCommit(projectPath, "Initial commit"); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}

// GetDefaultGitIgnore returns a default .gitignore content for Go projects
func GetDefaultGitIgnore() string {
	return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log

# Environment variables
.env
.env.local
.env.*.local

# Build artifacts
/bin/
/dist/
/build/

# Temporary files
/tmp/
*.tmp
*.temp

# Configuration files with secrets
config.local.yaml
config.secret.yaml
*.secret.yaml
`
}

// HasUncommittedChanges checks if there are uncommitted changes
func HasUncommittedChanges(projectPath string) (bool, error) {
	status, err := GitStatus(projectPath)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(status) != "", nil
}

// GetLastCommitHash returns the hash of the last commit
func GetLastCommitHash(projectPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit hash: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CheckGitInstallation checks if Git is properly installed and configured
func CheckGitInstallation() error {
	if !IsGitInstalled() {
		return fmt.Errorf("Git is not installed or not available in PATH")
	}

	config, err := GetGitConfig()
	if err != nil {
		return fmt.Errorf("failed to get Git configuration: %w", err)
	}

	var warnings []string
	if config.Name == "" {
		warnings = append(warnings, "Git user.name is not configured")
	}
	if config.Email == "" {
		warnings = append(warnings, "Git user.email is not configured")
	}

	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", strings.Join(warnings, ", "))
		fmt.Fprintf(os.Stderr, "You can configure Git with:\n")
		fmt.Fprintf(os.Stderr, "  git config --global user.name \"Your Name\"\n")
		fmt.Fprintf(os.Stderr, "  git config --global user.email \"your.email@example.com\"\n")
	}

	return nil
}
