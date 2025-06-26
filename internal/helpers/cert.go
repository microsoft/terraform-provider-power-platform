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

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func GetCertificateRawFromCertOrFilePath(certificate, certificateFilePath string) (string, error) {
	if certificate != "" {
		return strings.TrimSpace(certificate), nil
	}
	if certificateFilePath != "" {
		pfx, err := os.ReadFile(certificateFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to read certificate file '%s': %w", certificateFilePath, err)
		}
		certAsBase64 := base64.StdEncoding.EncodeToString(pfx)
		return strings.TrimSpace(certAsBase64), nil
	}
	return "", errors.New("either client_certificate base64 or certificate_file_path must be provided")
}

func ConvertBase64ToCert(b64, password string) ([]*x509.Certificate, crypto.PrivateKey, error) {
	pfx, err := convertBase64ToByte(b64)
	if err != nil {
		return nil, nil, err
	}

	certs, key, err := convertByteToCert(pfx, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert certificate bytes: %w", err)
	}

	return certs, key, nil
}

func convertBase64ToByte(b64 string) ([]byte, error) {
	if b64 == "" {
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
		return nil, nil, fmt.Errorf("failed to decode PKCS12 certificate chain: %w", err)
	}

	if cert == nil {
		return nil, nil, errors.New("found no certificate")
	}

	certs := []*x509.Certificate{cert}

	return certs, key, nil
}
