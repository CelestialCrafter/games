package common

import "flag"

var (
	EnableSsh = flag.Bool("ssh", false, "turns into a ssh server")
	// https://github.com/muesli/termenv/blob/51d72d34e2b9778a31aa5dd79fbdd8cdac50b4d5/profile.go#L12
	ColorProfile      = flag.Int("profile", 1, "chooses a color profile")
	EnableMultiplayer = flag.Bool("multiplayer", true, "enables multiplayer mode")
)
