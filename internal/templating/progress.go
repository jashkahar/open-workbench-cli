// Package templating provides the core templating system for the Open Workbench CLI.
// This package implements dynamic template discovery, parameter processing, and
// file generation capabilities with support for conditional logic and validation.
package templating

import (
	"fmt"
	"strings"
	"time"
)

// ProgressReporter handles progress reporting for long-running operations.
// This struct provides a consistent interface for reporting progress to users
// with spinners, progress bars, and detailed status messages.
type ProgressReporter struct {
	currentStep    string    // Current operation being performed
	totalSteps     int       // Total number of steps
	currentStepNum int       // Current step number
	startTime      time.Time // When the operation started
	spinnerChars   []string  // Characters for spinner animation
	spinnerIndex   int       // Current spinner character index
	lastUpdate     time.Time // Last progress update time
	verbose        bool      // Whether to show detailed progress
}

// NewProgressReporter creates a new progress reporter for tracking operations.
// This function initializes a progress reporter with default settings and
// prepares it for reporting progress during template operations.
//
// Parameters:
//   - totalSteps: The total number of steps in the operation
//   - verbose: Whether to show detailed progress information
//
// Returns:
//   - A pointer to the initialized ProgressReporter
func NewProgressReporter(totalSteps int, verbose bool) *ProgressReporter {
	return &ProgressReporter{
		totalSteps:   totalSteps,
		verbose:      verbose,
		startTime:    time.Now(),
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		lastUpdate:   time.Now(),
	}
}

// StartOperation begins a new operation and reports the initial status.
// This function should be called at the beginning of any long-running operation
// to set up progress tracking and provide initial feedback to the user.
//
// Parameters:
//   - operation: The name of the operation being performed
func (pr *ProgressReporter) StartOperation(operation string) {
	pr.currentStep = operation
	pr.currentStepNum = 0
	pr.startTime = time.Now()
	pr.lastUpdate = time.Now()

	if pr.verbose {
		fmt.Printf("🚀 Starting: %s\n", operation)
	} else {
		fmt.Printf("🔄 %s...", operation)
	}
}

// ReportProgress updates the progress for the current operation.
// This function should be called periodically during long-running operations
// to provide feedback to the user about the current progress.
//
// Parameters:
//   - step: The current step being performed
//   - current: The current step number
//   - total: The total number of steps
func (pr *ProgressReporter) ReportProgress(step string, current, total int) {
	pr.currentStep = step
	pr.currentStepNum = current

	// Only update if enough time has passed to avoid spam
	if time.Since(pr.lastUpdate) < 100*time.Millisecond {
		return
	}

	pr.lastUpdate = time.Now()

	if pr.verbose {
		percentage := float64(current) / float64(total) * 100
		fmt.Printf("  📋 %s (%d/%d) %.1f%%\n", step, current, total, percentage)
	} else {
		// Update spinner in place
		spinner := pr.spinnerChars[pr.spinnerIndex]
		fmt.Printf("\r%s %s (%d/%d)", spinner, step, current, total)
		pr.spinnerIndex = (pr.spinnerIndex + 1) % len(pr.spinnerChars)
	}
}

// CompleteStep marks the completion of a step and reports the result.
// This function should be called when a step is completed to provide
// feedback about the success or failure of the operation.
//
// Parameters:
//   - step: The step that was completed
//   - success: Whether the step completed successfully
//   - details: Optional details about the completion
func (pr *ProgressReporter) CompleteStep(step string, success bool, details string) {
	if pr.verbose {
		if success {
			fmt.Printf("  ✅ %s completed", step)
			if details != "" {
				fmt.Printf(": %s", details)
			}
			fmt.Println()
		} else {
			fmt.Printf("  ❌ %s failed", step)
			if details != "" {
				fmt.Printf(": %s", details)
			}
			fmt.Println()
		}
	} else {
		if success {
			fmt.Printf("\r✅ %s completed\n", step)
		} else {
			fmt.Printf("\r❌ %s failed\n", step)
		}
	}
}

