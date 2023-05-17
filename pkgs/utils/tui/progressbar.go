package tui

import (
	"fmt"

	"github.com/pterm/pterm"
)

const (
	BarLength int    = 192
	BarTitle  string = "GetFile|<%s %s/%sMB>"
)

/*
ProgressBar for downloading files.
*/
type ProgressBar struct {
	Bar             *pterm.ProgressbarPrinter
	Filename        string
	ContentLength   int
	Total           string
	CurrentReceived int
	divideBy        int
	writeCount      int64
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
	p.Bar = pterm.DefaultProgressbar.WithTotal(BarLength).WithTitle(fmt.Sprintf(BarTitle, p.Filename, current, total)).WithShowCount(false).WithShowPercentage(true)
	p.divideBy = length / BarLength
	return
}

func (that *ProgressBar) calcSize() (total, current string) {
	if that.Total == "" {
		that.Total = fmt.Sprintf("%.2f", float64(that.ContentLength)/float64(1024*1024))
	}
	total = that.Total
	if that.Total != "" {
		length := len(that.Total) - 3
		mb := float64(that.CurrentReceived) / float64(1024*1024)
		current = fmt.Sprintf("%"+fmt.Sprintf("%v", length)+".2f", mb)
	}
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
	if that.CurrentReceived == that.ContentLength || that.writeCount%10 == 0 {
		total, current := that.calcSize()
		that.Bar.UpdateTitle(fmt.Sprintf(BarTitle, that.Filename, current, total))
	}
	that.writeCount += 1
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
