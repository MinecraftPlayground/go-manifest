package main

import (
	"fmt"
	"time"

	"github.com/Kesuaheli/go-manifest"
)

func main() {
	// make sure the cache will be saved as files on exit
	defer manifest.SaveCache()

	for range make([]int, 5) {
		now := time.Now()
		_, err := manifest.GetVersion("1.7.10")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got version in %s\n", time.Since(now))
		time.Sleep(2 * time.Second)
	}
}