// CompleteOperation marks the completion of the entire operation.
// This function should be called when all steps are completed to provide
// a summary of the operation and clear any pending progress indicators.
//
// Parameters:
//   - success: Whether the entire operation completed successfully
//   - summary: Optional summary of the operation results
func (pr *ProgressReporter) CompleteOperation(success bool, summary string) {
	duration := time.Since(pr.startTime)

	if pr.verbose {
		if success {
			fmt.Printf("🎉 Operation completed successfully in %v\n", duration.Round(time.Millisecond))
		} else {
			fmt.Printf("💥 Operation failed after %v\n", duration.Round(time.Millisecond))
		}
		if summary != "" {
			fmt.Printf("📋 Summary: %s\n", summary)
		}
	} else {
		if success {
			fmt.Printf("✅ Completed in %v\n", duration.Round(time.Millisecond))
		} else {
			fmt.Printf("❌ Failed after %v\n", duration.Round(time.Millisecond))
		}
	}
}

// ReportFileOperation reports progress for file operations.
// This function provides specialized progress reporting for file copying,
// template processing, and other file-related operations.
//
// Parameters:
//   - operation: The file operation being performed
//   - filePath: The path of the file being processed
//   - current: Current file number
//   - total: Total number of files
func (pr *ProgressReporter) ReportFileOperation(operation, filePath string, current, total int) {
	step := fmt.Sprintf("%s: %s", operation, filePath)
	pr.ReportProgress(step, current, total)
}

// ReportCommandExecution reports progress for command execution.
// This function provides specialized progress reporting for shell commands
// and other external process execution.
//
// Parameters:
//   - command: The command being executed
//   - description: Human-readable description of the command
func (pr *ProgressReporter) ReportCommandExecution(command, description string) {
	if pr.verbose {
		fmt.Printf("  🔧 Executing: %s\n", description)
		fmt.Printf("    Command: %s\n", command)
	} else {
		fmt.Printf("🔧 %s...", description)
	}
}

// ReportCommandResult reports the result of a command execution.
// This function should be called after a command completes to report
// success or failure with appropriate details.
//
// Parameters:
//   - command: The command that was executed
//   - success: Whether the command executed successfully
//   - output: Optional command output or error message
func (pr *ProgressReporter) ReportCommandResult(command string, success bool, output string) {
	if pr.verbose {
		if success {
			fmt.Printf("  ✅ Command completed successfully")
			if output != "" {
				fmt.Printf(" (output: %s)", strings.TrimSpace(output))
			}
			fmt.Println()
		} else {
			fmt.Printf("  ❌ Command failed")
			if output != "" {
				fmt.Printf(" (error: %s)", strings.TrimSpace(output))
			}
			fmt.Println()
		}
	} else {
		if success {
			fmt.Printf("✅ %s completed\n", command)
		} else {
			fmt.Printf("❌ %s failed\n", command)
		}
	}
}

// ReportTemplateProcessing reports progress for template processing operations.
// This function provides specialized progress reporting for template variable
// substitution and conditional processing.
//
// Parameters:
//   - templateName: The name of the template being processed
//   - current: Current template file number
//   - total: Total number of template files
func (pr *ProgressReporter) ReportTemplateProcessing(templateName string, current, total int) {
	step := fmt.Sprintf("Processing template files (%s)", templateName)
	pr.ReportProgress(step, current, total)
}

// ReportParameterCollection reports progress for parameter collection.
// This function provides specialized progress reporting for user input
// collection during the scaffolding process.
//
// Parameters:
//   - parameterName: The name of the parameter being collected
//   - current: Current parameter number
//   - total: Total number of parameters
func (pr *ProgressReporter) ReportParameterCollection(parameterName string, current, total int) {
	step := fmt.Sprintf("Collecting parameter: %s", parameterName)
	pr.ReportProgress(step, current, total)
}

// GetElapsedTime returns the elapsed time since the operation started.
// This function provides timing information for performance monitoring
// and user feedback.
//
// Returns:
//   - The elapsed time as a duration
func (pr *ProgressReporter) GetElapsedTime() time.Duration {
	return time.Since(pr.startTime)
}

// IsVerbose returns whether the progress reporter is in verbose mode.
// This function allows other components to adjust their behavior based
// on the verbosity setting.
//
// Returns:
//   - true if verbose mode is enabled, false otherwise
func (pr *ProgressReporter) IsVerbose() bool {
	return pr.verbose
}
