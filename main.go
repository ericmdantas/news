package main

import (
	"flag"
	"fmt"
	"log"
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

func createNews(group, url string) news {
	return news{
		Group: group,
		URL:   url,
	}
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
	const name = "HN"
	const url = "https://news.ycombinator.com/news"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		panic(err)
	}

	ns := createNews(name, url)

	doc.Find(".athing").Each(func(_ int, s *goquery.Selection) {
		txt := s.Find(".title:last-of-type .storylink").Text()
		url := ""
		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func physOrg(wg *sync.WaitGroup) {
	const name = "phys.org"
	const url = "https://phys.org"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		panic(err)
	}

	ns := createNews(name, url)

	doc.Find(".news-box h3").Each(func(_ int, s *goquery.Selection) {
		txt := s.Text()
		url := ""
		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func lookupAtSpace(wg *sync.WaitGroup) {
	const name = "lookupat.space"
	const url = "http://lookupat.space/fetch/5b242458c632ddb37d2e00b2f30dca144b8de39d08d3b20fdadbfcff0fc8a25f"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		panic(err)
	}

	ns := createNews(name, url)

	doc.Find(".post").Each(func(_ int, s *goquery.Selection) {
		txt := s.Find(".post-title").Eq(0).Text()
		url, _ := s.Find(".btn.btn-default.border").Eq(0).Attr("href")

		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func run(t string, wg *sync.WaitGroup) {
	if t == "" {
		log.Fatal("choose the type")
	}

	switch t {
	case "hn":
		wg.Add(1)
		go hn(wg)
		break
	case "science":
		wg.Add(1)
		go physOrg(wg)
		go lookupAtSpace(wg)
		break
	case "all":
		wg.Add(3)
		go hn(wg)
		go physOrg(wg)
		go lookupAtSpace(wg)
		break
	default:
		log.Fatalf("%s: type doesn't exist", t)
	}
}

func logStats() {
	fmt.Printf("\nElapsed time: %v", time.Since(start))
}

func logNews() {
	for _, news := range newsDB {
		fmt.Printf("\n%s - %s\n\n", news.Group, news.URL)

		for _, newsInfo := range news.Info {
			fmt.Printf("%d. %s", newsInfo.Index, newsInfo.Title)

			if newsInfo.URL != "" {
				fmt.Printf(" - %s", newsInfo.URL)
			}

			fmt.Println()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	t := flag.String("t", "", "news type")
	flag.Parse()

	run(*t, &wg)
	wg.Wait()
	logNews()
	logStats()
}
