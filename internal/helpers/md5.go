// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func CalculateMd5(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// we return empty Md5 value if the file does not exist yet.
		return constants.EMPTY, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return constants.EMPTY, err
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return constants.EMPTY, err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
