package manifest

import (
	"fmt"
	"io"
	"net/http"
)

type Version struct {
	Arguments  VersionArguments  `json:"arguments"`
	AssetIndex VersionAssetIndex `json:"assetIndex"`
	Assets     string            `json:"assets"`
	Downloads  VersionDownloads  `json:"downloads"`
	// ...
}

type VersionArguments struct {
}

type VersionAssetIndex struct {
}

type VersionDownloads map[string]VersionDownload

type VersionDownload struct {
	Hash     string `json:"sha1"`
	FileSize int    `json:"size"`
	URL      string `json:"url"`
}

func (version Version) GetClient() (io.Reader, error) {
	resp, err := http.Get(version.Downloads["client"].URL)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	// TODO: handle response: check status code, verify response version.Downloads["client"].FileSize and version.Downloads["client"].Hash

	return resp.Body, nil
}
