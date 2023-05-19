package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pterm/pterm"
)

type InputItem struct {
	Title   string
	Value   string
	Default string
	Must    bool
}

func (that *InputItem) String() string {
	if that.Value != "" {
		return strings.TrimSpace(that.Value)
	}
	if that.Default != "" {
		fmt.Println(pterm.Cyan(fmt.Sprintf("[%s] ", that.Title)), pterm.Green(fmt.Sprintf("uses default value: %s", that.Default)))
	}
	return that.Default
}

func (that *InputItem) Int() int {
	if that.Value != "" {
		return gconv.Int(that.Value)
	}
	return gconv.Int(that.Default)
}

func (that *InputItem) Bool() bool {
	if that.Value != "" {
		return gconv.Bool(that.Value)
	}
	return gconv.Bool(that.Default)
}

type Input struct {
	Items []*InputItem
}

func NewInput(items []*InputItem) *Input {
	return &Input{Items: items}
}

func (that *Input) Render() {
	iput := pterm.DefaultInteractiveTextInput
	for _, item := range that.Items {
		if !strings.Contains(strings.ToLower(item.Title), "password") {
			item.Value, _ = iput.WithDefaultText(item.Default).Show(item.Title)
		} else {
			item.Value, _ = iput.WithMask("*").WithDefaultText(item.Default).Show(item.Title)
		}
		if item.Must && item.Value == "" && item.Default == "" {
			fmt.Println(pterm.Red(fmt.Sprintf("%s cannot be empty.", item.Title)))
			os.Exit(1)
		}
	}
}
