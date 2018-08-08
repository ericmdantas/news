package web

import (
	"sync"
)

func GetWeb(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherNewsFromSmashingMagazine(wg)
	wg.Add(1)
	go GatherNewsFromAListApart(wg)
}