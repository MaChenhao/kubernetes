/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dynamiccertificates

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
)

type staticCAContent struct {
	name     string
	caBundle []byte
}

// NewStaticCAContentFromFile returns a CAContentProvider based on a filename
func NewStaticCAContentFromFile(filename string) (CAContentProvider, error) {
	if len(filename) == 0 {
		return nil, fmt.Errorf("missing filename for ca bundle")
	}

	caBundle, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewStaticCAContent(filename, caBundle), nil
}

// NewStaticCAContent returns a CAContentProvider that always returns the same value
func NewStaticCAContent(name string, caBundle []byte) CAContentProvider {
	return &staticCAContent{
		name:     name,
		caBundle: caBundle,
	}
}

// Name is just an identifier
func (c *staticCAContent) Name() string {
	return c.name
}

// CurrentCABundleContent provides ca bundle byte content
func (c *staticCAContent) CurrentCABundleContent() (cabundle []byte) {
	return c.caBundle
}

type staticCertKeyContent struct {
	name string
	cert []byte
	key  []byte
}

type staticSNICertKeyContent struct {
	staticCertKeyContent
	sniNames []string
}

// NewStaticCertKeyContentFromFiles returns a CertKeyContentProvider based on a filename
func NewStaticCertKeyContentFromFiles(certFile, keyFile string) (CertKeyContentProvider, error) {
	if len(certFile) == 0 {
		return nil, fmt.Errorf("missing filename for certificate")
	}
	if len(keyFile) == 0 {
		return nil, fmt.Errorf("missing filename for key")
	}

	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	return NewStaticCertKeyContent(fmt.Sprintf("cert: %s, key: %s", certFile, keyFile), certPEMBlock, keyPEMBlock)
}

// NewStaticSNICertKeyContentFromFiles returns a SNICertKeyContentProvider based on a filename
func NewStaticSNICertKeyContentFromFiles(certFile, keyFile string, sniNames ...string) (SNICertKeyContentProvider, error) {
	if len(certFile) == 0 {
		return nil, fmt.Errorf("missing filename for certificate")
	}
	if len(keyFile) == 0 {
		return nil, fmt.Errorf("missing filename for key")
	}

	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	return NewStaticSNICertKeyContent(fmt.Sprintf("cert: %s, key: %s", certFile, keyFile), certPEMBlock, keyPEMBlock, sniNames...)
}

// NewStaticCertKeyContent returns a CertKeyContentProvider that always returns the same value
func NewStaticCertKeyContent(name string, cert, key []byte) (CertKeyContentProvider, error) {
	// Ensure that the key matches the cert and both are valid
	_, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &staticCertKeyContent{
		name: name,
		cert: cert,
		key:  key,
	}, nil
}

// NewStaticSNICertKeyContent returns a SNICertKeyContentProvider that always returns the same value
func NewStaticSNICertKeyContent(name string, cert, key []byte, sniNames ...string) (SNICertKeyContentProvider, error) {
	// Ensure that the key matches the cert and both are valid
	_, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &staticSNICertKeyContent{
		staticCertKeyContent: staticCertKeyContent{
			name: name,
			cert: cert,
			key:  key,
		},
		sniNames: sniNames,
	}, nil
}

// Name is just an identifier
func (c *staticCertKeyContent) Name() string {
	return c.name
}

// CurrentCertKeyContent provides cert and key content
func (c *staticCertKeyContent) CurrentCertKeyContent() ([]byte, []byte) {
	return c.cert, c.key
}

func (c *staticSNICertKeyContent) SNINames() []string {
	return c.sniNames
}
