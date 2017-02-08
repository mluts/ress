package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mluts/ress/ui/web/ressweb/console"
)

// UI wraps corresponding javascript code
type UI struct {
	obj *js.Object
}

func newUI() *UI {
	return &UI{js.Global.Get("ui")}
}

func (ui *UI) render(feeds []*Feed) {
	console.Log("Rendering...", feeds)
	ui.obj.Call("render", map[string]interface{}{
		"feeds": feeds,
	})
}

func (ui *UI) onSelectFeed(fn func(feedID int)) {
	ui.obj.Call("registerHandler", "onSelectFeed", func(feedJS *js.Object) {
		var (
			id float64
			ok bool
		)

		f, ok := feedJS.Interface().(map[string]interface{})
		if !ok {
			console.Err("Wrong argument type for onSelectFeed callback:", feedJS)
			return
		}

		if id, ok = f["ID"].(float64); !ok {
			console.Err("onSelectFeed feed.ID is not a number:", feedJS)
			return
		}

		fn(int(id))
	})
}

func (ui *UI) onSubscribeToFeed(fn func(link string)) {
	ui.obj.Call("registerHandler", "onSubscribeToFeed", func(feedJS *js.Object) {
		var (
			ok   bool
			link string
		)
		f, ok := feedJS.Interface().(map[string]interface{})
		if !ok {
			console.Err("Wrong argument type for onSelectFeed callback")
			return
		}

		link, ok = f["Link"].(string)
		if !ok {
			console.Err("onSubscribeToFeed feed.Link should be a string")
		}

		fn(link)
	})
}
