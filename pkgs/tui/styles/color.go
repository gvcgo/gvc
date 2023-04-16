package styles

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func NewNoColor() lipgloss.NoColor {
	return lipgloss.NoColor{}
}

func NewColor(i int) lipgloss.Color {
	return lipgloss.Color(strconv.Itoa(i))
}

func NewColorFromHex(hex string) lipgloss.Color {
	return lipgloss.Color(hex)
}

func NewColorAdaptive(light, dark string) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: light,
		Dark:  dark,
	}
}

var (
	Magenta   = NewColor(5)
	Red       = NewColor(9)
	LightBlue = NewColor(12)
	Blank     = NewColor(30)
	Pink      = NewColor(13)
	Green     = NewColor(10)
	Cyan      = NewColor(14)
	Aqua      = NewColor(86)
	HotPink   = NewColor(201)
	Orange    = NewColor(202)
	RedPink   = NewColor(205)

	FullBlue = NewColorFromHex("#0000FF")
	DarkGray = NewColorFromHex("#3C3C3C")
	Gray     = NewColorFromHex("#808080")
	White    = NewColorFromHex("#ffffff")

	Special   = NewColorAdaptive("#43BF6D", "#73F59F")
	Highlight = NewColorAdaptive("#874BFD", "#7D56F4")
	Subtle    = NewColorAdaptive("#D9DCCF", "#383838")
)
