package main

import (
	"encoding/json"
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

	data, err := json.MarshalIndent(v, "", "	")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.WriteFile(fmt.Sprintf("version_%s.json", os.Args[1]), data, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Saved version data of %s to version_%[1]s.json\n", os.Args[1])
}
