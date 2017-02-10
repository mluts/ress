package main

import (
	"github.com/mluts/ress/ui/web/ressweb/console"
	"time"
)

const refreshRate = time.Second / 2
const downloadTimeout = time.Second * 30

func main() {
	throttle := time.Tick(refreshRate)
	downloadTick := time.Tick(downloadTimeout)

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
			console.Log("UI was updated")
			<-throttle
			ui.render(app.feeds)
		}
	}()

	go func() {
		for {
			<-downloadTick
			console.Log("Downloading feeds...")
			app.downloadFeeds()
		}
	}()

	app.downloadFeeds()
	ui.render(app.feeds)
}
