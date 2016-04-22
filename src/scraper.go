package ss13

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/PuerkitoBio/goquery"
)

var (
	RE_PLAYERS = regexp.MustCompile(`Logged in: (\d+) player`)
)

func (i *Instance) ScrapePage() ([]*RawServerData, error) {
	data, err := download_data(i.Debug)
	if err != nil {
		return nil, err
	}

	tmp, err := parse_data(data)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func download_data(debug bool) (*goquery.Document, error) {
	var r io.Reader
	if debug {
		f, err := os.Open("./tmp/dump.html")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = charmap.Windows1252.NewDecoder().Reader(f)
	} else {
		client := &http.Client{
			Timeout: time.Duration(1) * time.Minute,
		}
		resp, err := client.Get("http://www.byond.com/games/exadv1/spacestation13")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		// Yep, Byond serve's it's pages with Windows-1252 encoding...
		r = charmap.Windows1252.NewDecoder().Reader(resp.Body)
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func parse_data(data *goquery.Document) ([]*RawServerData, error) {
	var servers []*RawServerData
	data.Find(".live_game_entry").Each(func(i int, s *goquery.Selection) {
		tmp, err := parse_server_data(s)
		if !LogError(err) {
			if tmp != nil {
				servers = append(servers, tmp)
			}
		}
	})
	return servers, nil
}

func parse_server_data(raw *goquery.Selection) (*RawServerData, error) {
	s := raw.Find(".live_game_status")

	t := s.Find("b").First()
	if t.Find("b").Length() > 0 {
		t = t.Find("b").First()
	}
	title := strings.TrimSpace(t.Text())
	//title = toUtf8([]byte(title))
	title = strings.Replace(title, "\n", "", -1)
	if len(title) < 1 {
		// Yes, someone has made a public server without a server name at least once
		return nil, fmt.Errorf("Empty name for server")
	}

	game_url := s.Find("span.smaller").Find("nobr").Text()

	site_url := s.Find("a").First().AttrOr("href", "")
	if site_url == "http://" {
		site_url = ""
	}

	players := 0
	tmp := strings.Replace(raw.Find("div").Text(), "\n", "", -1)
	ret := RE_PLAYERS.FindStringSubmatch(tmp)
	// 2 = because the regexp returns wholestring + matched part
	// If it's less than 2 we couldn't find a match and if it's greater
	// than 2 there's multiple matches, which is fishy...
	if len(ret) == 2 {
		p, err := strconv.ParseInt(ret[1], 10, 0)
		if err != nil {
			return nil, err
		}
		players = int(p)
	}

	return &RawServerData{title, game_url, site_url, players, Now()}, nil
}
