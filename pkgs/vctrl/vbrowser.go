package vctrl

import (
	"github.com/moqsien/hackbrowser/browser"
	"github.com/moqsien/hackbrowser/item"
	"github.com/moqsien/hackbrowser/log"
)

type Browser struct{}

func Test() {
	log.SetVerbose()
	browsers, err := browser.PickBrowsers("chrome", "")
	if err != nil {
		log.Error(err)
	}

	for _, b := range browsers {
		b.OnlyToSave(item.DefaultOnlyToSave)
		data, err := b.BrowsingData(true)
		if err != nil {
			log.Error(err)
		}
		data.Output("/Users/moqsien/.gvc/backup", b.Name(), "json")
	}
}
