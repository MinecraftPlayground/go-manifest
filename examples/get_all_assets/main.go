package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Kesuaheli/go-manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: get_assets <version_id>")
		os.Exit(1)
	}

	fmt.Printf("Fetching assets for version %s...\n", os.Args[1])
	assets, err := manifest.GetAllAssets(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Got asset %d files of %s\n", len(assets), os.Args[1])
	lastLen := 0
	for path, asset := range assets {
		pathSplit := strings.Split(path, "/")
		err = os.MkdirAll("assets/"+strings.Join(pathSplit[:len(pathSplit)-1], "/"), 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		file, err := os.Create("assets/" + path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		_, err = file.Write(asset.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\r% *s\r", lastLen, "")
		log := fmt.Sprintf("- saved asset file %s to assets/%s!",
			pathSplit[len(pathSplit)-1],
			path,
		)
		fmt.Print(log)
		lastLen = len(log)
	}
	fmt.Println()
}
