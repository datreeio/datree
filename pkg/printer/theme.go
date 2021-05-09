package printer

import (
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

type theme struct {
	Colors struct {
		Green *color.Color
		Yello *color.Color
		Red   *color.Color
		White *color.Color
	}
	Spacing struct {
		Default string
	}
	Emoji struct {
		Error      string
		Suggestion string
	}
}

func createTheme() *theme {
	return &theme{
		Colors: struct {
			Green *color.Color
			Yello *color.Color
			Red   *color.Color
			White *color.Color
		}{
			Green: color.New(color.FgGreen),
			Yello: color.New(color.FgYellow),
			Red:   color.New(color.FgHiRed, color.Bold),
			White: color.New(color.FgHiWhite),
		},
		Spacing: struct{ Default string }{
			Default: strings.Join([]string{" "}, ""),
		},
		Emoji: struct {
			Error      string
			Suggestion string
		}{
			Error:      emoji.Sprint(":cross_mark:"),
			Suggestion: emoji.Sprint(":light_bulb:"),
		},
	}
}
