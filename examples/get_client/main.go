package main

import (
	"fmt"
	"os"

	"github.com/Kesuaheli/go-manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: get_client <version_id>")
		os.Exit(1)
	}

	fmt.Printf("Fetching client JAR for version %s...\n", os.Args[1])

	client, err := manifest.GetClient(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(fmt.Sprintf("client_%s.jar", os.Args[1]), client.Bytes(), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Saved client JAR for version %s to client_%[1]s.jar\n", os.Args[1])
}
