package powerplatform_helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func CalculateMd5(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
