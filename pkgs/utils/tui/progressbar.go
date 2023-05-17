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
	total, current := p.calcSize()
	p.Bar = pterm.DefaultProgressbar.WithTitle(fmt.Sprintf("Downloading <%s %s/%sMB>", filename, current, total)).WithTotal(BarLength).WithShowCount(false).WithShowPercentage(true)
	p.divideBy = length / BarLength
	return
}

func (that *ProgressBar) calcSize() (total, current string) {
	mb := float64(that.ContentLength) / float64(1024*1024)
	total = fmt.Sprintf("%.2f", mb)
	length := len(total) - 2
	mb = float64(that.CurrentReceived) / float64(1024*1024)
	current = fmt.Sprintf("%"+fmt.Sprintf("%v", length)+".2f", mb)
	return
}

func (that *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	if that.CurrentReceived == 0 {
		that.Bar.Increment()
	}
	that.CurrentReceived += n
	var increasement int
	if that.CurrentReceived == that.ContentLength {
		increasement = that.Bar.Total - that.Bar.Current
	} else {
		increasement = (that.CurrentReceived / that.divideBy) - that.Bar.Current
	}
	total, current := that.calcSize()
	that.Bar.UpdateTitle(fmt.Sprintf("Downloading <%s %s/%sMB>", that.Filename, current, total))
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
