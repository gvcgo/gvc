package tui

import (
	"strings"

	"github.com/pterm/pterm"
)

type FadeColors struct {
	content any
}

func NewFadeColors(content any) *FadeColors {
	return &FadeColors{
		content: content,
	}
}

func (that *FadeColors) Println() {
	from := pterm.NewRGB(0, 255, 255)
	to := pterm.NewRGB(255, 0, 255)
	var fadeInfo string
	if res, ok := that.content.(string); ok {
		strs := strings.Split(res, "")
		for i := 0; i < len(res); i++ {
			fadeInfo += from.Fade(0, float32(len(res)), float32(i), to).Sprint(strs[i])
		}
	} else if res, ok := that.content.([]string); ok {
		content := strings.Join(res, "  ")
		strs := strings.Split(content, "")
		for i := 0; i < len(content); i++ {
			fadeInfo += from.Fade(0, float32(len(content)), float32(i), to).Sprint(strs[i])
		}
	}
	pterm.Println(fadeInfo)
}
