package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"code.gitea.io/sdk/gitea"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Bold(true).
			Align(lipgloss.Center)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("13")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Margin(1, 0, 0, 0)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	authorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true)
)

type ListPRModel struct {
	table    table.Model
	prs      []*gitea.PullRequest
	owner    string
	repo     string
	selected int
	action   string
	done     bool
}

type ListPRResult struct {
	SelectedPR *gitea.PullRequest
	Action     string // "view", "checkout", "quit"
}

func NewListPRModel(prs []*gitea.PullRequest, owner, repo string) ListPRModel {
	columns := []table.Column{
		{Title: "PR", Width: 6},
		{Title: "Title", Width: 50},
		{Title: "Author", Width: 15},
		{Title: "Status", Width: 10},
		{Title: "Updated", Width: 12},
	}

	rows := make([]table.Row, len(prs))
	for i, pr := range prs {
		status := "Open"
		if pr.State == gitea.StateClosed {
			status = "Closed"
		} else if pr.Merged != nil && !pr.Merged.IsZero() {
			status = "Merged"
		}

		updatedTime := ""
		if pr.Updated != nil {
			updatedTime = pr.Updated.Format("2006-01-02")
		}

		title := pr.Title
		if len(title) > 47 {
			title = title[:44] + "..."
		}

		rows[i] = table.Row{
			fmt.Sprintf("#%d", pr.Index),
			title,
			pr.Poster.UserName,
			status,
			updatedTime,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle
	t.SetStyles(s)

	return ListPRModel{
		table: t,
		prs:   prs,
		owner: owner,
		repo:  repo,
	}
}

func (m ListPRModel) Init() tea.Cmd {
	return nil
}

func (m ListPRModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.action = "quit"
			m.done = true
			return m, tea.Quit

		case "enter":
			m.selected = m.table.Cursor()
			m.action = "checkout"
			m.done = true
			return m, tea.Quit

		case "v":
			m.selected = m.table.Cursor()
			m.action = "view"
			m.done = true
			return m, tea.Quit

		case "r":
			m.action = "refresh"
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ListPRModel) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	// Header
	title := fmt.Sprintf("📋 Pull Requests - %s/%s", m.owner, m.repo)
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Table
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteString("\n")

	// Info section
	if len(m.prs) > 0 {
		selected := m.prs[m.table.Cursor()]
		b.WriteString(infoStyle.Render("Selected PR Details:"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Title: %s\n", selected.Title))
		b.WriteString(fmt.Sprintf("Author: %s\n", authorStyle.Render(selected.Poster.UserName)))
		b.WriteString(fmt.Sprintf("Status: %s\n", statusStyle.Render("Open")))
		if selected.Updated != nil {
			b.WriteString(fmt.Sprintf("Updated: %s\n", selected.Updated.Format("2006-01-02 15:04")))
		}
		if selected.Body != "" {
			description := selected.Body
			if len(description) > 100 {
				description = description[:97] + "..."
			}
			b.WriteString(fmt.Sprintf("Description: %s\n", description))
		}
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter: checkout PR • v: view details • r: refresh • q/esc: quit"))

	return b.String()
}

func (m ListPRModel) GetResult() ListPRResult {
	if m.selected >= 0 && m.selected < len(m.prs) && (m.action == "checkout" || m.action == "view") {
		return ListPRResult{
			SelectedPR: m.prs[m.selected],
			Action:     m.action,
		}
	}
	return ListPRResult{
		SelectedPR: nil,
		Action:     m.action,
	}
}

func ShowPRList(prs []*gitea.PullRequest, owner, repo string) (*ListPRResult, error) {
	if len(prs) == 0 {
		fmt.Printf("📋 No open pull requests found for %s/%s\n", owner, repo)
		return &ListPRResult{Action: "quit"}, nil
	}

	model := NewListPRModel(prs, owner, repo)
	program := tea.NewProgram(model)
	
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	listModel := finalModel.(ListPRModel)
	result := listModel.GetResult()
	
	return &result, nil
}