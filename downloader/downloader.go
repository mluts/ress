package downloader

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

// Handler is called with downloaded feed data time after time
type Handler func(url string, feed *gofeed.Feed, err error)

// The Downloader cares about feed download scheduling
type Downloader struct {
	Logger *log.Logger

	period time.Duration
	active bool

	urls    map[string]bool
	handler Handler

	poolSize uint

	workChan    chan string
	handlerChan chan *downloadResult
	wg          *sync.WaitGroup
}

type downloadResult struct {
	url  string
	feed *gofeed.Feed
	err  error
}

// New initializes Downloader
func New(period time.Duration, poolSize uint, h Handler) *Downloader {
	return &Downloader{
		period:   period,
		handler:  h,
		urls:     make(map[string]bool),
		poolSize: poolSize,
		Logger:   log.New(ioutil.Discard, "Downloader: ", log.LstdFlags),
		wg:       &sync.WaitGroup{}}
}

// work is a goroutine for doing downloads
func (d *Downloader) work() {
	for {
		if url, open := <-d.workChan; open {
			d.handlerChan <- d.downloadURL(url)
		} else {
			return
		}
	}
}

// handle is a goroutine for calling handler successively
func (d *Downloader) handle() {
	for {
		if result, open := <-d.handlerChan; open {
			func() {
				d.handler(result.url, result.feed, result.err)
				// Be sure to mark this job as done in the case of panic
				defer d.wg.Done()
			}()
		} else {
			return
		}
	}
}

// downloadURL wraps feed downloading and generating some result from it
func (d *Downloader) downloadURL(url string) (result *downloadResult) {
	defer func() {
		if err := recover(); err != nil {
			d.Logger.Printf("Panic during feed download %s", url)
			msg := fmt.Sprintf("%v", err)
			result = &downloadResult{url, nil, errors.New(msg)}
			return
		}
	}()

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	result = &downloadResult{url, feed, err}
	return
}

// Download given url now
func (d *Downloader) Download(url string) {
	d.Logger.Print("Immediate download ", url)
	d.wg.Add(1)
	d.workChan <- url
}

// Enqueue schedules a download for this url
func (d *Downloader) Enqueue(url string) {
	d.Logger.Print("Enqueueing ", url)
	d.urls[url] = true
}

// Discard given url from downloading
func (d *Downloader) Discard(url string) {
	d.Logger.Print("Discarding ", url)
	delete(d.urls, url)
}

// Serve starts downloading urls, blocks until Cancel() will not be called
func (d *Downloader) Serve() {
	d.Logger.Print("Serving downloads")
	d.active = true

	d.handlerChan = make(chan *downloadResult)
	defer close(d.handlerChan)

	d.workChan = make(chan string)
	defer close(d.workChan)

	for i := uint(0); i < d.poolSize; i++ {
		go d.work()
	}
	go d.handle()

	for d.active {
		for url := range d.urls {
			d.wg.Add(1)
			d.workChan <- url
		}
		d.wg.Wait()
		time.Sleep(d.period)
	}
}

// Cancel disables further downloading
func (d *Downloader) Cancel() {
	d.Logger.Print("Cancelling downloads")
	d.active = false
}
