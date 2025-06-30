package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Bold(true).
			Margin(1, 0)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("13")).
				Padding(0, 1).
				Margin(0, 0, 1, 0)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("13")).
				Bold(true).
				Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Margin(1, 0)

	branchInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(false)
)

type CreatePRModel struct {
	titleInput   textinput.Model
	descInput    textarea.Model
	focused      int
	topicBranch  string
	targetBranch string
	err          error
	done         bool
	canceled     bool
}

type CreatePRResult struct {
	Title       string
	Description string
	Topic       string
	Target      string
	Canceled    bool
}

func NewCreatePRModel(topicBranch, targetBranch string) CreatePRModel {
	// Title input
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter PR title..."
	titleInput.Focus()
	titleInput.CharLimit = 100
	titleInput.Width = 60

	// Description textarea (multiline)
	descInput := textarea.New()
	descInput.Placeholder = "Enter PR description (optional)..."
	descInput.CharLimit = 1000
	descInput.SetWidth(60)
	descInput.SetHeight(4)

	return CreatePRModel{
		titleInput:   titleInput,
		descInput:    descInput,
		focused:      0,
		topicBranch:  topicBranch,
		targetBranch: targetBranch,
	}
}

func (m CreatePRModel) Init() tea.Cmd {
	return nil
}

func (m CreatePRModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.canceled = true
			m.done = true
			return m, tea.Quit

		case "enter":
			// Handle button actions
			if m.focused == 2 {
				// Create PR button
				m.done = true
				return m, tea.Quit
			} else if m.focused == 3 {
				// Cancel button
				m.canceled = true
				m.done = true
				return m, tea.Quit
			}
			// If we're on title field, move to description
			if m.focused == 0 {
				m.nextInput()
			}
			// If we're on description field, let Enter add newline (handled by textarea)

		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.prevInput()
			} else {
				m.nextInput()
			}

			return m, tea.Batch(cmds...)

		case "ctrl+enter":
			// Ctrl+Enter submits from anywhere
			m.done = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.titleInput.Width = msg.Width - 4
		m.descInput.SetWidth(msg.Width - 4)
	}

	// Only update the currently focused input
	var cmd tea.Cmd
	if m.focused == 0 {
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focused == 1 {
		m.descInput, cmd = m.descInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m CreatePRModel) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("🚀 Create Pull Request"))
	b.WriteString("\n\n")

	// Topic and target branch info
	b.WriteString(labelStyle.Render("Topic Branch: "))
	b.WriteString(branchInfoStyle.Render(m.topicBranch))
	b.WriteString("\n")
	b.WriteString(labelStyle.Render("Target Branch: "))
	b.WriteString(branchInfoStyle.Render(m.targetBranch))
	b.WriteString("\n\n")

	// Title input
	b.WriteString(labelStyle.Render("Title:"))
	b.WriteString("\n")
	if m.focused == 0 {
		b.WriteString(focusedInputStyle.Render(m.titleInput.View()))
	} else {
		b.WriteString(inputStyle.Render(m.titleInput.View()))
	}
	b.WriteString("\n")

	// Description textarea
	b.WriteString(labelStyle.Render("Description:"))
	b.WriteString("\n")
	if m.focused == 1 {
		b.WriteString(focusedInputStyle.Render(m.descInput.View()))
	} else {
		b.WriteString(inputStyle.Render(m.descInput.View()))
	}
	b.WriteString("\n")

	// Buttons
	var createButton, cancelButton string
	if m.focused == 2 {
		createButton = activeButtonStyle.Render("Create PR")
	} else {
		createButton = buttonStyle.Render("Create PR")
	}
	if m.focused == 3 {
		cancelButton = activeButtonStyle.Render("Cancel")
	} else {
		cancelButton = buttonStyle.Render("Cancel")
	}
	
	buttonsLine := lipgloss.JoinHorizontal(lipgloss.Top, createButton, "   ", cancelButton)
	b.WriteString(buttonsLine)
	b.WriteString("\n\n")

	// Help
	b.WriteString(helpStyle.Render("tab: navigate • enter: newline in description • ctrl+enter: submit • esc: cancel"))

	return b.String()
}

func (m *CreatePRModel) nextInput() {
	// Blur current input
	if m.focused == 0 {
		m.titleInput.Blur()
	} else if m.focused == 1 {
		m.descInput.Blur()
	}
	
	m.focused = (m.focused + 1) % 4 // 0: title, 1: desc, 2: create button, 3: cancel button
	
	// Focus new input
	if m.focused == 0 {
		m.titleInput.Focus()
	} else if m.focused == 1 {
		m.descInput.Focus()
	}
}

func (m *CreatePRModel) prevInput() {
	// Blur current input
	if m.focused == 0 {
		m.titleInput.Blur()
	} else if m.focused == 1 {
		m.descInput.Blur()
	}
	
	m.focused = (m.focused - 1 + 4) % 4 // 0: title, 1: desc, 2: create button, 3: cancel button
	
	// Focus new input
	if m.focused == 0 {
		m.titleInput.Focus()
	} else if m.focused == 1 {
		m.descInput.Focus()
	}
}

func (m CreatePRModel) GetResult() CreatePRResult {
	return CreatePRResult{
		Title:       m.titleInput.Value(),
		Description: m.descInput.Value(),
		Topic:       m.topicBranch,
		Target:      m.targetBranch,
		Canceled:    m.canceled,
	}
}

func ShowCreatePRDialog(topicBranch, targetBranch string) (*CreatePRResult, error) {
	model := NewCreatePRModel(topicBranch, targetBranch)
	program := tea.NewProgram(model)
	
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	createModel := finalModel.(CreatePRModel)
	result := createModel.GetResult()
	
	if result.Canceled {
		return nil, fmt.Errorf("canceled by user")
	}

	if result.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return &result, nil
}