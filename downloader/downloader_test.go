package downloader

import (
	"github.com/mmcdole/gofeed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func fixturesServer() *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir("./fixtures")))
}

var rubyFeedFixture = &gofeed.Feed{
	Title:       "Ruby News",
	Description: "The latest news from ruby-lang.org."}

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

func TestDownload(t *testing.T) {
	var downloader *Downloader
	server := fixturesServer()
	givenURL := strings.Join([]string{server.URL, "ruby.rss"}, "/")

	handler := func(url string, feed *gofeed.Feed, err error) {
		downloader.Cancel()

		if err != nil {
			t.Error("Expected no errors, but have", err)
		}

		if url != givenURL {
			t.Errorf("Expected to receive %s, but have %s", givenURL, url)
		}

		assertSameFeeds(t, rubyFeedFixture, feed)
	}

	downloader = New(time.Microsecond, handler)
	downloader.Download(givenURL)
	downloader.Serve()

	defer server.Close()
}
