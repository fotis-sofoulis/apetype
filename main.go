package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// -------- Tree Node --------

type MenuNode struct {
	Label    string
	Children []*MenuNode
	Expanded bool
}

// -------- Flattened View Row --------

type flatNode struct {
	Node  *MenuNode
	Depth int
}

// -------- Model --------

type model struct {
	view      string // "menu" or "test"
	testValue string

	cursor int

	nodes []*MenuNode // root menu tree
	flat  []flatNode  // computed visible rows
}

func initialModel() model {
	m := model{
		view: "menu",

		nodes: []*MenuNode{
			{
				Label: "A",
				Children: []*MenuNode{
					{Label: "aa"},
					{Label: "ab"},
					{Label: "ac"},
				},
			},
			{
				Label: "B",
				Children: []*MenuNode{
					{Label: "ba"},
					{Label: "bb"},
					{Label: "bc"},
				},
			},
			{Label: "C"},
		},
	}

	m.recomputeFlat()
	return m
}

// -------- Tree → Visible List --------

func flatten(nodes []*MenuNode, depth int, out *[]flatNode) {
	for _, n := range nodes {
		*out = append(*out, flatNode{Node: n, Depth: depth})
		if n.Expanded {
			flatten(n.Children, depth+1, out)
		}
	}
}

func (m *model) recomputeFlat() {
	m.flat = nil
	flatten(m.nodes, 0, &m.flat)

	if m.cursor >= len(m.flat) {
		 m.cursor = len(m.flat) - 1
	}
	if m.cursor < 0 && len(m.flat) > 0 {
		m.cursor = 0
	}
}

func (m model) Init() tea.Cmd { return nil }

// -------- Update --------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()

		switch m.view {

		// ===== MENU VIEW =====
		case "menu":
			switch key {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "up":
				if m.cursor > 0 {
					m.cursor++
					m.cursor--
				}
				if m.cursor > 0 {
					m.cursor--
				}

			case "down":
				if m.cursor < len(m.flat)-1 {
					m.cursor++
				}

			case "enter":
				if len(m.flat) == 0 {
					return m, nil
				}

				selected := m.flat[m.cursor].Node

				// Parent → toggle expansion
				if len(selected.Children) > 0 {
					selected.Expanded = !selected.Expanded
					m.recomputeFlat()
					return m, nil
				}

				// Leaf → go to test view
				m.testValue = selected.Label
				m.view = "test"
				return m, nil
			}

		// ===== TEST VIEW =====
		case "test":
			switch key {
			case "q", "esc":
				m.view = "menu"
				return m, nil

			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// -------- View --------

func (m model) View() string {
	switch m.view {

	// ----- MENU -----
	case "menu":
		var b strings.Builder
		b.WriteString("Menu (↑/↓ move, Enter select, q quit)\n\n")

		for i, f := range m.flat {

			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}

			indent := strings.Repeat("  ", f.Depth)

			state := ""
			if len(f.Node.Children) > 0 {
				if f.Node.Expanded {
					state = " (-)"
				} else {
					state = " (+)"
				}
			}

			b.WriteString(fmt.Sprintf("%s %s%s%s\n",
				cursor,
				indent,
				f.Node.Label,
				state,
			))
		}

		return b.String()

	// ----- TEST VIEW -----
	case "test":
		return fmt.Sprintf(
			"Test View\n\nSelected: %s\n\nPress q to return\n",
			m.testValue,
		)
	}

	return ""
}

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Println("Error:", err)
	}
}

