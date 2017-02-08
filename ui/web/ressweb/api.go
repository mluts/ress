package main

import (
	"fmt"
	"github.com/mluts/ress/ui/web/ressweb/ajax"
	"github.com/mluts/ress/ui/web/ressweb/json"
	"strconv"
	"strings"
)

type jsonRequester interface {
	JSONRequest(method, url string, data ...interface{}) chan *ajax.Response
}

type api struct {
	basePath string
	r        jsonRequester
}

// Feed represents a single feed
type Feed struct {
	ID       int
	Title    string
	Link     string
	Selected bool
}

// Item represents a single feed item
type Item struct {
	ID    int
	Title string
	Link  string
}

func (a *api) getFeeds() ([]*Feed, error) {
	responseChan := a.r.JSONRequest("GET", a.withBasePath("feeds"))
	response := <-responseChan
	if response.Error != nil {
		return nil, response.Error
	} else if response.Code != 200 {
		return nil, fmt.Errorf(
			"Can't get feeds, server returned response code %d", response.Code)
	}
	return parseFeeds(response.Body)
}

func (a *api) getFeed(id int) (*Feed, error) {
	responseChan := a.r.JSONRequest("GET", a.withBasePath(
		"feeds", strconv.Itoa(id)))

	response := <-responseChan
	if response.Error != nil {
		return nil, response.Error
	} else if response.Code != 200 {
		return nil, fmt.Errorf(
			"Can't get feed %d, server responded with code %d", id, response.Code)
	}

	return parseFeed(response.Body)
}

func (a *api) addFeed(f *Feed) error {
	responseChan := a.r.JSONRequest(
		"POST", a.withBasePath("feeds"),
		json.Stringify(map[string]string{"link": f.Link}),
	)

	response := <-responseChan
	if response.Error != nil {
		return response.Error
	} else if response.Code != 200 {
		return fmt.Errorf(
			"Can't craete a feed %s, server responded with code %d", f.Link, response.Code)
	}

	return nil
}

func (a *api) getItems(feedID int) ([]*Item, error) {
	responseChan := a.r.JSONRequest("GET", a.withBasePath(
		"feeds", strconv.Itoa(feedID), "items"))
	response := <-responseChan

	if response.Error != nil {
		return nil, response.Error
	} else if response.Code != 200 {
		return nil, fmt.Errorf(
			"Can't get items for feed %d, server responded with %d", feedID, response.Code)
	}

	return parseItems(response.Body)
}

func (a *api) withBasePath(path ...string) string {
	out := []string{a.basePath}
	out = append(out, path...)
	return strings.Join(out, "/")
}

func parseFeeds(json interface{}) ([]*Feed, error) {
	feeds := make([]*Feed, 0)

	collection, ok := json.([]interface{})
	if !ok {
		return feeds, fmt.Errorf("Bad feeds response")
	}

	for i := range collection {
		feed, err := parseFeed(collection[i])
		if err != nil {
			return make([]*Feed, 0), fmt.Errorf(
				"Can't parse single feed at %d: %v", i, err)
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func parseFeed(json interface{}) (*Feed, error) {
	feedJSON, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected feed format")
	}

	feed := &Feed{}
	feed.ID = feedJSON["id"].(int)
	feed.Title = feedJSON["title"].(string)
	feed.Link = feedJSON["link"].(string)

	return feed, nil
}

func parseItems(json interface{}) ([]*Item, error) {
	items := make([]*Item, 0)

	collection, ok := json.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Bad items response")
	}

	for i := range collection {
		item, err := parseItem(collection[i])
		if err != nil {
			return make([]*Item, 0), err
		}
		items = append(items, item)
	}

	return items, nil
}

func parseItem(json interface{}) (*Item, error) {
	itemJSON, ok := json.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Unexpected item format")
	}

	item := &Item{}
	item.ID = itemJSON["id"].(int)
	item.Link = itemJSON["link"].(string)
	item.Title = itemJSON["title"].(string)

	return item, nil
}
