package common

const ICON_WIDTH, ICON_HEIGHT = 14, 7

type Metadata struct {
	Name string
	// icons should be 13x7
	Icon string
	// loading, saving, and any other common features the game may support
	Features []string
	ID       uint
}
