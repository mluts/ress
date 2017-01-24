package downloader

import (
	"github.com/mmcdole/gofeed"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func fixturesServer() *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir("./fixtures")))
}

var rubyFeedFixture = &gofeed.Feed{
	Title:       "Ruby News",
	Description: "The latest news from ruby-lang.org."}

var rubyItemFixture = &gofeed.Item{
	Title: "Ruby 2.4.0 Released",
	Link:  "https://www.ruby-lang.org/en/news/2016/12/25/ruby-2-4-0-released/"}

func assertSameFeeds(t *testing.T, expected, actual *gofeed.Feed) {
	if expected.Title != actual.Title {
		t.Errorf("Expected to have title %s, but have %s",
			expected.Title, actual.Title)
	}

	if expected.Description != actual.Description {
		t.Errorf("Expected to have description %s, but have %s",
			expected.Description, actual.Description)
	}
}

func assertSameItems(t *testing.T, expected, actual *gofeed.Item) {
	if expected.Title != actual.Title {
		t.Errorf("Expected to have title %s, but have %s",
			expected.Title, actual.Title)
	}

	if expected.Link != actual.Link {
		t.Errorf("Expected to have Link %s, but have %s",
			expected.Link, actual.Link)
	}
}

func testURL(givenURL string, h Handler) {
	var (
		downloader *Downloader
		wg         = &sync.WaitGroup{}
	)

	wg.Add(1)
	downloader = New(time.Nanosecond, 1, func(url string, feed *gofeed.Feed, err error) {
		downloader.Cancel()
		h(url, feed, err)
		wg.Done()
	})
	downloader.Logger.SetOutput(os.Stdout)

	downloader.Enqueue(givenURL)
	go downloader.Serve()
	wg.Wait()
}

func TestDownload_good_url(t *testing.T) {
	i := 0

	server := fixturesServer()
	defer server.Close()

	givenURL := strings.Join([]string{server.URL, "ruby.rss"}, "/")

	testURL(givenURL, func(url string, feed *gofeed.Feed, err error) {
		i++

		if err != nil {
			t.Error("Expected no errors, but have", err)
		}

		if url != givenURL {
			t.Errorf("Expected to receive %s, but have %s", givenURL, url)
		}

		assertSameFeeds(t, rubyFeedFixture, feed)

		if len(feed.Items) != 10 {
			t.Error("Expected to have 10 items, but have", len(feed.Items))
		} else {
			assertSameItems(t, rubyItemFixture, feed.Items[0])
		}
	})

	if i != 1 {
		t.Errorf("Expected to have been server once, but was served %d times", i)
	}
}

func TestDownload_bad_url(t *testing.T) {
	server := fixturesServer()
	defer server.Close()

	givenURL := strings.Join([]string{server.URL, "bad.rss"}, "/")

	testURL(givenURL, func(url string, feed *gofeed.Feed, err error) {
		if err == nil {
			t.Error("Expected some error, but have nil")
		}

		if url != givenURL {
			t.Errorf("Expected to see %s, but having %s", givenURL, url)
		}

		if feed != nil {
			t.Error("Expected feed to be nil")
		}
	})
}

func TestDownload_discard(t *testing.T) {
	server := fixturesServer()
	defer server.Close()

	givenURL := strings.Join([]string{server.URL, "ruby.rss"}, "/")

	var downloader *Downloader

	quit := make(chan int)

	downloader = New(time.Nanosecond, 3, func(url string, feed *gofeed.Feed, err error) {
		downloader.Discard(url)

		if _, present := downloader.urls[url]; present {
			t.Error("Url should be nil after discarding it")
		}

		quit <- 1
	})

	downloader.Enqueue(givenURL)
	go downloader.Serve()
	<-quit
}

func TestDownload_immediate_download(t *testing.T) {
	var (
		downloader *Downloader
		wg         = &sync.WaitGroup{}
		once       = &sync.Once{}
	)
	server := fixturesServer()
	defer server.Close()

	givenURL := strings.Join([]string{server.URL, "ruby.rss"}, "/")

	// Expecting exactly TWO downloads
	wg.Add(2)

	download := func() {
		downloader.Download(givenURL)
	}

	downloader = New(time.Nanosecond, 3, func(url string, feed *gofeed.Feed, err error) {
		downloader.Cancel()
		once.Do(download)
		wg.Done()
	})

	downloader.Enqueue(givenURL)
	go downloader.Serve()
	wg.Wait()
}
