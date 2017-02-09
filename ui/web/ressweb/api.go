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

var (
	jsonFeedID    = "ID"
	jsonFeedTitle = "Title"
	jsonFeedLink  = "Link"
	jsonFeedError = "Error"

	jsonItemID     = jsonFeedID
	jsonItemTitle  = jsonFeedTitle
	jsonItemLink   = jsonFeedLink
	jsonItemUnread = "Unread"
)

// Feed represents a single feed
type Feed struct {
	ID       int
	Title    string
	Link     string
	Selected bool
	Error    string

	Items []*Item
}

// Item represents a single feed item
type Item struct {
	ID       int
	Title    string
	Link     string
	Selected bool
	Unread   bool
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

func (a *api) addFeed(link string) error {
	responseChan := a.r.JSONRequest(
		"POST", a.withBasePath("feeds"),
		json.Stringify(map[string]string{"link": link}),
	)

	response := <-responseChan
	if response.Error != nil {
		return response.Error
	} else if response.Code != 200 {
		return fmt.Errorf(
			"Can't craete a feed %s, server responded with code %d", link, response.Code)
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
	var id float64
	feedJSON, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected feed format")
	}

	feed := &Feed{}
	if id, ok = feedJSON[jsonFeedID].(float64); !ok {
		return nil, fmt.Errorf("feed.%s is %T, not float64", jsonFeedID, feedJSON[jsonFeedID])
	}
	feed.ID = int(id)

	if feed.Title, ok = feedJSON[jsonFeedTitle].(string); !ok {
		return nil, fmt.Errorf("feed.%s is %T, not string", jsonFeedTitle, feedJSON[jsonFeedTitle])
	}

	if feed.Link, ok = feedJSON[jsonFeedLink].(string); !ok {
		return nil, fmt.Errorf("feed.%s is %T, not string", jsonFeedLink, feedJSON[jsonFeedLink])
	}

	if len(feed.Title) == 0 {
		feed.Title = feed.Link
	}

	if feed.Error, ok = feedJSON[jsonFeedError].(string); !ok {
		return nil, fmt.Errorf("feed.%s is %T, not string", jsonFeedError, feedJSON[jsonFeedError])
	}

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
	var id float64
	itemJSON, ok := json.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Unexpected item format")
	}

	item := &Item{}
	if id, ok = itemJSON[jsonItemID].(float64); !ok {
		return nil, fmt.Errorf("item.%s is %T, not float64", jsonItemID, itemJSON[jsonItemID])
	}
	item.ID = int(id)

	if item.Title, ok = itemJSON[jsonItemTitle].(string); !ok {
		return nil, fmt.Errorf("item.%s is %T, not string", jsonItemTitle, itemJSON[jsonItemTitle])
	}

	if item.Link, ok = itemJSON[jsonItemLink].(string); !ok {
		return nil, fmt.Errorf("item.%s is %T, not string", jsonItemLink, itemJSON[jsonItemLink])
	}

	if item.Unread, ok = itemJSON[jsonItemUnread].(bool); !ok {
		return nil, fmt.Errorf("item.%s is %T, not bool", jsonItemUnread, itemJSON[jsonItemUnread])
	}

	return item, nil
}
