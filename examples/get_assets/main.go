package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MinecraftPlayground/go-manifest"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: get_assets <version_id> <asset_id>")
		os.Exit(1)
	}

	v, err := manifest.GetAssetFile(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Got asset file %s of %s\n", os.Args[2], os.Args[1])
	path := strings.Split(os.Args[2], "/")
	err = os.MkdirAll("assets/"+strings.Join(path[:len(path)-1], "/"), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Create("assets/" + os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = file.Write(v.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Saved asset file %s to assets/%s\n", path[len(path)-1], os.Args[2])
}
