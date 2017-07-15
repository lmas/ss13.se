package ss13_se

import (
	"net/http"
	"time"

	"github.com/SlyMarbo/rss"
)

const redditURL string = "https://www.reddit.com/r/SS13/search.rss?q=ss13.se&restrict_sr=on&t=year&sort=new"

func (a *App) runRedditWatcher(webClient *http.Client) {
	f := func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", userAgent)
		return webClient.Do(req)
	}

	for {
		start := time.Now()
		feed, err := rss.FetchByFunc(f, redditURL)
		dur := time.Since(start)
		a.Log("Updated reddit in %s, errors: %v", dur, err)

		if err == nil {
			a.news = feed.Items
		}
		time.Sleep(a.conf.ScrapeTimeout)
	}
}
