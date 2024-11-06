package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Version struct {
	Arguments              VersionArguments `json:"arguments"`
	AssetIndex             Download         `json:"assetIndex"`
	Assets                 string           `json:"assets"`
	Downloads              VersionDownloads `json:"downloads"`
	ID                     string           `json:"id"`
	JavaVersion            VersionJava      `json:"javaVersion"`
	Libraries              []VersionLibrary `json:"libraries"`
	Logging                VersionLogging   `json:"logging"`
	MainClass              string           `json:"mainClass"`
	MinimumLauncherVersion int              `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time        `json:"releaseTime"`
	Type                   string           `json:"type"`
}

func (version Version) GetClient() (client *bytes.Buffer, err error) {
	client, err = version.Downloads.Client.download()
	if err != nil {
		err = fmt.Errorf("get client: %w", err)
	}
	return client, err
}

func (version Version) GetAssetIndex() (assets Assets, err error) {
	assetsData, err := version.AssetIndex.download()
	if err != nil {
		return Assets{}, fmt.Errorf("get asset index: %w", err)
	}

	err = json.NewDecoder(assetsData).Decode(&assets)
	if err != nil {
		err = fmt.Errorf("get asset index: %w", err)
	}
	return assets, err
}

func (version Version) GetAssetFile(path string) (buf *bytes.Buffer, err error) {
	assets, err := version.GetAssetIndex()
	if err != nil {
		return nil, err
	}

	assetDownload, ok := assets.Objects[path]
	if !ok {
		return nil, fmt.Errorf("get asset: file '%s' not found in version %s", path, version.ID)
	}
	buf, err = assetDownload.download()
	if err != nil {
		err = fmt.Errorf("get asset: %w", err)
	}
	return buf, err
}

func (version Version) GetAllAssets() (assetMap map[string]*bytes.Buffer, err error) {
	assets, err := version.GetAssetIndex()
	if err != nil {
		return nil, err
	}

	assetMap = make(map[string]*bytes.Buffer, len(assets.Objects))
	i := 0
	lastLen := 0
	defer fmt.Println()
	for path, asset := range assets.Objects {
		i++
		fmt.Printf("\r% *s\r", lastLen, "")
		log := fmt.Sprintf("- downloading %0*d/%d (%03d%%): %s...",
			len(fmt.Sprintf("%d", len(assets.Objects))), i,
			len(assets.Objects),
			i*100/len(assets.Objects),
			path,
		)
		fmt.Print(log)
		lastLen = len(log)
		buf, err := asset.download()
		if err != nil {
			return nil, fmt.Errorf("get all assets ('%s'): %w", path, err)
		}
		assetMap[path] = buf
	}
	return assetMap, nil

}

type VersionArguments struct {
}

type VersionDownloads struct {
	Client         Download `json:"client"`
	ClientMappings Download `json:"client_mappings"`
	Server         Download `json:"server"`
	ServerMappings Download `json:"server_mappings"`
}

type VersionJava struct {
	Component string `json:"component"`
	Major     int    `json:"majorVersion"`
}

type VersionLibrary struct {
	Name      string
	Downloads VersionLibraryDownload `json:"downloads"`
}

type VersionLibraryDownload struct {
	Artifact VersionLibraryDownloadArtifact `json:"artifact"`
}

type VersionLibraryDownloadArtifact struct {
	Path string `json:"path"`
	Download
}

type VersionLogging struct {
	Client VersionLogger `json:"client"`
}

type VersionLogger struct {
	Argument string   `json:"argument"`
	File     Download `json:"file"`
	Type     string   `json:"type"`
}
