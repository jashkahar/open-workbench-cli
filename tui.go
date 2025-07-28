// Package main provides the Terminal User Interface (TUI) functionality for the
// Open Workbench CLI. This file contains the TUI implementation using the
// Bubble Tea framework for creating interactive terminal interfaces.
package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jashkahar/open-workbench-cli/internal/templating"
)

// Define visual styles for the TUI components
var (
	// titleStyle defines the styling for the main title text
	titleStyle = lipgloss.NewStyle().MarginLeft(2)

	// itemStyle defines the styling for list items in the template selection
	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	// quitTextStyle defines the styling for the quit/confirmation text
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// item represents a template in the TUI list selection.
// This struct implements the list.Item interface required by the Bubble Tea list component.
type item struct {
	name        string // The template name (e.g., "nextjs-golden-path")
	description string // The template description for display
}

// FilterValue returns the value used for filtering in the list.
// This method is required by the list.Item interface.
func (i item) FilterValue() string {
	return i.name
}

// Title returns the display title for the list item.
// This method is required by the list.Item interface.
func (i item) Title() string {
	return i.name
}

// Description returns the display description for the list item.
// This method is required by the list.Item interface.
func (i item) Description() string {
	return i.description
}

// model represents the state of the TUI application.
// This struct holds all the data needed to render and update the interface.
type model struct {
	list     list.Model // The list component for template selection
	choice   string     // The selected template name
	quitting bool       // Whether the user has chosen to quit
}

// newModel initializes a new TUI model with available templates.
// This function discovers all available templates and creates a list
// component for user selection.
//
// Returns:
//   - A pointer to the initialized model
//   - An error if template discovery fails
func newModel() (*model, error) {
	var items []list.Item

	// Discover all available templates using the dynamic template system
	templates, err := templating.DiscoverTemplates(templatesFS)
	if err != nil {
		return nil, fmt.Errorf("could not discover templates: %w", err)
	}

	// Convert discovered templates to list items for display
	for _, tmpl := range templates {
		items = append(items, item{
			name:        tmpl.Name,
			description: tmpl.Description,
		})
	}

	// Create and configure the list component
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Please choose a project template"
	l.SetSpinner(spinner.Dot) // Use dot spinner for loading animation

	return &model{list: l}, nil
}

// Init is the first command that is run when the program starts.
// This method is required by the tea.Model interface and is called
// when the TUI program begins execution.
//
// Returns:
//   - A tea.Cmd that will be executed on initialization
func (m model) Init() tea.Cmd {
	return nil // No initial command needed
}

// Update handles all incoming events and updates the model accordingly.
// This method is required by the tea.Model interface and processes
// all user input and system events.
//
// Parameters:
//   - msg: The incoming message/event to process
//
// Returns:
//   - The updated model
//   - A tea.Cmd to execute (if any)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize events
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		// Handle keyboard input events
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			// Handle quit commands
			m.quitting = true
			return m, tea.Quit

		case "enter":
			// Handle template selection
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.name
			}
			return m, tea.Quit
		}
	}

	// Update the list component with the current message
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the current state of the TUI.
// This method is required by the tea.Model interface and returns
// the string representation of the current UI state.
//
// Returns:
//   - The string representation of the current UI state
func (m model) View() string {
	if m.choice != "" {
		// Display confirmation message when template is selected
		return quitTextStyle.Render(fmt.Sprintf("Great! Scaffolding project from template: %s", m.choice))
	}
	if m.quitting {
		// Display cancellation message when user quits
		return quitTextStyle.Render("Operation cancelled.")
	}
	// Display the main list interface
	return "\n" + m.list.View()
}

// runTUI launches the Terminal User Interface for template selection.
// This function creates and runs the Bubble Tea program, providing
// an interactive interface for users to select from available templates.
//
// Returns:
//   - The name of the selected template (empty string if none selected)
//   - An error if the TUI fails to start or run
func runTUI() (string, error) {
	// Initialize the TUI model with available templates
	m, err := newModel()
	if err != nil {
		return "", err
	}

	// Create and configure the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the TUI program and get the final model state
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running TUI: %w", err)
	}

	// Extract the selected template name from the final model state
	if m, ok := finalModel.(model); ok {
		return m.choice, nil
	}

	// Return empty string if no template was selected
	return "", nil
}
