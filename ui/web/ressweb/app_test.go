package main

import "testing"

func prepareApp() *App {
	app := &App{
		&api{}, []*Feed{},
	}

	app.feeds = append(app.feeds,
		&Feed{ID: 1, Title: "The Title 1",
			Link: "https://example1.com/"},
		&Feed{ID: 2, Title: "The Title 2",
			Link: "https://example2.com/"},
		&Feed{ID: 3, Title: "The Title 3",
			Link: "https://example3.com/"},
	)

	return app
}

func assertSelected(t *testing.T, f *Feed) {
	if !f.Selected {
		t.Errorf("Feed %d should be selected", f.ID)
	}
}

func assertNotSelected(t *testing.T, f *Feed) {
	if f.Selected {
		t.Errorf("Feed %d should not be selected", f.ID)
	}
}

func TestApp_select_feed(t *testing.T) {
	app := prepareApp()
	app.selectFeed(1)

	assertSelected(t, app.feeds[0])
	assertNotSelected(t, app.feeds[1])
	assertNotSelected(t, app.feeds[2])

	app.selectFeed(3)

	assertSelected(t, app.feeds[2])
	assertNotSelected(t, app.feeds[0])
	assertNotSelected(t, app.feeds[1])
}
