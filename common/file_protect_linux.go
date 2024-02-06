//go:build linux
// +build linux

package common

import (
	"os"
	"path/filepath"

	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

// On non Windows platforms, we don't encrypt the cache file
type FileProtectData struct {
}

func (pd *FileProtectData) GetOrCreateCacheFileDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homedir, constants.CACHE_SAVE_FOLDER_PATH_LINUX)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}
	return dir, nil
}

func (pd *FileProtectData) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (pd *FileProtectData) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}
