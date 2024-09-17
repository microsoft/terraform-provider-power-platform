// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func GetCertificateRawFromCertOrFilePath(certificate, certificateFilePath string) (string, error) {
	if certificate != constants.EMPTY {
		return strings.TrimSpace(certificate), nil
	}
	if certificateFilePath != constants.EMPTY {
		pfx, err := os.ReadFile(certificateFilePath)
		if err != nil {
			return constants.EMPTY, err
		}
		certAsBase64 := base64.StdEncoding.EncodeToString(pfx)
		return strings.TrimSpace(certAsBase64), nil
	}
	return constants.EMPTY, errors.New("either client_certificate base64 or certificate_file_path must be provided")
}

func ConvertBase64ToCert(b64, password string) ([]*x509.Certificate, crypto.PrivateKey, error) {
	pfx, err := convertBase64ToByte(b64)
	if err != nil {
		return nil, nil, err
	}

	certs, key, err := convertByteToCert(pfx, password)
	if err != nil {
		return nil, nil, err
	}

	return certs, key, nil
}

func convertBase64ToByte(b64 string) ([]byte, error) {
	if b64 == constants.EMPTY {
		return nil, errors.New("got empty base64 certificate data")
	}

	pfx, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return pfx, fmt.Errorf("could not decode base64 certificate data: %w", err)
	}

	return pfx, nil
}

func convertByteToCert(certData []byte, password string) ([]*x509.Certificate, crypto.PrivateKey, error) {
	var key crypto.PrivateKey

	key, cert, _, err := pkcs12.DecodeChain(certData, password)
	if err != nil {
		return nil, nil, err
	}

	if cert == nil {
		return nil, nil, errors.New("found no certificate")
	}

	certs := []*x509.Certificate{cert}

	return certs, key, nil
}
