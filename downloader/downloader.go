package downloader

import (
	"github.com/mmcdole/gofeed"
	"time"
)

const concurrency = 3

// The Downloader cares about feed download scheduling
type Downloader struct {
	period  time.Duration
	urls    map[string]bool
	handler Handler

	jobsChan chan string
	active   bool
}

// Handler is called with downloaded feed data time after time
type Handler func(url string, feed *gofeed.Feed, err error)

// New initializes Downloader
func New(period time.Duration, h Handler) *Downloader {
	return &Downloader{
		period:  period,
		handler: h,
		urls:    make(map[string]bool)}
}

// Download given the url in future
func (d *Downloader) Download(url string) {
	d.urls[url] = true
}

// Serve starts downloading urls, blocks until Cancel() will not be called
func (d *Downloader) Serve() {
	d.active = true
	d.jobsChan = make(chan string)
	defer close(d.jobsChan)

	for i := 0; i < concurrency; i++ {
		go d.serve()
	}

	for d.active {
		for url := range d.urls {
			d.jobsChan <- url
		}
		time.Sleep(d.period)
	}
}

func (d *Downloader) serve() {
	for {
		if url, ok := <-d.jobsChan; ok {
			d.download(url)
		} else {
			return
		}
	}
}

func (d *Downloader) download(url string) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	d.handler(url, feed, err)
}

// Cancel disables further downloading
func (d *Downloader) Cancel() {
	d.active = false
}
