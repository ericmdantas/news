package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type news struct {
	Group string
	URL   string
	Info  []struct {
		Index int
		Title string
		URL   string
	}
}

func (ns *news) addGroup(group, url string) {
	ns.Group = group
	ns.URL = url
}

func (ns *news) addNews(title, url string) {
	ns.Info = append(ns.Info, struct {
		Index int
		Title string
		URL   string
	}{len(ns.Info), title, url})
}

var newsDB = []news{}
var start = time.Now()

func hn(wg *sync.WaitGroup) {
	const hnName = "HN"
	const hnURL = "https://news.ycombinator.com/news"

	defer wg.Done()

	doc, err := goquery.NewDocument(hnURL)

	if err != nil {
		panic(err)
	}

	ns := news{
		Group: hnName,
		URL:   hnURL,
	}

	doc.Find(".athing").Each(func(_ int, s *goquery.Selection) {
		txt := s.Find(".title:last-of-type .storylink").Text()
		url := s.Find(".title:last-of-type .storylink").Text()
		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func runAll(wg *sync.WaitGroup) {
	wg.Add(1)
	go hn(wg)
}

func logStats() {
	fmt.Printf("\nElapsed time: %v", time.Since(start))
}

func logNews() {
	for _, news := range newsDB {
		fmt.Printf("\n%s (%s)\n\n", news.Group, news.URL)

		for _, newsInfo := range news.Info {
			fmt.Printf("%d. %s (%s)\n", newsInfo.Index, newsInfo.Title, newsInfo.URL)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	runAll(&wg)
	wg.Wait()
	logNews()
	logStats()
}
