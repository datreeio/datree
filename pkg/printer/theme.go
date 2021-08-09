package printer

import (
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

type Theme struct {
	Colors struct {
		Green   *color.Color
		Yellow  *color.Color
		RedBold *color.Color
		White   *color.Color
		Error   *color.Color
	}
	ColorsAttributes struct {
		Green color.Attribute
		Red   color.Attribute
	}
	Spacing struct {
		Default string
	}
	Emoji struct {
		Error      string
		Suggestion string
	}
}

func createDefaultTheme() *Theme {
	return &Theme{
		Colors: struct {
			Green   *color.Color
			Yellow  *color.Color
			RedBold *color.Color
			White   *color.Color
			Error   *color.Color
		}{
			Green:   color.New(color.FgGreen),
			Yellow:  color.New(color.FgYellow),
			RedBold: color.New(color.FgHiRed, color.Bold),
			Error:   color.New(color.FgHiRed),
			White:   color.New(color.FgHiWhite),
		},
		ColorsAttributes: struct {
			Green color.Attribute
			Red   color.Attribute
		}{
			Green: color.FgGreen,
			Red:   color.FgRed,
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
func CreateSimpleTheme() *Theme {
	return &Theme{
		Colors: struct {
			Green   *color.Color
			Yellow  *color.Color
			RedBold *color.Color
			White   *color.Color
			Error   *color.Color
		}{
			Green:   color.New(),
			Yellow:  color.New(),
			RedBold: color.New(),
			Error:   color.New(),
			White:   color.New(),
		},
		ColorsAttributes: struct {
			Green color.Attribute
			Red   color.Attribute
		}{
			Green: color.Reset,
			Red:   color.Reset,
		},
		Spacing: struct{ Default string }{
			Default: strings.Join([]string{" "}, ""),
		},
		Emoji: struct {
			Error      string
			Suggestion string
		}{
			Error:      "[X] ",
			Suggestion: "[*] ",
		},
	}
}
