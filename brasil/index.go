package brasil

import (
	"sync"
)

func GetBrasil(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherNewsFromG1Geral(wg)
	wg.Add(1)
	go GatherNewsFromG1RegiaoSerrana(wg)
}