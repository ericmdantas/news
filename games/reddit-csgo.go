package games

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sync"
)

const (
	SLEEP_TIME = 3 * time.Second
	URL_THREAD = "https://www.reddit.com/r/globaloffensive/.json"
)

var (
	startTime = time.Now()
)

type payloadThread struct {
	Data struct {
		Children []struct {
			Data struct {
				ID        string `json:"id"`
				Title     string `json:"title"`
				Permalink string `json:"permalink"`
				Ups       int    `json:"ups"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (p payloadThread) IsEmpty() bool {
	return len(p.Data.Children) == 0
}

func (p payloadThread) log() {
	fmt.Println("\n\n[Reddit]\n")

	for _, child := range p.Data.Children {
		fmt.Printf("\n-> %s", child.Data.Title)
	}

	fmt.Println()
}

func searchThreads() payloadThread {
	var p payloadThread
	res, err := http.Get(URL_THREAD)

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(res.Body).Decode(&p)

	if err != nil {
		panic(err)
	}

	return p
}

func GatherCSGoNewsFromReddit(wg *sync.WaitGroup) {
	p := searchThreads()	

	for {
		if !p.IsEmpty() {
			break
		}

		time.Sleep(SLEEP_TIME)
		p = searchThreads()
	}

	p.log()
	
	wg.Done()
}
