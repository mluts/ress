package downloader

import (
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"log"
	"time"
)

// The Downloader cares about feed download scheduling
type Downloader struct {
	period   time.Duration
	shutDown bool

	urls    map[string]chan int
	handler Handler

	pool   chan int
	Logger *log.Logger
}

// Handler is called with downloaded feed data time after time
type Handler func(url string, feed *gofeed.Feed, err error)

// New initializes Downloader
func New(period time.Duration, poolSize uint, h Handler) *Downloader {
	return &Downloader{
		period:  period,
		handler: h,
		urls:    make(map[string]chan int),
		pool:    make(chan int, poolSize),
		Logger:  log.New(ioutil.Discard, "", log.LstdFlags)}
}

// Download given the url in future
func (d *Downloader) Download(url string) {
	if d.urls[url] == nil {
		d.Logger.Print("Enqueueing ", url)
		d.urls[url] = make(chan int, 1)
	}
}

// Discard given url from downloading
func (d *Downloader) Discard(url string) {
	d.Logger.Print("Discarding ", url)
	delete(d.urls, url)
}

// Serve starts downloading urls, blocks until Cancel() will not be called
func (d *Downloader) Serve() {
	d.Logger.Print("Serving downloads")
	d.shutDown = false

	for !d.shutDown {
		for url := range d.urls {
			if len(d.urls[url]) > 0 {
				continue
			}

			d.pool <- 1
			d.urls[url] <- 1
			go d.download(url)
		}
		time.Sleep(d.period)
	}
}

func (d *Downloader) download(url string) {
	d.Logger.Print("Downloading ", url)

	defer func() {
		<-d.pool
		<-d.urls[url]

		if err := recover(); err != nil {
			log.Print("Panic on url ", url)
			log.Print(err)
		}
	}()

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	d.handler(url, feed, err)
}

// Cancel disables further downloading
func (d *Downloader) Cancel() {
	d.Logger.Print("Cancelling downloads")
	d.shutDown = true
}
