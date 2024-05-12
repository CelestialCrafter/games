package twenty48

import (
	"fmt"
	"math"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	status := ""

	if m.finished {
		status = styles.StatusStyle.Render("you lose!")
	}

	board := common.RenderBoard(m.board, func(cell uint16) string {
		color := lipgloss.Color(fmt.Sprint(math.Log2(float64(cell))))
		cellString := fmt.Sprint(cell)
		return cellStyle.Background(color).Render(cellString)
	})

	return fmt.Sprintf("%v\n\n%v\n%v", board, status, m.help.View(m.keys))
}
