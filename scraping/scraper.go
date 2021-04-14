package main

import (
	"database/sql"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/html"
)

type apartment struct {
	ObjNr          string         `db:"obj_nr"`
	Hood           string         `db:"hood"`
	AptType        string         `db:"type"`
	Address        string         `db:"address"`
	AptNr          string         `db:"apt_nr"`
	AvailableUntil time.Time      `db:"available_until"`
	BestPoints     int            `db:"best_points"`
	Bookers        int            `db:"bookers"`
	InfoLink       string         `db:"info_link"`
	FloorPlanLink  string         `db:"floor_plan_link"`
	PlanLink       string         `db:"plan_link"`
	MoveIn         time.Time      `db:"move_in"`
	Rent           int            `db:"rent"`
	Sqm            int            `db:"sqm"`
	Special        sql.NullString `db:"special"`
}

func fullScrape() {
	aptsListLink := "https://sssb.se/widgets/?callback=a&widgets%5B%5D=objektlista%40lagenheter"

	resp, err := http.Get(aptsListLink)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	var aptRefIDs []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					re := regexp.MustCompile(`(?:refid=)(.*)(?:\\)`)
					link := a.Val
					refID := re.FindString(link)
					aptRefIDs = append(aptRefIDs, refID[6:len(refID)-1])
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func main() {
	fullScrape()
}
