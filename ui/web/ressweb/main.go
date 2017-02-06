package main

import (
	"github.com/gopherjs/gopherjs/js"
	// "github.com/mluts/ress/ui/web/ressweb/ajax"
	// "github.com/mluts/ress/ui/web/ressweb/console"
)

func main() {
	app := &App{}
	for i := 0; i < 10; i++ {
		app.feeds = append(app.feeds, &Feed{
			Title: "The Feed 1",
			Link:  "https://example.com/",
			ID:    i})
	}

	ui := js.Global.Get("ui")
	ui.Call("renderFeeds", app.feeds)
	ui.Call("registerHandler", "onSelectFeed", func(feed *js.Object) {
		app.selectFeed(feed.Get("ID").Int())
		ui.Call("renderFeeds", app.feeds)
	})
}
