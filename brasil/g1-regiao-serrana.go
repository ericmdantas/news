package brasil

import (
	"fmt"
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func GatherNewsFromG1RegiaoSerrana(wg *sync.WaitGroup) {
	defer wg.Done()

	var newsList []string
	
	res, err := http.Get("http://g1.globo.com/rj/regiao-serrana/")
	
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	
	if err != nil {
		panic(err)
	}
	
	doc.Find(".feed-post-body-title").Each(func(i int, s *goquery.Selection) {
		news := s.Text()
		newsList = append(newsList, news)
	})
	
	fmt.Println("\n\n[G1 Região Serrana]\n")
	
	for _, v := range newsList {
		fmt.Printf("-> %s\n", v)
	}
}