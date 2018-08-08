package science

import (
	"sync"
)

func GetScience(wg *sync.WaitGroup) {
	wg.Add(1)
	go GatherNewsFromPhysOrg(wg)
}