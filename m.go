package main

import (
	"sync"
	"flag"
	"fmt"
	"time"
	
	"github.com/ericmdantas/news/games"
	"github.com/ericmdantas/news/science"
	"github.com/ericmdantas/news/hn"
	"github.com/ericmdantas/news/brasil"
	"github.com/ericmdantas/news/web"
)

const (
	typeWeb = "web"
	typeGames = "gg"
	typeBrasil = "br"
	typeScience = "sci"
	typeHN = "hn"
	typeAll = "all"
)

func getAll(wg *sync.WaitGroup) {
	hn.GetHN(wg)
	science.GetScience(wg)
	brasil.GetBrasil(wg)
	games.GetGames(wg)		
	web.GetWeb(wg)
}

func main() {
	var wg sync.WaitGroup
	
	start := time.Now()
	
	t := flag.String("type", "", "Type of news to be shown")
	flag.Parse()

	if (*t == "") {
		panic("escolhe um tipo de notícia ai, maluco")
		return 
	}
	
	switch (*t) {
	case typeAll:
		getAll(&wg)
		break
	case typeGames:
		games.GetGames(&wg)
		break
	case typeBrasil:
		brasil.GetBrasil(&wg)
		break
	case typeScience:
		science.GetScience(&wg)
		break
	case typeHN:
		hn.GetHN(&wg)
		break
	case typeWeb:
		web.GetWeb(&wg)
		break
	default:
		panic("escolhe um tipo de notícia ai, maluco")
	}
	
	wg.Wait()
	
	fmt.Printf("\n\n---> it took %v to gather it all :D\n", time.Since(start))
}