package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Version represents a Minecraft version definition.
type Version struct {
	// Arguments contains the arguments for this version.
	Arguments VersionArguments `json:"arguments"`

	// AssetIndex contains the asset index file for this version.
	AssetIndex Download `json:"assetIndex"`

	// Assets is the asset id for this version.
	// FIXME Need more info
	Assets string `json:"assets"`

	// Downloads contains the version downloads like client and server.
	Downloads VersionDownloads `json:"downloads"`

	// ID is the version name.
	ID string `json:"id"`

	// JavaVersion contains the java version information for this version.
	JavaVersion VersionJava `json:"javaVersion"`

	// Libraries contains the libraries for this version.
	Libraries []VersionLibrary `json:"libraries"`

	// Logging contains the logging information for this version.
	// FIXME Need more info
	Logging VersionLogging `json:"logging"`

	// MainClass is the name of the main class to run the client with.
	MainClass string `json:"mainClass"`

	// MinimumLauncherVersion is the minimum launcher version to run this version.
	// FIXME Need more info
	MinimumLauncherVersion int `json:"minimumLauncherVersion"`

	// ReleaseTime is the time this version was released.
	ReleaseTime time.Time `json:"releaseTime"`

	// Type is the type of this version like "release" or "snapshot".
	Type string `json:"type"`
}

// GetClient fetches the client JAR for this version.
func (version Version) GetClient() (client *bytes.Buffer, err error) {
	client, err = version.Downloads.Client.download()
	if err != nil {
		err = fmt.Errorf("get client: %w", err)
	}
	return client, err
}

// GetAssetIndex fetches the asset index for this version.
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

// GetAssetFile fetches the given asset file for this version.
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

// GetAllAssets fetches all assets for this version.
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

// VersionArguments represents the JAVA arguments for this version.
type VersionArguments struct {
	// TODO implement
}

// VersionDownloads represents the downloads for this version.
type VersionDownloads struct {
	// Client is the client JAR for this version.
	Client Download `json:"client"`

	// ClientMappings is the client mappings file for the client JAR for this version.
	ClientMappings Download `json:"client_mappings"`

	// Server is the server JAR for this version.
	Server Download `json:"server"`

	// ServerMappings is the server mappings file for the server JAR for this version.
	ServerMappings Download `json:"server_mappings"`
}

// VersionJava represents the Java version for a version.
type VersionJava struct {
	// Component is the component of the Java version.
	// FIXME Need more info
	Component string `json:"component"`

	// Major is the major version of the Java version.
	Major int `json:"majorVersion"`
}

// VersionLibrary represents a library for this version.
type VersionLibrary struct {
	// Name is the name of the library.
	Name string

	// Downloads contains the downloads for the library.
	Downloads VersionLibraryDownload `json:"downloads"`
}

// VersionLibraryDownload represents the downloads for a library.
type VersionLibraryDownload struct {
	// Artifact is the artifact for a library.
	Artifact VersionLibraryDownloadArtifact `json:"artifact"`
}

// VersionLibraryDownloadArtifact represents the artifact for a library.
type VersionLibraryDownloadArtifact struct {
	// Path is the path for a library artifact.
	Path string `json:"path"`

	Download
}

// VersionLogging represents the logging for this version.
// FIXME Need more info
type VersionLogging struct {
	// Client contains the client logger for this version.
	// FIXME Need more info
	Client VersionLogger `json:"client"`
}

// VersionLogger represents a logger for this version.
type VersionLogger struct {
	// Argument is the java argument for the logger.
	Argument string `json:"argument"`

	// File is the download for the logger.
	File Download `json:"file"`

	// Type is the type of the logger.
	Type string `json:"type"`
}
