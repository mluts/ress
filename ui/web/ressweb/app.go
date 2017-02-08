package main

// App is an application state container
type App struct {
	api   *api
	feeds []*Feed
}

func (a *App) selectFeed(id int) {
	for i := range a.feeds {
		a.feeds[i].Selected = a.feeds[i].ID == id
	}
}

func (a *App) subscribeToFeed(f *Feed) {
	a.api.addFeed(f)
}
