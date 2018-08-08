package science

import (
	"fmt"
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func GatherNewsFromPhysOrg(wg *sync.WaitGroup) {
	defer wg.Done()

	var newsList []string
	
	res, err := http.Get("https://www.phys.org")
	
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	
	if err != nil {
		panic(err)
	}
	
	doc.Find(".news-box h3 a").Each(func(i int, s *goquery.Selection) {
		news := s.Text()
		newsList = append(newsList, news)
	})
	
	fmt.Println("\n\n[PHYS ORG]\n")
	
	for _, v := range newsList {
		fmt.Printf("-> %s\n", v)
	}
}