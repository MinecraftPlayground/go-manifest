package manifest

import (
	"fmt"
	"time"
)

// Manifest represents a version manifest.
type Manifest struct {
	// Latest contains the latest release and snapshot.
	Latest ManifestLatest `json:"latest"`

	// Versions contains all versions.
	Versions []ManifestVersion `json:"versions"`
}

// GetLatestRelease returns the latest release version.
func (manifest Manifest) GetLatestRelease() (version ManifestVersion) {
	version, err := manifest.FindVersion(manifest.Latest.Release)
	if err != nil {
		panic(fmt.Errorf("get latest release: %w", err))
	}
	return version
}

// GetLatestSnapshot returns the latest snapshot version.
func (manifest Manifest) GetLatestSnapshot() (version ManifestVersion) {
	version, err := manifest.FindVersion(manifest.Latest.Snapshot)
	if err != nil {
		panic(fmt.Errorf("get latest snapshot: %w", err))
	}
	return version
}

// FindVersion returns the version with the given id.
//
// v can be any version id like "1.7.10" or "24w10a" or even "a1.0.4".
func (manifest Manifest) FindVersion(v string) (version ManifestVersion, err error) {
	for _, version := range manifest.Versions {
		if version.ID == v {
			return version, nil
		}
	}
	return ManifestVersion{}, fmt.Errorf("find version: version named '%s' does not exist in this manifest", v)
}

// ManifestLatest contains the latest release and snapshot.
type ManifestLatest struct {
	// Release is the name of the latest release version.
	Release string `json:"release"`

	// Snapshot is the name of the latest snapshot version.
	Snapshot string `json:"snapshot"`
}

// ManifestVersion represents a version definition in the manifest.
type ManifestVersion struct {
	// ID is the name of the version.
	ID string `json:"id"`

	// Type is the type of the version.
	// Can be "release", "snapshot", "old_beta" or "old_alpha".
	VersionType string `json:"type"`

	// URL is the download URL of the version details.
	URL string `json:"url"`

	// Hash is the hash of the version details file. It is used to verify the file.
	Hash string `json:"sha1"`

	// ReleaseTime is the release time of the version.
	ReleaseTime time.Time `json:"releaseTime"`
}
