// Package templating provides the core templating system for the Open Workbench CLI.
// This package implements dynamic template discovery, parameter processing, and
// file generation capabilities with support for conditional logic and validation.
package templating

import (
	"fmt"
	"os"
)

// ErrorType represents different categories of errors that can occur
// during template processing, allowing for appropriate error handling
// and user-friendly error messages.
type ErrorType int

const (
	// ErrorTypeTemplateNotFound indicates that a requested template could not be found
	ErrorTypeTemplateNotFound ErrorType = iota
	// ErrorTypeInvalidManifest indicates that a template.json file is malformed or invalid
	ErrorTypeInvalidManifest
	// ErrorTypeParameterValidation indicates that a parameter value failed validation
	ErrorTypeParameterValidation
	// ErrorTypeFileSystem indicates a file system operation failed
	ErrorTypeFileSystem
	// ErrorTypeCommandExecution indicates that a post-scaffolding command failed
	ErrorTypeCommandExecution
	// ErrorTypeTemplateProcessing indicates that template processing failed
	ErrorTypeTemplateProcessing
	// ErrorTypePermission indicates a permission-related error
	ErrorTypePermission
	// ErrorTypeNetwork indicates a network-related error (e.g., git operations)
	ErrorTypeNetwork
)

// TemplateError represents a structured error with context and user-friendly messages
type TemplateError struct {
	Type        ErrorType // The category of error
	Message     string    // User-friendly error message
	Details     string    // Technical details for debugging
	Template    string    // Template name where error occurred
	Parameter   string    // Parameter name if applicable
	FilePath    string    // File path if applicable
	Command     string    // Command if applicable
	OriginalErr error     // Original error for debugging
}

// Error returns the user-friendly error message
func (e *TemplateError) Error() string {
	return e.Message
}

// Unwrap returns the original error for debugging
func (e *TemplateError) Unwrap() error {
	return e.OriginalErr
}

// NewTemplateError creates a new structured template error with appropriate
// user-friendly messages based on the error type and context.
func NewTemplateError(errType ErrorType, message, details string, originalErr error) *TemplateError {
	return &TemplateError{
		Type:        errType,
		Message:     message,
		Details:     details,
		OriginalErr: originalErr,
	}
}

// NewTemplateNotFoundError creates an error for when a template cannot be found
func NewTemplateNotFoundError(templateName string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Template '%s' was not found.", templateName)
	details := fmt.Sprintf("Template directory 'templates/%s' does not exist or is not accessible", templateName)

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check that the template name is spelled correctly\n"
	suggestions += "• Run 'open-workbench-cli ui' to see available templates\n"
	suggestions += "• Ensure the template is properly installed\n"

	return &TemplateError{
		Type:        ErrorTypeTemplateNotFound,
		Message:     message + suggestions,
		Details:     details,
		Template:    templateName,
		OriginalErr: originalErr,
	}
}

// NewInvalidManifestError creates an error for malformed template.json files
func NewInvalidManifestError(templateName, details string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Template '%s' has an invalid configuration file.", templateName)
	message += "\n\nThe template.json file is malformed or contains invalid settings."

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check the template.json file for syntax errors\n"
	suggestions += "• Verify that all required fields are present\n"
	suggestions += "• Ensure parameter definitions are valid\n"

	return &TemplateError{
		Type:        ErrorTypeInvalidManifest,
		Message:     message + suggestions,
		Details:     details,
		Template:    templateName,
		OriginalErr: originalErr,
	}
}

// NewParameterValidationError creates an error for parameter validation failures
func NewParameterValidationError(paramName, value, reason string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Invalid value '%s' for parameter '%s'.", value, paramName)
	message += fmt.Sprintf("\n\nReason: %s", reason)

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check the parameter requirements\n"
	suggestions += "• Ensure the value matches the expected format\n"
	suggestions += "• Try a different value that meets the requirements\n"

	return &TemplateError{
		Type:        ErrorTypeParameterValidation,
		Message:     message + suggestions,
		Details:     fmt.Sprintf("Parameter '%s' validation failed: %s", paramName, reason),
		Parameter:   paramName,
		OriginalErr: originalErr,
	}
}

