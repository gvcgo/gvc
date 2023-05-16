package tui

import (
	"fmt"

	"github.com/pterm/pterm"
)

const (
	BarLength int = 192
)

/*
ProgressBar for downloading files.
*/
type ProgressBar struct {
	Bar             *pterm.ProgressbarPrinter
	Filename        string
	ContentLength   int
	CurrentReceived int
	divideBy        int
}

func NewProgressBar(filename string, length int) (p *ProgressBar) {
	if length <= 0 {
		return
	}
	p = &ProgressBar{
		ContentLength: length,
		Filename:      filename,
	}
	p.Bar = pterm.DefaultProgressbar.WithTitle(fmt.Sprintf("Downloading <%s %sMB>", filename, p.calcSize())).WithTotal(BarLength).WithShowCount(false).WithShowPercentage(true)
	p.divideBy = length / BarLength
	return
}

func (that *ProgressBar) calcSize() string {
	mb := float64(that.ContentLength) / float64(1024*1024)
	return fmt.Sprintf("%.2f", mb)
}

func (that *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	that.CurrentReceived += n
	var increasement int
	if that.CurrentReceived == that.ContentLength {
		increasement = that.Bar.Total - that.Bar.Current
	} else {
		increasement = (that.CurrentReceived / that.divideBy) - that.Bar.Current
	}
	if that.CurrentReceived == that.ContentLength {
		that.Bar.UpdateTitle(fmt.Sprintf("Downloading <%s %sMB>", that.Filename, that.calcSize()))
	}
	if increasement > 0 {
		that.Bar.Add(increasement)
	}
	return
}

func (that *ProgressBar) Start() (err error) {
	bar := *that.Bar
	that.Bar, err = bar.Start()
	return
}
