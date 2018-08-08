package web

import (
	"fmt"
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func GatherNewsFromAListApart(wg *sync.WaitGroup) {
	defer wg.Done()

	var newsList []string
	
	res, err := http.Get("https://alistapart.com/articles")
	
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	
	if err != nil {
		panic(err)
	}
	
	doc.Find(".entry-title").Each(func(i int, s *goquery.Selection) {
		news := s.Text()
		newsList = append(newsList, news)
	})
	
	fmt.Println("\n\n[ALISTAPART]\n")
	
	for _, v := range newsList {
		fmt.Printf("-> %s\n", v)
	}
}