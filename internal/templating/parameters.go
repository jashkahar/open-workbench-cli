// Package templating provides the core templating system for the Open Workbench CLI.
// This package implements dynamic template discovery, parameter processing, and
// file generation capabilities with support for conditional logic and validation.
package templating

import (
	"fmt"
	"regexp"
	"strings"
)

// ParameterProcessor handles parameter collection, validation, and processing.
// This struct manages the complete lifecycle of template parameters, including
// conditional visibility, validation, and value tracking.
type ParameterProcessor struct {
	manifest *TemplateManifest      // The template manifest containing parameter definitions
	values   map[string]interface{} // Collected parameter values
}

// NewParameterProcessor creates a new parameter processor for a template manifest.
// This function initializes a parameter processor with the given manifest and
// prepares it for parameter collection and validation.
//
// Parameters:
//   - manifest: The template manifest containing parameter definitions
//
// Returns:
//   - A pointer to the initialized ParameterProcessor
func NewParameterProcessor(manifest *TemplateManifest) *ParameterProcessor {
	return &ParameterProcessor{
		manifest: manifest,
		values:   make(map[string]interface{}),
	}
}

// GetVisibleParameters returns parameters that should be shown to the user based on conditions.
// This function filters the template's parameters based on their conditional visibility
// rules, ensuring that only relevant parameters are presented to the user.
//
// Returns:
//   - A slice of parameters that should be shown to the user
func (pp *ParameterProcessor) GetVisibleParameters() []Parameter {
	var visibleParams []Parameter

	// Check each parameter's visibility based on its condition
	for _, param := range pp.manifest.Parameters {
		if pp.shouldShowParameter(param) {
			visibleParams = append(visibleParams, param)
		}
	}

	return visibleParams
}

// shouldShowParameter evaluates if a parameter should be shown based on its condition.
// This function evaluates the conditional logic for a parameter to determine
// whether it should be visible to the user based on previously collected values.
//
// Parameters:
//   - param: The parameter to evaluate for visibility
//
// Returns:
//   - true if the parameter should be shown, false otherwise
func (pp *ParameterProcessor) shouldShowParameter(param Parameter) bool {
	// If no condition is specified, always show the parameter
	if param.Condition == "" {
		return true
	}

	// Evaluate the condition against current parameter values
	result, err := pp.evaluateCondition(param.Condition)
	if err != nil {
		// If we can't evaluate the condition, show the parameter to be safe
		return true
	}

	return result
}

// evaluateCondition evaluates a condition string against current parameter values.
// This function implements a simple condition evaluator that supports equality
// and inequality comparisons for boolean and string values.
//
// Parameters:
//   - condition: The condition string to evaluate (e.g., "IncludeTesting == true")
//
// Returns:
//   - true if the condition is met, false otherwise
//   - An error if the condition cannot be evaluated
func (pp *ParameterProcessor) evaluateCondition(condition string) (bool, error) {
	// Simple condition evaluation for now
	// In the future, this could be enhanced with a proper expression parser

	condition = strings.TrimSpace(condition)

	// Handle simple equality conditions like "IncludeTesting == true"
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}

		paramName := strings.TrimSpace(parts[0])
		expectedValue := strings.TrimSpace(parts[1])

		// Get the actual value for the parameter
		actualValue, exists := pp.values[paramName]
		if !exists {
			return false, nil
		}

		// Convert expected value to appropriate type and compare
		switch expectedValue {
		case "true":
			return actualValue == true, nil
		case "false":
			return actualValue == false, nil
		default:
			// String comparison
			return fmt.Sprintf("%v", actualValue) == expectedValue, nil
		}
	}

	// Handle inequality conditions like "TestingFramework != 'Jest'"
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}

		paramName := strings.TrimSpace(parts[0])
		expectedValue := strings.TrimSpace(parts[1])

		// Get the actual value for the parameter
		actualValue, exists := pp.values[paramName]
		if !exists {
			return true, nil // If parameter doesn't exist, condition is true
		}

		// Remove quotes from expected value
		expectedValue = strings.Trim(expectedValue, "'\"")

		return fmt.Sprintf("%v", actualValue) != expectedValue, nil
	}

	return false, fmt.Errorf("unsupported condition format: %s", condition)
}

// SetValue sets a parameter value and updates the processor state.
// This function stores a parameter value in the processor's internal state,
// which is used for conditional logic evaluation.
//
// Parameters:
//   - paramName: The name of the parameter to set
//   - value: The value to store for the parameter
func (pp *ParameterProcessor) SetValue(paramName string, value interface{}) {
	pp.values[paramName] = value
}

// GetValue retrieves a parameter value from the processor state.
// This function returns the stored value for a parameter, if it exists.
//
// Parameters:
//   - paramName: The name of the parameter to retrieve
//
// Returns:
//   - The parameter value and a boolean indicating if the value exists
func (pp *ParameterProcessor) GetValue(paramName string) (interface{}, bool) {
	value, exists := pp.values[paramName]
	return value, exists
}

