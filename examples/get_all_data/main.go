package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MinecraftPlayground/go-manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: get_all_data <version_id>")
		os.Exit(1)
	}
	start := time.Now()

	fmt.Printf("Fetching client JAR for version %s...\n", os.Args[1])
	client, err := manifest.GetClient(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Extracting files...")
	r := bytes.NewReader(client.Bytes())
	zipReader, err := zip.NewReader(r, r.Size())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	for _, file := range zipReader.File {
		if !strings.HasPrefix(file.Name, "data/") || strings.HasPrefix(file.Name, "data/.") {
			continue
		}
		wg.Add(1)
		go writeFile(file, wg)
	}
	wg.Wait()
	fmt.Printf("\r%-*s\n", 160,
		fmt.Sprintf("Done! (%s elapsed)", time.Since(start).Round(10*time.Millisecond)),
	)
}

func writeFile(file *zip.File, wg *sync.WaitGroup) {
	splitPath := strings.Split(file.Name, "/")
	err := os.MkdirAll(strings.Join(splitPath[:len(splitPath)-1], "/"), 0755)
	if err != nil {
		fmt.Println("Error creating directory: ", err)
		os.Exit(1)
	}

	fileReader, err := file.OpenRaw()
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	data, err := io.ReadAll(fileReader)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}
	err = os.WriteFile(file.Name, data, 0644)
	if err != nil {
		fmt.Println("Error writing file: ", err)
		return
	}

	fmt.Print(fmt.Sprintf("\rExtracted: %-*s", 150, file.Name))
	wg.Done()
}
