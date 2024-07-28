package styles

import "github.com/charmbracelet/lipgloss"

type AppColors struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Muted     lipgloss.Color
}

var Colors = AppColors{
	Primary:   lipgloss.Color("213"),
	Secondary: lipgloss.Color("211"),
	Accent:    lipgloss.Color("204"),
	Muted:     lipgloss.Color("244"),
}

// https://github.com/fidian/ansi/blob/master/ansi
// ansi --color-codes
var CellColors = []lipgloss.Color{
	lipgloss.Color("214"),
	lipgloss.Color("215"),
	lipgloss.Color("216"),
	lipgloss.Color("217"),
	lipgloss.Color("182"),
	lipgloss.Color("183"),
	lipgloss.Color("218"),
	lipgloss.Color("219"),

	lipgloss.Color("208"),
	lipgloss.Color("209"),
	lipgloss.Color("210"),
	lipgloss.Color("211"),
	lipgloss.Color("176"),
	lipgloss.Color("177"),
	lipgloss.Color("212"),
	lipgloss.Color("213"),
}
