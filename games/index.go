package games

import (
	"sync"
)

func GetGames(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherCSGoNewsFromReddit(wg)
	wg.Add(1)
	go GatherCSGoNewsFromHLTV(wg)
}