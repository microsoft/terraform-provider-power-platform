// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func CalculateMd5(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		//we return empty Md5 value if the file does not exist yet
		return "", nil
	}

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
