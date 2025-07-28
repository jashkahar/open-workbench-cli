// Package main provides shared type definitions for the Open Workbench CLI.
// This file contains data structures used throughout the application for
// template processing, parameter handling, and configuration management.
package main

// TemplateManifest represents the structure of a template.json file.
// This struct defines the complete configuration for a template, including
// metadata, parameters, validation rules, and post-scaffolding actions.
//
// The manifest is loaded from JSON files in each template directory and
// provides the blueprint for template processing and parameter collection.
type TemplateManifest struct {
	Name         string        `json:"name"`                   // Display name for the template
	Description  string        `json:"description"`            // Human-readable description
	Parameters   []Parameter   `json:"parameters"`             // List of parameters to collect
	PostScaffold *PostScaffold `json:"postScaffold,omitempty"` // Post-processing actions
}

// Parameter represents a single parameter that the user needs to provide.
// This struct defines the configuration for collecting user input during
// the template scaffolding process, including validation rules and UI options.
type Parameter struct {
	Name       string      `json:"name"`                 // Unique parameter identifier
	Prompt     string      `json:"prompt"`               // User-facing question text
	HelpText   string      `json:"helpText,omitempty"`   // Additional help information
	Group      string      `json:"group,omitempty"`      // Group for organizing parameters
	Type       string      `json:"type"`                 // Parameter type (string, boolean, select, multiselect)
	Required   bool        `json:"required,omitempty"`   // Whether parameter is mandatory
	Default    any         `json:"default,omitempty"`    // Default value for the parameter
	Options    []string    `json:"options,omitempty"`    // Available options for select/multiselect
	Condition  string      `json:"condition,omitempty"`  // Conditional visibility rule
	Validation *Validation `json:"validation,omitempty"` // Validation rules for the parameter
}

// Validation represents validation rules for string parameters.
// This struct defines regex patterns and error messages for validating
// user input during parameter collection.
type Validation struct {
	Regex        string `json:"regex"`        // Regular expression pattern for validation
	ErrorMessage string `json:"errorMessage"` // Custom error message for validation failures
}

// PostScaffold represents actions to perform after scaffolding.
// This struct defines cleanup and setup actions that should be executed
// after the main template processing is complete.
type PostScaffold struct {
	FilesToDelete []FileAction    `json:"filesToDelete,omitempty"` // Files to delete based on conditions
	Commands      []CommandAction `json:"commands,omitempty"`      // Commands to execute after scaffolding
}

// FileAction represents a file or directory to delete based on a condition.
// This struct defines conditional file deletion rules that can be applied
// after template processing based on user parameter choices.
type FileAction struct {
	Path      string `json:"path"`      // Path to the file or directory to delete
	Condition string `json:"condition"` // Condition that must be true for deletion
}

// CommandAction represents a command to execute after scaffolding.
// This struct defines shell commands that should be run after template
// processing, such as dependency installation or git initialization.
type CommandAction struct {
	Command     string `json:"command"`             // Shell command to execute
	Description string `json:"description"`         // Human-readable description of the command
	Condition   string `json:"condition,omitempty"` // Optional condition for execution
}

// TemplateInfo represents metadata about a discovered template.
// This struct contains information about a template discovered during
// the template discovery process, including its manifest and metadata.
type TemplateInfo struct {
	Name        string            // The template name (directory name)
	Description string            // Human-readable description from manifest
	Path        string            // Path to the template directory
	Manifest    *TemplateManifest // Parsed template manifest
}

// ParameterValue represents a collected parameter value with metadata.
// This struct is used internally to track parameter values during
// the collection and validation process.
type ParameterValue struct {
	Parameter Parameter // The parameter definition
	Value     any       // The collected value
	Valid     bool      // Whether the value passed validation
	Error     string    // Error message if validation failed
}
