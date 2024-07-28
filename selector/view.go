package selector

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	greeting := lipgloss.NewStyle().
		BorderForeground(styles.Colors.Accent).
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		Padding(0, 2).
		Render("welcome back!! <3")
	greeting = lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Top, greeting)

	selectedBar := fmt.Sprintf(
		"\n%v", lipgloss.NewStyle().
			Margin(0, 1).
			Foreground(styles.Colors.Secondary).
			Render(strings.Repeat("─", common.ICON_WIDTH+2)),
	)

	rowAmount := int(math.Ceil(float64(len(common.Games)) / float64(m.rowLength)))
	menu := make([][]string, rowAmount)
	for i := range menu {
		menu[i] = make([]string, m.rowLength)
	}

	menuRows := make([]string, rowAmount)

	for i := 0; i < len(common.Games); i++ {
		currentRow := i / m.rowLength
		currentColumn := i % m.rowLength
		current := &menu[currentRow][currentColumn]

		*current = common.Games[uint(i)].Icon
		if m.selectedGame == i {
			*current = strings.TrimRightFunc(*current, unicode.IsSpace) + selectedBar
		}
	}

	for i := 0; i < len(menu); i++ {
		menuRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, menu[i]...)
	}

	menuJoined := lipgloss.JoinVertical(lipgloss.Left, menuRows...)
	help := m.help.View(m.keys)

	availableHeight := m.height
	availableHeight -= lipgloss.Height(greeting)
	availableHeight -= lipgloss.Height(help)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		greeting,
		lipgloss.NewStyle().Height(availableHeight).Render(menuJoined),
		help,
	)
}
