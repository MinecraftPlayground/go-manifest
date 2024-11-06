package manifest

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BASE_URL     = "https://piston-meta.mojang.com"
	RESOURCE_URL = "https://resources.download.minecraft.net"
)

func GetManifest() (manifest *Manifest, err error) {
	resp, err := http.Get(BASE_URL + "/mc/game/version_manifest_v2.json")
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

func GetClient(v string) (client *bytes.Buffer, err error) {
	version, err := GetVersion(v)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}
	return version.GetClient()

}

type Download struct {
	Hash     string `json:"sha1"`
	FileSize int    `json:"size"`
	URL      string `json:"url"`
}

func (d Download) download() (*bytes.Buffer, error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("download: read data (%s) %w", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download: invalid response: expected %d, got %s: %s", http.StatusOK, resp.Status, body)
	}

	// verify file size
	if d.FileSize != len(body) {
		return nil, fmt.Errorf("download: invalid file size: expected %db, got %db", d.FileSize, len(body))
	}

	// verify hash
	hash := fmt.Sprintf("%x", sha1.Sum(body))
	if d.Hash != hash {
		return nil, fmt.Errorf("download: invalid hash: expected %s, got %s", d.Hash, hash)
	}

	return bytes.NewBuffer(body), nil
}
