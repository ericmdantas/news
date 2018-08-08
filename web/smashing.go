package web

import (
	"fmt"
	"strings"
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func GatherNewsFromSmashingMagazine(wg *sync.WaitGroup) {
	defer wg.Done()

	var newsList []string
	
	res, err := http.Get("https://www.smashingmagazine.com/")
	
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	
	if err != nil {
		panic(err)
	}
	
	doc.Find(".featured-article__title>a").Each(func(i int, s *goquery.Selection) {
		news := s.Text()
		if (strings.Contains(news, "Meow")) {
			return 
		}
		
		newsList = append(newsList, news)
	})
	
	fmt.Println("\n\n[SMASHING MAGAZINE]\n")
	
	for _, v := range newsList {
		fmt.Printf("-> %s\n", v)
	}
}