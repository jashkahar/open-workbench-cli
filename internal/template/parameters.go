package template

import (
	"fmt"
	"regexp"
	"strings"
)

// ParameterProcessor handles parameter collection, validation, and processing
type ParameterProcessor struct {
	manifest *TemplateManifest
	values   map[string]interface{}
}

// NewParameterProcessor creates a new parameter processor for a template manifest
func NewParameterProcessor(manifest *TemplateManifest) *ParameterProcessor {
	return &ParameterProcessor{
		manifest: manifest,
		values:   make(map[string]interface{}),
	}
}

// GetVisibleParameters returns parameters that should be shown to the user based on conditions
func (pp *ParameterProcessor) GetVisibleParameters() []Parameter {
	var visibleParams []Parameter

	for _, param := range pp.manifest.Parameters {
		if pp.shouldShowParameter(param) {
			visibleParams = append(visibleParams, param)
		}
	}

	return visibleParams
}

// shouldShowParameter evaluates if a parameter should be shown based on its condition
func (pp *ParameterProcessor) shouldShowParameter(param Parameter) bool {
	if param.Condition == "" {
		return true
	}

	result, err := pp.evaluateCondition(param.Condition)
	if err != nil {
		// If we can't evaluate the condition, show the parameter to be safe
		return true
	}

	return result
}

// evaluateCondition evaluates a condition string against current parameter values
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

		actualValue, exists := pp.values[paramName]
		if !exists {
			return false, nil
		}

		// Convert expected value to appropriate type
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

// SetValue sets a parameter value and updates the processor state
func (pp *ParameterProcessor) SetValue(paramName string, value interface{}) {
	pp.values[paramName] = value
}

// GetValue retrieves a parameter value
func (pp *ParameterProcessor) GetValue(paramName string) (interface{}, bool) {
	value, exists := pp.values[paramName]
	return value, exists
}

// ValidateParameter validates a parameter value against its validation rules
func (pp *ParameterProcessor) ValidateParameter(param Parameter, value interface{}) error {
	// Type validation
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

	// Custom validation for string parameters
	if param.Type == "string" && param.Validation != nil {
		strValue, _ := value.(string)
		return pp.validateStringValue(param, strValue)
	}

	return nil
}

// validateSelectValue validates a select parameter value
func (pp *ParameterProcessor) validateSelectValue(param Parameter, value string) error {
	for _, option := range param.Options {
		if option == value {
			return nil
		}
	}
	return fmt.Errorf("parameter %s: value '%s' is not a valid option", param.Name, value)
}

// validateMultiSelectValue validates a multiselect parameter value
func (pp *ParameterProcessor) validateMultiSelectValue(param Parameter, values []string) error {
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

// validateStringValue validates a string parameter value against regex
func (pp *ParameterProcessor) validateStringValue(param Parameter, value string) error {
	if param.Validation == nil || param.Validation.Regex == "" {
		return nil
	}

	matched, err := regexp.MatchString(param.Validation.Regex, value)
	if err != nil {
		return fmt.Errorf("invalid regex pattern in parameter %s: %w", param.Name, err)
	}

	if !matched {
		if param.Validation.ErrorMessage != "" {
			return fmt.Errorf("parameter %s: %s", param.Name, param.Validation.ErrorMessage)
		}
		return fmt.Errorf("parameter %s: value does not match required pattern", param.Name)
	}

	return nil
}

// GetParameterGroups returns parameters organized by their groups
func (pp *ParameterProcessor) GetParameterGroups() map[string][]Parameter {
	groups := make(map[string][]Parameter)

	for _, param := range pp.GetVisibleParameters() {
		group := param.Group
		if group == "" {
			group = "General"
		}
		groups[group] = append(groups[group], param)
	}

	return groups
}

// GetRequiredParameters returns only the required parameters
func (pp *ParameterProcessor) GetRequiredParameters() []Parameter {
	var requiredParams []Parameter

	for _, param := range pp.GetVisibleParameters() {
		if param.Required {
			requiredParams = append(requiredParams, param)
		}
	}

	return requiredParams
}

// GetAllValues returns all collected parameter values
func (pp *ParameterProcessor) GetAllValues() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range pp.values {
		result[k] = v
	}
	return result
}

// IsComplete checks if all required parameters have been provided
func (pp *ParameterProcessor) IsComplete() bool {
	for _, param := range pp.GetRequiredParameters() {
		if _, exists := pp.values[param.Name]; !exists {
			return false
		}
	}
	return true
}
