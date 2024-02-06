//go:build windows
// +build windows

// https://stackoverflow.com/questions/33516053/windows-encrypted-rdp-passwords-in-golang
package common

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

const (
	CRYPTPROTECT_UI_FORBIDDEN = 0x1
)

var (
	dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

	procEncryptData = dllcrypt32.NewProc("CryptProtectData")
	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree   = dllkernel32.NewProc("LocalFree")
)

type PROTECT_DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func (b *PROTECT_DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

type FileProtectData struct {
}

func (pd *FileProtectData) GetOrCreateCacheFileDir() (string, error) {
	appDataLocal, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(appDataLocal, constants.CACHE_SAVE_FOLDER_PATH_WINDOWS)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}

	return dir, nil
}

func (pd *FileProtectData) newBlob(d []byte) *PROTECT_DATA_BLOB {
	if len(d) == 0 {
		return &PROTECT_DATA_BLOB{}
	}
	return &PROTECT_DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (pd *FileProtectData) Encrypt(data []byte) ([]byte, error) {
	var out_blob PROTECT_DATA_BLOB
	r, _, err := procEncryptData.Call(uintptr(unsafe.Pointer(pd.newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&out_blob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(out_blob.pbData)))
	return out_blob.ToByteArray(), nil
}

func (pd *FileProtectData) Decrypt(data []byte) ([]byte, error) {
	var out_blob PROTECT_DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(pd.newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&out_blob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(out_blob.pbData)))
	return out_blob.ToByteArray(), nil
}
