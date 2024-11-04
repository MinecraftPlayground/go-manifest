package manifest

import (
	"fmt"
	"time"
)

type Manifest struct {
	Latest   ManifestLatest    `json:"latest"`
	Versions []ManifestVersion `json:"versions"`
}

func (manifest Manifest) GetLatestRelease() (version ManifestVersion) {
	version, err := manifest.FindVersion(manifest.Latest.Release)
	if err != nil {
		panic(fmt.Errorf("get latest release: %w", err))
	}
	return version
}

func (manifest Manifest) GetLatestSnapshot() (version ManifestVersion) {
	version, err := manifest.FindVersion(manifest.Latest.Snapshot)
	if err != nil {
		panic(fmt.Errorf("get latest snapshot: %w", err))
	}
	return version
}

func (manifest Manifest) FindVersion(v string) (version ManifestVersion, err error) {
	for _, version := range manifest.Versions {
		if version.ID == v {
			return version, nil
		}
	}
	return ManifestVersion{}, fmt.Errorf("find version: version named '%s' does not exist in this manifest", v)
}

type ManifestLatest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type ManifestVersion struct {
	ID          string    `json:"id"`
	VersionType string    `json:"type"`
	URL         string    `json:"url"`
	ReleaseTime time.Time `json:"releaseTime"`
	Hash        string    `json:"sha1"`
}
