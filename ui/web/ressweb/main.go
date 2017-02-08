package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mluts/ress/ui/web/ressweb/ajax"
	"github.com/mluts/ress/ui/web/ressweb/console"
)

func main() {
	app := &App{}
	app.api = &api{"/api", ajax.New()}

	for i := 0; i < 10; i++ {
		app.feeds = append(app.feeds, &Feed{
			Title: "The Feed 1",
			Link:  "https://example.com/",
			ID:    i})
	}

	ui := js.Global.Get("ui")

	render := func() {
		ui.Call("render", map[string]interface{}{
			"feeds": app.feeds,
		})
	}

	render()

	ui.Call("registerHandler", "onSelectFeed", func(feed *js.Object) {
		app.selectFeed(feed.Get("ID").Int())
		render()
	})

	ui.Call("registerHandler", "onSubscribeToFeed", func(feed *js.Object) {
		console.Log("Subscribing to ", feed)
		go func() {
			app.subscribeToFeed(&Feed{Link: feed.Get("Link").String()})
			console.Log("Done")
		}()
	})
}
