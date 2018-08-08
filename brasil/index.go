package brasil

import (
	"sync"
)

func GetBrasil(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherNewsFromG1(wg)
}