// NewFileSystemError creates an error for file system operation failures
func NewFileSystemError(operation, filePath string, originalErr error) *TemplateError {
	var message string
	var suggestions string

	switch {
	case os.IsNotExist(originalErr):
		message = fmt.Sprintf("File or directory not found: %s", filePath)
		suggestions = "\n\nPossible solutions:\n"
		suggestions += "• Check that the file path is correct\n"
		suggestions += "• Ensure the file exists and is accessible\n"
		suggestions += "• Verify file permissions\n"

	case os.IsPermission(originalErr):
		message = fmt.Sprintf("Permission denied: %s", filePath)
		suggestions = "\n\nPossible solutions:\n"
		suggestions += "• Run the command with appropriate permissions\n"
		suggestions += "• Check file and directory permissions\n"
		suggestions += "• Ensure you have write access to the target directory\n"

	default:
		message = fmt.Sprintf("File system error during %s: %s", operation, filePath)
		suggestions = "\n\nPossible solutions:\n"
		suggestions += "• Check available disk space\n"
		suggestions += "• Ensure the directory is writable\n"
		suggestions += "• Try running the command again\n"
	}

	return &TemplateError{
		Type:        ErrorTypeFileSystem,
		Message:     message + suggestions,
		Details:     fmt.Sprintf("Operation: %s, Path: %s", operation, filePath),
		FilePath:    filePath,
		OriginalErr: originalErr,
	}
}

// NewCommandExecutionError creates an error for failed command execution
func NewCommandExecutionError(command, description string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Command failed: %s", description)
	message += fmt.Sprintf("\n\nCommand: %s", command)

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check that the required tools are installed\n"
	suggestions += "• Verify that you have the necessary permissions\n"
	suggestions += "• Try running the command manually to see the error\n"
	suggestions += "• Check your internet connection if the command requires it\n"

	return &TemplateError{
		Type:        ErrorTypeCommandExecution,
		Message:     message + suggestions,
		Details:     fmt.Sprintf("Command: %s, Description: %s", command, description),
		Command:     command,
		OriginalErr: originalErr,
	}
}

// NewTemplateProcessingError creates an error for template processing failures
func NewTemplateProcessingError(templateName, details string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Failed to process template '%s'.", templateName)
	message += "\n\nThere was an error during template processing."

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check that the template files are not corrupted\n"
	suggestions += "• Verify that all template variables are properly defined\n"
	suggestions += "• Try using a different template\n"
	suggestions += "• Report this issue to the template maintainer\n"

	return &TemplateError{
		Type:        ErrorTypeTemplateProcessing,
		Message:     message + suggestions,
		Details:     details,
		Template:    templateName,
		OriginalErr: originalErr,
	}
}

// NewPermissionError creates an error for permission-related issues
func NewPermissionError(operation, resource string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Permission denied: %s", operation)
	message += fmt.Sprintf("\n\nResource: %s", resource)

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Run the command with administrator privileges\n"
	suggestions += "• Check file and directory permissions\n"
	suggestions += "• Ensure you have write access to the target location\n"
	suggestions += "• Try running in a different directory\n"

	return &TemplateError{
		Type:        ErrorTypePermission,
		Message:     message + suggestions,
		Details:     fmt.Sprintf("Operation: %s, Resource: %s", operation, resource),
		OriginalErr: originalErr,
	}
}

// NewNetworkError creates an error for network-related issues
func NewNetworkError(operation string, originalErr error) *TemplateError {
	message := fmt.Sprintf("Network error during %s", operation)

	suggestions := "\n\nPossible solutions:\n"
	suggestions += "• Check your internet connection\n"
	suggestions += "• Verify that the required services are accessible\n"
	suggestions += "• Try again in a few moments\n"
	suggestions += "• Check firewall and proxy settings\n"

	return &TemplateError{
		Type:        ErrorTypeNetwork,
		Message:     message + suggestions,
		Details:     fmt.Sprintf("Operation: %s", operation),
		OriginalErr: originalErr,
	}
}

// FormatErrorForUser formats an error for user display, providing
// context and suggestions based on the error type.
func FormatErrorForUser(err error) string {
	if templateErr, ok := err.(*TemplateError); ok {
		return templateErr.Message
	}

	// For generic errors, provide a basic user-friendly message
	return fmt.Sprintf("An unexpected error occurred: %v\n\nPlease try again or contact support if the problem persists.", err)
}

// IsTemplateError checks if an error is a structured template error
func IsTemplateError(err error) bool {
	_, ok := err.(*TemplateError)
	return ok
}

// GetErrorType returns the error type for a template error
func GetErrorType(err error) ErrorType {
	if templateErr, ok := err.(*TemplateError); ok {
		return templateErr.Type
	}
	return ErrorTypeTemplateProcessing // Default for unknown errors
}

// ShouldRetry determines if an operation should be retried based on the error type
func ShouldRetry(err error) bool {
	if templateErr, ok := err.(*TemplateError); ok {
		switch templateErr.Type {
		case ErrorTypeNetwork, ErrorTypeFileSystem:
			return true
		default:
			return false
		}
	}
	return false
}
