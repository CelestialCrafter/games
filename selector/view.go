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
	selectedBar := fmt.Sprintf(
		"\n%v", lipgloss.NewStyle().
			Width(common.ICON_WIDTH).
			Height(1).
			Margin(0, 1).
			Render(strings.Repeat("━", common.ICON_WIDTH)),
	)

	rowAmount := int(math.Ceil(float64(len(m.gamesMetadata)) / float64(m.rowLength)))
	menu := make([][]string, rowAmount)
	for i := range menu {
		menu[i] = make([]string, m.rowLength)
	}

	menuRows := make([]string, rowAmount)

	for i := 0; i < len(m.gamesMetadata); i++ {
		currentRow := int(math.Floor(float64(i) / float64(m.rowLength)))
		currentColumn := i % m.rowLength
		current := &menu[currentRow][currentColumn]

		*current = m.gamesMetadata[i].Icon
		if m.selectedGame == i {
			*current = strings.TrimRightFunc(*current, unicode.IsSpace) + selectedBar
		}
	}

	for i := 0; i < len(menu); i++ {
		menuRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, menu[i]...)
	}

	menuString := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	return fmt.Sprintf("%v\n%v", menuString, m.help.View(m.keys))
}
