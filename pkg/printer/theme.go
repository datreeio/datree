package printer

import (
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

type theme struct {
	Colors struct {
		Warning *color.Color
		Error   *color.Color
		Plain   *color.Color
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
			Warning *color.Color
			Error   *color.Color
			Plain   *color.Color
		}{
			Warning: color.New(color.FgYellow),
			Error:   color.New(color.FgHiRed, color.Bold),
			Plain:   color.New(color.FgHiWhite),
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
