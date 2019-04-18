package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
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

type reddit struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (r reddit) isEmpty() bool {
	return len(r.Data.Children) == 0
}

func (r reddit) retrieveNews() []simpleInfo {
	var sList = []simpleInfo{}

	for _, v := range r.Data.Children {
		sList = append(sList, simpleInfo{
			Title: v.Data.Title,
			URL:   v.Data.URL,
		})
	}

	return sList
}

type simpleInfo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type fetcher struct {
	token     string
	fetchList []func()
}

var newsDB = []news{}
var fetchersDB = []fetcher{}
var start = time.Now()

func registerFetchers(wg *sync.WaitGroup) {
	fetchers := []fetcher{
		{
			token: "hn",
			fetchList: []func(){
				func() {
					hn(wg)
				},
			},
		},
		{
			token: "science",
			fetchList: []func(){
				func() {
					lookupAtSpace(wg)
				},
				func() {
					physOrg(wg)
				},
			},
		},
		{
			token: "br",
			fetchList: []func(){
				func() {
					g1(wg)
				},
			},
		},
		{
			token: "cs",
			fetchList: []func(){
				func() {
					cs(wg)
				},
			},
		},
		{
			token: "motivation",
			fetchList: []func(){
				func() {
					motivation(wg)
				},
			},
		},
		{
			token: "mma",
			fetchList: []func(){
				func() {
					mmafighting(wg)
				},
			},
		},
	}

	fetchAll := fetcher{
		token: "all",
	}

	var allFetchList []func()

	for _, v := range fetchers {
		for _, cb := range v.fetchList {
			allFetchList = append(allFetchList, cb)
		}
	}

	fetchAll.fetchList = allFetchList

	fetchersDB = append(fetchers, fetchAll)
}

func mmafighting(wg *sync.WaitGroup) {
	const name = "MMAFighting"
	const url = "https://mmafighting.com/news"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		panic(err)
	}

	ns := createNews(name, url)

	doc.Find(".c-entry-box--compact__title a").Each(func(_ int, s *goquery.Selection) {
		txt := strings.TrimSpace(s.Text())
		url, _ := s.Attr("href")
		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

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
		url, _ := s.Find(".title:last-of-type .storylink").Attr("href")
		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func g1(wg *sync.WaitGroup) {
	const name = "g1"
	const url = "https://g1.globo.com/"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
	}

	ns := createNews(name, url)

	doc.Find(".feed-post-body .feed-post-body-title a").Each(func(_ int, s *goquery.Selection) {
		txt := s.Text()
		url, _ := s.Attr("href")

		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func physOrg(wg *sync.WaitGroup) {
	const name = "phys.org"
	const url = "https://phys.org/latest-news/"

	defer wg.Done()

	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
	}

	ns := createNews(name, url)

	doc.Find(".sorted-article-content .news-link").Each(func(_ int, s *goquery.Selection) {
		txt := s.Text()
		url, _ := s.Attr("href")
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
		log.Fatal(err)
	}

	ns := createNews(name, url)

	doc.Find(".post").Each(func(_ int, s *goquery.Selection) {
		txt := s.Find(".post-title").Eq(0).Text()
		url, _ := s.Find(".btn.btn-default.border").Eq(0).Attr("href")

		ns.addNews(txt, url)
	})

	newsDB = append(newsDB, ns)
}

func cs(wg *sync.WaitGroup) {
	const name = "Global Offensive"
	const url = "http://reddit.com/r/GlobalOffensive/.json"

	var csStuff reddit

	defer wg.Done()

	ns := createNews(name, url)

	for csStuff.isEmpty() {
		time.Sleep(1500 * time.Millisecond)
		r, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		err = json.NewDecoder(r.Body).Decode(&csStuff)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range csStuff.retrieveNews() {
		ns.addNews(v.Title, v.URL)
	}

	newsDB = append(newsDB, ns)
}

func motivation(wg *sync.WaitGroup) {
	const name = "Get Motivated"
	const url = "http://reddit.com/r/GetMotivated/.json"

	var csStuff reddit

	defer wg.Done()

	ns := createNews(name, url)

	for csStuff.isEmpty() {
		time.Sleep(1500 * time.Millisecond)
		r, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		err = json.NewDecoder(r.Body).Decode(&csStuff)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range csStuff.retrieveNews() {
		ns.addNews(v.Title, v.URL)
	}

	newsDB = append(newsDB, ns)
}

func run(t string, wg *sync.WaitGroup) {
	if t == "" {
		log.Fatal("choose the type")
	}

	for _, v := range fetchersDB {
		if v.token == t {
			for _, cb := range v.fetchList {
				wg.Add(1)
				go cb()
			}
		}
	}
}

func logStats() {
	elapsed := color.New(color.BgRed)
	fmt.Printf("\nElapsed time: ")
	elapsed.Printf(" %v ", time.Since(start))
	fmt.Println()
}

func logNews() {
	for _, news := range newsDB {
		group := color.New(color.Bold, color.FgYellow)
		group.Printf("\n%s - %s\n\n", news.Group, news.URL)

		for _, newsInfo := range news.Info {
			n := color.New(color.FgYellow)
			t := color.New(color.FgWhite)
			indexFmt := "0%d. "
			if newsInfo.Index > 9 {
				indexFmt = "%d. "
			}

			n.Printf(indexFmt, newsInfo.Index)
			t.Printf("%s", newsInfo.Title)

			if newsInfo.URL != "" {
				color.White("\n    %s", newsInfo.URL)
			}

			fmt.Println()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	t := flag.String("t", "", "type")
	flag.Parse()

	registerFetchers(&wg)
	run(*t, &wg)
	wg.Wait()
	logNews()
	logStats()
}
