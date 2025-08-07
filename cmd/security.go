package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	MaxPathLength     int
	MaxNameLength     int
	AllowedCharacters *regexp.Regexp
	ForbiddenPatterns []*regexp.Regexp
}

// DefaultSecurityConfig returns the default security configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		MaxPathLength:     255,
		MaxNameLength:     50,
		AllowedCharacters: regexp.MustCompile(`^[a-z0-9-]+$`),
		ForbiddenPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\.\.`),                                  // Path traversal
			regexp.MustCompile(`^(con|prn|aux|nul|com[1-9]|lpt[1-9])$`), // Windows reserved names
		},
	}
}

// ValidateAndSanitizePath validates and sanitizes a file path
func ValidateAndSanitizePath(path string, config *SecurityConfig) (string, error) {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	// Check for empty path
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Check path length
	if len(path) > config.MaxPathLength {
		return "", fmt.Errorf("path too long (max %d characters)", config.MaxPathLength)
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("path traversal not allowed")
	}

	// Check for absolute paths
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths not allowed")
	}

	// Check for forbidden patterns
	for _, pattern := range config.ForbiddenPatterns {
		if pattern.MatchString(path) {
			return "", fmt.Errorf("path contains forbidden pattern: %s", pattern.String())
		}
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Ensure the cleaned path doesn't start with ../
	if strings.HasPrefix(cleanPath, "..") {
		return "", fmt.Errorf("path traversal not allowed")
	}

	return cleanPath, nil
}

// ValidateAndSanitizeName validates and sanitizes a project/service name
func ValidateAndSanitizeName(name string, config *SecurityConfig) (string, error) {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	// Check name length
	if len(name) > config.MaxNameLength {
		return "", fmt.Errorf("name too long (max %d characters)", config.MaxNameLength)
	}

	// Check for forbidden patterns
	for _, pattern := range config.ForbiddenPatterns {
		if pattern.MatchString(name) {
			return "", fmt.Errorf("name contains forbidden pattern: %s", pattern.String())
		}
	}

	// Validate character set
	if !config.AllowedCharacters.MatchString(name) {
		return "", fmt.Errorf("name can only contain lowercase letters, numbers, and hyphens")
	}

	// Must start with a letter
	if len(name) > 0 && !unicode.IsLetter(rune(name[0])) {
		return "", fmt.Errorf("name must start with a letter")
	}

	// Must end with a letter or number
	if len(name) > 0 {
		lastChar := rune(name[len(name)-1])
		if !unicode.IsLetter(lastChar) && !unicode.IsNumber(lastChar) {
			return "", fmt.Errorf("name must end with a letter or number")
		}
	}

	return strings.ToLower(strings.TrimSpace(name)), nil
}

// ValidateDirectorySafety performs comprehensive directory safety checks
func ValidateDirectorySafety(dirPath string) error {
	// Check if directory exists and is accessible
	info, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("cannot access directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Check directory permissions
	if info.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("directory is not writable")
	}

	// Check for symbolic links (potential security risk)
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symbolic links not allowed for security reasons")
	}

	return nil
}

// ValidateTemplateName validates a template name for security
func ValidateTemplateName(templateName string) error {
	if strings.TrimSpace(templateName) == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	// Check for path traversal
	if strings.Contains(templateName, "..") || strings.Contains(templateName, "/") || strings.Contains(templateName, "\\") {
		return fmt.Errorf("template name contains invalid characters")
	}

	// Check length
	if len(templateName) > 100 {
		return fmt.Errorf("template name too long")
	}

	// Check for allowed characters only
	allowedPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !allowedPattern.MatchString(templateName) {
		return fmt.Errorf("template name contains invalid characters")
	}

	return nil
}

// SanitizeForFileSystem ensures a string is safe for filesystem operations
func SanitizeForFileSystem(input string) string {
	// Remove or replace potentially dangerous characters
	replacer := strings.NewReplacer(
		"<", "",
		">", "",
		":", "",
		"\"", "",
		"|", "",
		"?", "",
		"*", "",
		"\\", "",
		"/", "",
	)

	result := replacer.Replace(input)

	// Additional cleanup for script tags
	result = strings.ReplaceAll(result, "script", "")
	result = strings.ReplaceAll(result, "alert", "")
	result = strings.ReplaceAll(result, "eval", "")
	result = strings.ReplaceAll(result, "exec", "")
	result = strings.ReplaceAll(result, "system", "")

	return result
}

// CheckForSuspiciousPatterns checks for potentially malicious patterns
func CheckForSuspiciousPatterns(input string) error {
	suspiciousPatterns := []string{
		"javascript:",
		"data:",
		"vbscript:",
		"onload=",
		"onerror=",
		"<script",
		"</script>",
		"eval(",
		"exec(",
		"system(",
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(inputLower, pattern) {
			return fmt.Errorf("input contains suspicious pattern: %s", pattern)
		}
	}

	return nil
}
