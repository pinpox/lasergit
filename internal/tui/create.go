package tui

import (
	"fmt"
	"strings"

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
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("15")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("15")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("13")).
				Padding(0, 1).
				Margin(0, 0, 1, 0)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("8")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 2).
			Margin(0, 1)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("13")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("13")).
				Padding(0, 2).
				Margin(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Margin(1, 0)

	branchInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(false)
)

type CreatePRModel struct {
	inputs   []textinput.Model
	focused  int
	topicBranch    string
	targetBranch   string
	err      error
	done     bool
	canceled bool
}

type CreatePRResult struct {
	Title       string
	Description string
	Topic       string
	Target      string
	Canceled    bool
}

func NewCreatePRModel(topicBranch, targetBranch string) CreatePRModel {
	inputs := make([]textinput.Model, 2)

	// Title input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter PR title..."
	inputs[0].Focus()
	inputs[0].CharLimit = 100
	inputs[0].Width = 60

	// Description input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Enter PR description (optional)..."
	inputs[1].CharLimit = 500
	inputs[1].Width = 60

	return CreatePRModel{
		inputs:       inputs,
		focused:      0,
		topicBranch:  topicBranch,
		targetBranch: targetBranch,
	}
}

func (m CreatePRModel) Init() tea.Cmd {
	return textinput.Blink
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
			if m.focused == len(m.inputs)-1 {
				m.done = true
				return m, tea.Quit
			}
			m.nextInput()

		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.prevInput()
			} else {
				m.nextInput()
			}

			for i := range m.inputs {
				if i == m.focused {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		for i := range m.inputs {
			m.inputs[i].Width = msg.Width - 4
		}
	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
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
		b.WriteString(focusedInputStyle.Render(m.inputs[0].View()))
	} else {
		b.WriteString(inputStyle.Render(m.inputs[0].View()))
	}
	b.WriteString("\n")

	// Description input
	b.WriteString(labelStyle.Render("Description:"))
	b.WriteString("\n")
	if m.focused == 1 {
		b.WriteString(focusedInputStyle.Render(m.inputs[1].View()))
	} else {
		b.WriteString(inputStyle.Render(m.inputs[1].View()))
	}
	b.WriteString("\n")

	// Buttons
	if m.focused == len(m.inputs) {
		b.WriteString(activeButtonStyle.Render("[ Create PR ]"))
	} else {
		b.WriteString(buttonStyle.Render("[ Create PR ]"))
	}
	b.WriteString("  ")
	b.WriteString(buttonStyle.Render("[ Cancel ]"))
	b.WriteString("\n\n")

	// Help
	b.WriteString(helpStyle.Render("tab/shift+tab: navigate • enter: next/submit • esc: cancel"))

	return b.String()
}

func (m *CreatePRModel) nextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

func (m *CreatePRModel) prevInput() {
	m.focused = (m.focused - 1 + len(m.inputs) + 1) % (len(m.inputs) + 1)
}

func (m CreatePRModel) GetResult() CreatePRResult {
	return CreatePRResult{
		Title:       m.inputs[0].Value(),
		Description: m.inputs[1].Value(),
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