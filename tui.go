package main

import (
	"fmt"
	"io/fs"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define styles
var (
	titleStyle    = lipgloss.NewStyle().MarginLeft(2)
	itemStyle     = lipgloss.NewStyle().PaddingLeft(4)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// item represents a template in our list.
type item string

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "A project template" }

// model is the state of our TUI application.
type model struct {
	list     list.Model
	choice   string
	quitting bool
}

// newModel initializes the TUI model.
func newModel() (*model, error) {
	var items []list.Item

	// Read directories directly from the embedded templatesFS variable
	templateEntries, err := fs.ReadDir(templatesFS, "templates")
	if err != nil {
		// This error would mean the embedding failed during build, which is a critical issue.
		return nil, fmt.Errorf("could not read embedded templates directory: %w", err)
	}

	for _, entry := range templateEntries {
		if entry.IsDir() {
			items = append(items, item(entry.Name()))
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Please choose a project template"
	l.SetSpinner(spinner.Dot)

	return &model{list: l}, nil
}

// Init is the first command that is run when the program starts.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles all incoming events and updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Great! Scaffolding project from template: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Operation cancelled.")
	}
	return "\n" + m.list.View()
}

// runTUI is the function we will call from main.go to start the TUI.
func runTUI() (string, error) {
	m, err := newModel()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running TUI: %w", err)
	}

	// Type assertion to get the final state of the model
	if m, ok := finalModel.(model); ok {
		return m.choice, nil
	}

	return "", nil
}
