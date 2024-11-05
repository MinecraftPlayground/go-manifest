package manifest

import (
	"bytes"
	"fmt"
)

type Version struct {
	Arguments  VersionArguments  `json:"arguments"`
	AssetIndex VersionAssetIndex `json:"assetIndex"`
	Assets     string            `json:"assets"`
	Downloads  VersionDownloads  `json:"downloads"`
}

type VersionArguments struct {
}

type VersionAssetIndex struct {
}

type VersionDownloads struct {
	Client         Download `json:"client"`
	ClientMappings Download `json:"client_mappings"`
	Server         Download `json:"server"`
	ServerMappings Download `json:"server_mappings"`
}

func (version Version) GetClient() (client *bytes.Buffer, err error) {
	client, err = version.Downloads.Client.download()
	if err != nil {
		err = fmt.Errorf("get client: %w", err)
	}
	return client, err
}
