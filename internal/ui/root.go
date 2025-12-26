package ui

type Model struct {
	current Screen
	menu	MenuModel
}

func (m Model) View() string {
	switch m.current {
	case MenuModel:
		return m.menu.View()
	}
	return ""
}
