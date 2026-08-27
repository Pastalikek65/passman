package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Pastalikek65/passman/store"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Padding(0,1)
	selStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	normStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	faintStyle = lipgloss.NewStyle().Faint(true)
)

type Model struct {
	store *store.Store
	entries []store.Entry
	cursor int
}

func New(s *store.Store) Model {
	ents, _ := s.List()
	return Model{store: s, entries: ents}
}

func (m Model) Init() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q","ctrl+c": return m, tea.Quit
		case "j","down": if m.cursor < len(m.entries)-1 { m.cursor++ }
		case "k","up": if m.cursor>0 { m.cursor-- }
		}
	}
	return m, nil
}
func (m Model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(" passman — Termux ") + "\n")
	if len(m.entries)==0 {
		b.WriteString("  no entries — passman add <service> <username>\n")
	}
	for i,e := range m.entries {
		line := fmt.Sprintf(" %-20s %-15s", e.Service, e.Username)
		if i==m.cursor { b.WriteString(selStyle.Render("▶"+line)+"\n") } else { b.WriteString(normStyle.Render(" "+line)+"\n") }
	}
	b.WriteString("\n"+faintStyle.Render("j/k:nav  q:quit  (copy with passman get <service>)")+"\n")
	return b.String()
}
