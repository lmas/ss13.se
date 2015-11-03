package ss13

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	RE_PLAYERS = regexp.MustCompile(`Logged in: (\d+) player`)
)

func ScrapePage() []*RawServerData {
	data := download_data()
	return parse_data(data)
}

func download_data() *goquery.Document {
	var (
		doc *goquery.Document
		err error
	)
	if IsDebugging() {
		fmt.Println("Scraper data source: ./dump.html")
		f, err := os.Open("./tmp/dump.html")
		check_error(err)
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		check_error(err)
	} else {
		doc, err = goquery.NewDocument("http://www.byond.com/games/exadv1/spacestation13")
		check_error(err)
	}
	return doc
}

func parse_data(data *goquery.Document) []*RawServerData {
	var servers []*RawServerData
	data.Find(".live_game_entry").Each(func(i int, s *goquery.Selection) {
		tmp := parse_server_data(s)
		if tmp != nil {
			servers = append(servers, tmp)
		}
	})
	return servers
}

func parse_server_data(raw *goquery.Selection) *RawServerData {
	s := raw.Find(".live_game_status")

	t := s.Find("b").First()
	if t.Find("b").Length() > 0 {
		t = t.Find("b").First()
	}
	title := strings.TrimSpace(t.Text())
	title = strings.Replace(title, "\n", "", -1)
	if len(title) < 1 {
		// Yes, someone has made a public server without a server name at least once
		return nil
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
		check_error(err)
		players = int(p)
	}

	return &RawServerData{title, game_url, site_url, players, Now()}
}
