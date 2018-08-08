package hn

import (
	"fmt"
	"sync"
	"github.com/PuerkitoBio/goquery"
)

const (
	hnBaseURL = "https://news.ycombinator.com"
	hnItemPath = "/item?id="
	maxNews = 30
)

var (
	hnNewsURL = hnBaseURL + "/news"
)

type news struct {
	Title string
	ID string
	ItemID string
	ItemURL string
	SourceURL string
}

func (n *news) parse() {
	n.ItemURL = hnBaseURL + hnItemPath + n.ItemID
}

func (n *news) show() {
	fmt.Printf("-> %s\n", n.Title)
}

func fillList(s *goquery.Selection, newsList *[]news) {
	n := news{
		Title: s.Find(".title:last-of-type .storylink").Text(),
	}

	if n.Title != "" {
		n.ItemID, _ = s.Attr("id")
		n.SourceURL, _ = s.Find(".title:last-of-type .storylink").Attr("href")
		n.parse()		
		
		*newsList = append(*newsList, n)
	}	
}

func GatherNewsFromHackerNews(wg *sync.WaitGroup) {
	var newsList []news
	
	defer wg.Done()

	doc, err := goquery.NewDocument(hnNewsURL)
	
	if err != nil {
		panic(err)
	}
	
	doc.Find(".athing").Each(func(_ int, s *goquery.Selection) {
		if len(newsList) < maxNews {
			fillList(s, &newsList)
		}
	})
	
	fmt.Println("\n\n[HN]\n")
	
	for _, n := range newsList {
		n.show()
	}
}