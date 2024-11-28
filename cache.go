package manifest

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/patrickmn/go-cache"
)

var (
	versionCache *cache.Cache
)

func init() {
	gob.Register(&Version{})

	finfo, _ := os.Stat(".cache/versions")
	if finfo == nil {
		versionCache = cache.New(cache.NoExpiration, cache.NoExpiration)
		return
	}
	f, err := os.Open(".cache/versions")
	if err != nil {
		fmt.Printf("Failed to open cache file: %+v\n", err.(*os.PathError))
		return
	}
	defer f.Close()

	items := make(map[string]cache.Item)
	err = gob.NewDecoder(f).Decode(&items)
	if err != nil {
		fmt.Printf("Failed to decode cache file: %+v\n", err.(*os.PathError))
		return
	}
	versionCache = cache.NewFrom(cache.NoExpiration, cache.NoExpiration, items)
	fmt.Printf("Loaded cache from .cache/versions (%d items)\n", versionCache.ItemCount())
}

// SaveCache saves the cache to the .cache folder
func SaveCache() {
	versionCache.DeleteExpired()
	finfo, _ := os.Stat(".cache")
	var err error
	if finfo == nil {
		err = os.Mkdir(".cache", 0755)
		if err != nil {
			fmt.Printf("Failed to create cache directory: %+v\n", err.(*os.PathError))
			return
		}
	}

	finfo, _ = os.Stat(".cache/versions")
	var f *os.File
	if finfo == nil {
		f, err = os.Create(".cache/versions")
		if err != nil {
			fmt.Printf("Failed to create cache file: %+v\n", err.(*os.PathError))
			return
		}
	} else {
		f, err = os.OpenFile(".cache/versions", os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to create cache file: %+v\n", err.(*os.PathError))
			return
		}
	}
	defer f.Close()

	err = gob.NewEncoder(f).Encode(versionCache.Items())
	if err != nil {
		fmt.Printf("Failed to save cache: %+v\n", err)
		return
	}
	fmt.Printf("Saved cache to .cache/versions (%d items)\n", versionCache.ItemCount())
}
