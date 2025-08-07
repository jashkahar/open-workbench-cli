package templating

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// PlatformUtils provides cross-platform utility functions
type PlatformUtils struct{}

// NewPlatformUtils creates a new platform utilities instance
func NewPlatformUtils() *PlatformUtils {
	return &PlatformUtils{}
}

// CopyFile copies a file from source to destination in a cross-platform way
func (pu *PlatformUtils) CopyFile(source, dest string) error {
	// Read source file
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", source, err)
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(dest, content, 0644); err != nil {
		return fmt.Errorf("failed to write destination file %s: %w", dest, err)
	}

	return nil
}

// FileExists checks if a file exists in a cross-platform way
func (pu *PlatformUtils) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetCrossPlatformCopyCommand returns a command that works on all platforms
func (pu *PlatformUtils) GetCrossPlatformCopyCommand(source, dest string) string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("copy \"%s\" \"%s\" 2>nul", source, dest)
	default:
		return fmt.Sprintf("cp \"%s\" \"%s\" 2>/dev/null", source, dest)
	}
}

// GetCrossPlatformEchoCommand returns a cross-platform echo command
func (pu *PlatformUtils) GetCrossPlatformEchoCommand(message string) string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("echo %s", message)
	default:
		return fmt.Sprintf("echo '%s'", message)
	}
}

// GetCrossPlatformCommandWithFallback returns a command that tries multiple approaches
func (pu *PlatformUtils) GetCrossPlatformCommandWithFallback(baseCommand string) string {
	switch runtime.GOOS {
	case "windows":
		// For Windows, try Windows commands first, then fallback to Unix-like commands
		if strings.Contains(baseCommand, "cp ") {
			// Replace cp with copy for Windows
			parts := strings.SplitN(baseCommand, " ", 3)
			if len(parts) >= 3 {
				return fmt.Sprintf("copy \"%s\" \"%s\" 2>nul || %s", parts[1], parts[2], baseCommand)
			}
		}
		return baseCommand
	default:
		// For Unix-like systems, the command should work as-is
		return baseCommand
	}
}

// IsWindows returns true if running on Windows
func (pu *PlatformUtils) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsUnixLike returns true if running on Unix-like systems (Linux, macOS)
func (pu *PlatformUtils) IsUnixLike() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}

// GetShellCommand returns the appropriate shell and arguments for the current platform
func (pu *PlatformUtils) GetShellCommand(command string) (string, []string) {
	switch runtime.GOOS {
	case "windows":
		return "cmd", []string{"/C", command}
	case "darwin":
		return "bash", []string{"-c", command}
	default:
		return "sh", []string{"-c", command}
	}
}
