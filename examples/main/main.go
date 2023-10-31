package main

import (
	"sync"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	user.StartServer(wg)

	wg.Wait()
}
