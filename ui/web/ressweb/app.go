package main

import (
	"github.com/mluts/ress/ui/web/ressweb/ajax"
	"github.com/mluts/ress/ui/web/ressweb/console"
)

// App is an application state container
type App struct {
	api                            *api
	feeds                          []*Feed
	feedsMutex                     chan int
	update                         chan int
	selectedFeedID, selectedItemID int
}

func newApp() *App {
	return &App{
		api:        &api{"/api", ajax.New()},
		feeds:      make([]*Feed, 0),
		feedsMutex: make(chan int, 1),
		update:     make(chan int),
	}
}

func (a *App) selectFeed(id int) {
	a.feedsMutex <- 1
	defer func() { <-a.feedsMutex }()

	a.selectedFeedID = id
	for _, feed := range a.feeds {
		feed.Selected = feed.ID == id
	}
	a.update <- 1
}

func (a *App) subscribeToFeed(link string) {
	a.api.addFeed(link)
}

func (a *App) downloadFeeds() {
	a.feedsMutex <- 1
	defer func() { <-a.feedsMutex }()

	feeds, err := a.api.getFeeds()
	if err != nil {
		console.Err("Can't download feeds:", err.Error())
		return
	}

	for i := range feeds {
		a.downloadItems(feeds[i])
		if feeds[i].ID == a.selectedFeedID {
			feeds[i].Selected = true
		}
	}

	a.feeds = feeds
	a.update <- 1
}

func (a *App) downloadItems(feed *Feed) {
	items, err := a.api.getItems(feed.ID)
	if err != nil {
		console.Err("Can't download feed items:", err.Error())
		return
	}

	for i := range items {
		if items[i].ID == a.selectedItemID {
			items[i].Selected = true
		}
	}

	feed.Items = items
}

func (a *App) selectItem(id int) {
	a.feedsMutex <- 1
	defer func() { <-a.feedsMutex }()

	a.selectedItemID = id
	for _, feed := range a.feeds {
		if feed.Selected {
			for _, item := range feed.Items {
				item.Selected = item.ID == id
			}
		}
	}
	a.update <- 1
}
