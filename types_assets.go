package manifest

import (
	"bytes"
	"fmt"
)

// Assets represents the assets of a version
type Assets struct {
	// Objects is a map from file path of the assets to its download
	Objects map[string]AssetDownload `json:"objects"`
}

// AssetDownload is similar	to [Download], but has a different implementation of its download.
type AssetDownload struct {
	// Hash is the SHA1 hash of the asset file. It is used to construct the URL and verify the file content.
	Hash string `json:"hash"`

	Download
}

func (asset AssetDownload) download() (*bytes.Buffer, error) {
	asset.Download.Hash = asset.Hash
	asset.URL = fmt.Sprintf("%s/%s/%s", RESOURCE_URL, asset.Hash[:2], asset.Hash)
	return asset.Download.download()
}
