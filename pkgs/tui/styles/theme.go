package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	PromptStyle                  *Style
	MultiSelectedHintSymbolStyle *Style
	ChoiceTextStyle              *Style
	CursorSymbolStyle            *Style
	UnHintSymbolStyle            *Style
	SpinnerShapeStyle            *Style
	PlaceholderStyle             *Style

	FocusSymbol     string
	UnFocusSymbol   string
	FocusInterval   string
	UnFocusInterval string

	FocusSymbolStyle     *Style
	UnFocusSymbolStyle   *Style
	FocusIntervalStyle   *Style
	UnFocusIntervalStyle *Style
}

var (
	DefaultTheme = Theme{
		PromptStyle:                  NewStyle().Fg(White).Bold(),
		MultiSelectedHintSymbolStyle: NewStyle().Fg(Special),
		ChoiceTextStyle:              NewStyle().Fg(Highlight).Bold(),
		CursorSymbolStyle:            NewStyle(),
		UnHintSymbolStyle:            NewStyle().Fg(Red),
		SpinnerShapeStyle:            NewStyle().Fg(RedPink),
		PlaceholderStyle:             NewStyle().Fg(lipgloss.Color("240")),
		FocusSymbol:                  "? ",
		UnFocusSymbol:                "√ ",
		FocusInterval:                " » ",
		UnFocusInterval:              " … ",
		FocusSymbolStyle:             NewStyle().Fg(Cyan),
		UnFocusSymbolStyle:           NewStyle().Fg(Green),
		FocusIntervalStyle:           NewStyle().Fg(Gray),
		UnFocusIntervalStyle:         NewStyle().Fg(Gray).Bold(),
	}

	// fix https://github.com/fzdwx/infinite/issues/5
	_ = DefaultTheme.PromptStyle.Render("123")
	_ = DefaultTheme.MultiSelectedHintSymbolStyle.Render("123")
	_ = DefaultTheme.ChoiceTextStyle.Render("123")
	_ = DefaultTheme.CursorSymbolStyle.Render("123")
	_ = DefaultTheme.UnHintSymbolStyle.Render("123")
	_ = DefaultTheme.SpinnerShapeStyle.Render("123")
	_ = DefaultTheme.PlaceholderStyle.Render("123")
)