// ValidateParameter validates a parameter value against its validation rules.
// This function performs type checking and custom validation based on the
// parameter's configuration, including regex validation for string parameters.
//
// Parameters:
//   - param: The parameter definition containing validation rules
//   - value: The value to validate
//
// Returns:
//   - An error if validation fails, nil if validation passes
func (pp *ParameterProcessor) ValidateParameter(param Parameter, value interface{}) error {
	// Perform type validation based on parameter type
	switch param.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("parameter %s expects string value", param.Name)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("parameter %s expects boolean value", param.Name)
		}
	case "select":
		if strValue, ok := value.(string); ok {
			return pp.validateSelectValue(param, strValue)
		}
		return fmt.Errorf("parameter %s expects string value", param.Name)
	case "multiselect":
		if strValues, ok := value.([]string); ok {
			return pp.validateMultiSelectValue(param, strValues)
		}
		return fmt.Errorf("parameter %s expects array of strings", param.Name)
	}

	// Apply custom validation for string parameters
	if param.Type == "string" && param.Validation != nil {
		strValue, _ := value.(string)
		return pp.validateStringValue(param, strValue)
	}

	return nil
}

// validateSelectValue validates a select parameter value.
// This function ensures that the selected value is one of the valid options
// defined in the parameter's options list.
//
// Parameters:
//   - param: The parameter definition containing the options list
//   - value: The selected value to validate
//
// Returns:
//   - An error if the value is not a valid option
func (pp *ParameterProcessor) validateSelectValue(param Parameter, value string) error {
	// Check if the value is in the list of valid options
	for _, option := range param.Options {
		if option == value {
			return nil
		}
	}
	return fmt.Errorf("parameter %s: value '%s' is not a valid option", param.Name, value)
}

// validateMultiSelectValue validates a multiselect parameter value.
// This function ensures that all selected values are valid options
// defined in the parameter's options list.
//
// Parameters:
//   - param: The parameter definition containing the options list
//   - values: The selected values to validate
//
// Returns:
//   - An error if any value is not a valid option
func (pp *ParameterProcessor) validateMultiSelectValue(param Parameter, values []string) error {
	// Check each selected value against the valid options
	for _, value := range values {
		found := false
		for _, option := range param.Options {
			if option == value {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("parameter %s: value '%s' is not a valid option", param.Name, value)
		}
	}
	return nil
}

// validateStringValue validates a string parameter value against regex patterns.
// This function applies custom regex validation to string parameters,
// providing detailed error messages for validation failures.
//
// Parameters:
//   - param: The parameter definition containing validation rules
//   - value: The string value to validate
//
// Returns:
//   - An error if validation fails, nil if validation passes
func (pp *ParameterProcessor) validateStringValue(param Parameter, value string) error {
	// Skip validation if no regex pattern is specified
	if param.Validation == nil || param.Validation.Regex == "" {
		return nil
	}

	// Compile and apply the regex pattern
	matched, err := regexp.MatchString(param.Validation.Regex, value)
	if err != nil {
		return fmt.Errorf("invalid regex pattern in parameter %s: %w", param.Name, err)
	}

	// Return appropriate error message if validation fails
	if !matched {
		if param.Validation.ErrorMessage != "" {
			return fmt.Errorf("parameter %s: %s", param.Name, param.Validation.ErrorMessage)
		}
		return fmt.Errorf("parameter %s: value does not match required pattern", param.Name)
	}

	return nil
}

// GetParameterGroups returns parameters organized by their groups.
// This function organizes parameters into logical groups for better user experience,
// using the group field from parameter definitions.
//
// Returns:
//   - A map of group names to parameter slices
func (pp *ParameterProcessor) GetParameterGroups() map[string][]Parameter {
	groups := make(map[string][]Parameter)

	// Get visible parameters and organize them by group
	for _, param := range pp.GetVisibleParameters() {
		group := param.Group
		if group == "" {
			group = "General" // Default group for ungrouped parameters
		}
		groups[group] = append(groups[group], param)
	}

	return groups
}

// GetRequiredParameters returns only the required parameters.
// This function filters the visible parameters to return only those
// that are marked as required in their definition.
//
// Returns:
//   - A slice of required parameters
func (pp *ParameterProcessor) GetRequiredParameters() []Parameter {
	var requiredParams []Parameter

	// Filter visible parameters to include only required ones
	for _, param := range pp.GetVisibleParameters() {
		if param.Required {
			requiredParams = append(requiredParams, param)
		}
	}

	return requiredParams
}

// GetAllValues returns all collected parameter values.
// This function returns a copy of all parameter values collected so far,
// which can be used for template processing or debugging.
//
// Returns:
//   - A map of parameter names to their collected values
func (pp *ParameterProcessor) GetAllValues() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range pp.values {
		result[k] = v
	}
	return result
}

// IsComplete checks if all required parameters have been provided.
// This function verifies that all required parameters have been collected
// and are ready for template processing.
//
// Returns:
//   - true if all required parameters are complete, false otherwise
func (pp *ParameterProcessor) IsComplete() bool {
	// Check each required parameter to ensure it has been provided
	for _, param := range pp.GetRequiredParameters() {
		if _, exists := pp.values[param.Name]; !exists {
			return false
		}
	}
	return true
}
