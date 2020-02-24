package loxone_stats

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//
type Miniserver struct {
	Protocol string
	Host     string
	User     MiniserverUser
}

//
type MiniserverUser struct {
	Username string
	Password string
}

//
func NewMiniserver(host string, username string, password string) Miniserver {
	return Miniserver{
		Protocol: "http",
		Host:     host,
		User: MiniserverUser{
			Username: username,
			Password: password,
		},
	}
}

// return a slice of strings containing the uris to a single
// statistics file on the loxone miniserver
func (m Miniserver) FetchStatistics() (stats []Statistic, err error) {
	//
	var url = fmt.Sprintf("%s://%s/stats", m.Protocol, m.Host)
	client := &http.Client{}
	req, err := m.authenticatedRequest("GET", url)

	// perform the request
	res, err := client.Do(req)
	if err != nil {
		return []Statistic{}, err
	}

	//
	links, err := parseStatistics(res.Body)
	if err != nil {
		return []Statistic{}, err
	}

	return links, nil
}

// Return an request with the correct basic
// authentication params set
func (m Miniserver) authenticatedRequest(method string, url string) (req *http.Request, err error) {
	//
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	//
	req.SetBasicAuth(m.User.Username, m.User.Password)

	return req, nil
}

// return a slice of statistic structs
func parseStatistics(r io.Reader) (statsLinks []Statistic, err error) {
	//
	var stats []Statistic
	tokenizer := html.NewTokenizer(r)

	//
	for {
		token := tokenizer.Next()

		switch {

		case token == html.ErrorToken:
			// end of document, we're done
			return stats, nil

		case token == html.StartTagToken:
			//
			t := tokenizer.Token()

			// link found...
			if t.Data == "a" {
				// try getting the href value from the statistic link
				href := parseStatisticsHref(t)

				// check if value exists
				if len(href) > 0 {
					// get the month and year from the statistic link
					month, year := monthYear(href)

					//
					s := Statistic{
						Uri:   href,
						Month: month,
						Year:  year,
					}

					stats = append(stats, s)
				}
			}
		}
	}
}

//
func parseStatisticsHref(token html.Token) string {
	//
	for _, a := range token.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}

// 0d01a765-026e-085a-ffff6f4bfad385ea.201703.xml
func monthYear(link string) (month int, year int) {
	// get positions of the dots which are
	// around the month and year strings
	runes := []rune(link)
	start := strings.Index(link, ".")
	end := strings.LastIndex(link, ".")

	// get the substring containing the year and month
	yearMonth := string(runes[start+1 : end])
	yearMonthRunes := []rune(yearMonth)

	// get the string value of the month and year values
	sYear := string(yearMonthRunes[0:4])
	sMonth := string(yearMonthRunes[4:len(yearMonth)])

	// convert the string values to int's
	year, _ = strconv.Atoi(sYear)
	month, _ = strconv.Atoi(sMonth)

	return month, year
}
