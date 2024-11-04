package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BASE_URL     = "https://piston-meta.mojang.com/mc/game"
	RESOURCE_URL = "https://resources.download.minecraft.net"
)

func GetManifest() (manifest *Manifest, err error) {
	resp, err := http.Get(BASE_URL + "/version_manifest_v2.json")
	if err != nil {
		return &Manifest{}, fmt.Errorf("get manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return &Manifest{}, fmt.Errorf("get manifest: failed to read body on invalid response (%s): %w", resp.Status, err)
		}
		return &Manifest{}, fmt.Errorf("get manifest: invalid response: expected %d, got %s: %s", http.StatusOK, resp.Status, string(data))
	}

	manifest = &Manifest{}
	err = json.NewDecoder(resp.Body).Decode(manifest)
	if err != nil {
		err = fmt.Errorf("get manifest: failed to decode body: %w", err)
	}
	return manifest, err
}

func GetVersion(v string) (version *Version, err error) {
	manifest, err := GetManifest()
	if err != nil {
		return &Version{}, fmt.Errorf("get version: %w", err)
	}
	manifestVersion, err := manifest.FindVersion(v)
	if err != nil {
		return &Version{}, fmt.Errorf("get version: %w", err)
	}

	resp, err := http.Get(manifestVersion.URL)
	if err != nil {
		return &Version{}, fmt.Errorf("get version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Version{}, fmt.Errorf("get version: invalid response: expected %d, got %s", http.StatusOK, resp.Status)
	}

	version = &Version{}
	err = json.NewDecoder(resp.Body).Decode(version)
	if err != nil {
		err = fmt.Errorf("get version: failed to decode body: %w", err)
	}
	return version, err
}
