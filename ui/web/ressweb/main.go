package main

import (
	"github.com/mluts/ress/ui/web/ressweb/console"
	"time"
)

const rate = time.Second / 2

func main() {
	throttle := time.Tick(rate)

	app := newApp()

	ui := newUI()

	ui.onSelectFeed(func(feedID int) {
		console.Log("Selected feed:", feedID)
		go app.selectFeed(feedID)
	})

	ui.onSubscribeToFeed(func(link string) {
		console.Log("Subscribing to:", link)
		go func() {
			app.subscribeToFeed(link)
			console.Log("Subscribed, downloading feeds...")
			app.downloadFeeds()
			console.Log("Done")
		}()
	})

	ui.onSelectItem(func(id int) {
		console.Log("Selecting an item:", id)
		go app.selectItem(id)
	})

	go func() {
		for {
			<-app.update
			console.Log("Having an update")
			<-throttle
			ui.render(app.feeds)
		}
	}()

	app.downloadFeeds()
	ui.render(app.feeds)
}
