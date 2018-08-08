package hn

import (
	"sync"
)

func GetHN(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherNewsFromHackerNews(wg)
}