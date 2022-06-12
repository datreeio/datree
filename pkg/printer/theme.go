package printer

import (
	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

type Theme struct {
	Name   string
	Colors struct {
		Green     *color.Color
		Yellow    *color.Color
		RedBold   *color.Color
		White     *color.Color
		WhiteBold *color.Color
		Error     *color.Color
		Cyan      *color.Color
		CyanBold  *color.Color
	}
	ColorsAttributes struct {
		Cyan  color.Attribute
		Green color.Attribute
		Red   color.Attribute
	}
	Spacing struct {
		Default string
	}
	Emoji struct {
		Error      string
		Suggestion string
		Skip       string
	}
}

func createDefaultTheme() *Theme {
	return &Theme{
		Name: "Default",
		Colors: struct {
			Green     *color.Color
			Yellow    *color.Color
			RedBold   *color.Color
			White     *color.Color
			WhiteBold *color.Color
			Error     *color.Color
			Cyan      *color.Color
			CyanBold  *color.Color
		}{
			Green:     color.New(color.FgGreen),
			Yellow:    color.New(color.FgYellow),
			RedBold:   color.New(color.FgHiRed, color.Bold),
			Error:     color.New(color.FgHiRed),
			White:     color.New(color.FgHiWhite),
			WhiteBold: color.New(color.FgHiWhite, color.Bold),
			Cyan:      color.New(color.FgCyan),
			CyanBold:  color.New(color.FgCyan, color.Bold),
		},
		ColorsAttributes: struct {
			Cyan  color.Attribute
			Green color.Attribute
			Red   color.Attribute
		}{
			Cyan:  color.FgCyan,
			Green: color.FgGreen,
			Red:   color.FgHiRed,
		},
		Spacing: struct{ Default string }{
			Default: strings.Join([]string{" "}, ""),
		},
		Emoji: struct {
			Error      string
			Suggestion string
			Skip       string
		}{
			Error:      emoji.Sprint(":cross_mark:"),
			Suggestion: emoji.Sprint(":light_bulb:"),
			Skip:       emoji.Sprint(":fast_forward:"),
		},
	}
}
func CreateSimpleTheme() *Theme {
	return &Theme{
		Name: "Simple",
		Colors: struct {
			Green     *color.Color
			Yellow    *color.Color
			RedBold   *color.Color
			White     *color.Color
			WhiteBold *color.Color
			Error     *color.Color
			Cyan      *color.Color
			CyanBold  *color.Color
		}{
			Green:     color.New(),
			Yellow:    color.New(),
			RedBold:   color.New(),
			Error:     color.New(),
			White:     color.New(),
			WhiteBold: color.New(),
			Cyan:      color.New(),
			CyanBold:  color.New(),
		},
		ColorsAttributes: struct {
			Cyan  color.Attribute
			Green color.Attribute
			Red   color.Attribute
		}{
			Cyan:  color.Reset,
			Green: color.Reset,
			Red:   color.Reset,
		},
		Spacing: struct{ Default string }{
			Default: strings.Join([]string{" "}, ""),
		},
		Emoji: struct {
			Error      string
			Suggestion string
			Skip       string
		}{
			Error:      "[X] ",
			Suggestion: "[*] ",
			Skip:       "[>>]",
		},
	}
}
