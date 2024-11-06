package manifest

import (
	"bytes"
	"fmt"
)

type Assets struct {
	Objects map[string]AssetDownload `json:"objects"`
}

type AssetDownload struct {
	Hash string `json:"hash"`
	Download
}

func (asset AssetDownload) download() (*bytes.Buffer, error) {
	asset.Download.Hash = asset.Hash
	asset.URL = fmt.Sprintf("%s/%s/%s", RESOURCE_URL, asset.Hash[:2], asset.Hash)
	return asset.Download.download()
}
