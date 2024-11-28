package main

import (
	"fmt"
	"os"

	"github.com/MinecraftPlayground/go-manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: get_version <version_id>")
		os.Exit(1)
	}

	v, err := manifest.GetVersion(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Version data of %s:\n%+v\n", os.Args[1], *v)
}
