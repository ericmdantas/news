package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/ericmdantas/ns/brasil"
	"github.com/ericmdantas/ns/games"
	"github.com/ericmdantas/ns/hn"
	"github.com/ericmdantas/ns/science"
	"github.com/ericmdantas/ns/web"
)

const (
	typeWeb     = "web"
	typeGames   = "gg"
	typeBrasil  = "br"
	typeScience = "sci"
	typeHN      = "hn"
	typeAll     = "all"
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

	t := flag.String("type", "", "Type of ns to be shown")
	flag.Parse()

	if *t == "" {
		panic("escolhe um tipo de notícia ai, maluco")
		return
	}

	switch *t {
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
		panic("-type=hn|br|gg|sci|web|all")
	}

	wg.Wait()

	fmt.Printf("\n\n---> it took %v to gather it all :D\n", time.Since(start))
}
