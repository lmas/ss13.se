package ss13_se

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/text/encoding/charmap"
)

const (
	byondURL string = "http://www.byond.com/games/Exadv1/SpaceStation13"
	//byondURL  string = "./tmp/dump.html" // For testing
	userAgent string = "ss13hub/2.0pre"
)

var (
	rePlayers = regexp.MustCompile(`Logged in: (\d+) player`)
	//rePlayers = regexp.MustCompile(`<br/>\s*<br/>\s*Logged in: (\d+) player.*<a href`)
)

func scrapeByond(webClient *http.Client, now time.Time) ([]ServerEntry, error) {
	var body io.ReadCloser
	if byondURL == "./tmp/dump.html" {
		r, err := os.Open(byondURL)
		if err != nil {
			return nil, err
		}
		body = r
	} else {

		r, err := openPage(webClient, byondURL)
		if err != nil {
			return nil, err
		}
		body = r
	}
	defer body.Close()

	servers, err := parseByondPage(now, body)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func openPage(webClient *http.Client, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http.Response.Status: %s", resp.Status)
	}
	return resp.Body, nil
}

func parseByondPage(now time.Time, body io.Reader) ([]ServerEntry, error) {
	// Yep, Byond serves it's pages with Windows-1252 encoding...
	r := charmap.Windows1252.NewDecoder().Reader(body)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var servers []ServerEntry
	doc.Find(".live_game_entry").Each(func(i int, s *goquery.Selection) {
		entry, err := parseEntry(s.Find(".live_game_status"))
		if err != nil {
			log.Println("Error parsing entry:", err)
			return
		}
		if entry.IsZero() {
			return
		}

		// Make sure we only try to add only one instance of a server.
		// And since byond orders the most popular servers up on top,
		// we get a small protection from bad guys who's trying to
		// influent the history of a server.
		for _, s := range servers {
			if s.ID == entry.ID {
				return
			}
		}

		entry.Time = now
		servers = append(servers, entry)
	})

	return servers, nil
}

func parseEntry(s *goquery.Selection) (ServerEntry, error) {
	// Try find a player count (really tricky since it's not in a valid
	// html tag by itself)
	tmp := strings.TrimSpace(strings.Replace(s.Text(), "\n", "", -1))
	r := rePlayers.FindStringSubmatch(tmp)
	// 2 == because the regexp returns wholestring + matched part
	// If it's less than 2 we couldn't find a match and if it's greater
	// than 2 there's multiple matches, which is fishy...
	players := 0
	if len(r) == 2 {
		p, err := strconv.Atoi(r[1])
		if err != nil {
			return ServerEntry{}, err
		}
		players = p
	}

	// Grab and sanitize the server's name
	title := s.Find("b").First().Text()
	title = strings.Replace(strings.TrimSpace(title), "\n", "", -1)
	if len(title) < 1 {
		// the byond page sometimes has server entries that's basiclly
		// blank, no server name or player count (just some byond url)
		return ServerEntry{}, nil
	}
	id := makeID(title)

	gameURL := s.Find("span.smaller").Find("nobr").Text()
	siteURL := s.Find("a").First().AttrOr("href", "")
	if siteURL == "http://" {
		siteURL = ""
	}

	return ServerEntry{
		ID:      id,
		Title:   title,
		SiteURL: siteURL,
		GameURL: gameURL,
		Players: players,
	}, nil
}

func makeID(title string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(title)))
}
