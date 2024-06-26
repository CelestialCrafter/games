package selector

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	greeting := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("6")).
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		Margin(1).
		Padding(0, 2).
		Render("welcome back!! <3")

	selectedBar := fmt.Sprintf(
		"\n%v", lipgloss.NewStyle().
			Width(common.ICON_WIDTH+2).
			Height(1).
			Margin(1).
			Foreground(lipgloss.Color("6")).
			Render(strings.Repeat("━", common.ICON_WIDTH+2)),
	)

	rowAmount := int(math.Ceil(float64(len(common.Games)) / float64(m.rowLength)))
	menu := make([][]string, rowAmount)
	for i := range menu {
		menu[i] = make([]string, m.rowLength)
	}

	menuRows := make([]string, rowAmount)

	for i := 0; i < len(common.Games); i++ {
		currentRow := int(math.Floor(float64(i) / float64(m.rowLength)))
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

	menuString := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	return fmt.Sprintf(
		"%v%v\n%v",
		lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Top, greeting),
		menuString,
		m.help.View(m.keys),
	)
}
