// SPDX-License-Identifier: Apache-2.0

package meta

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/opensbom-generator/parsers/internal/license"
)

// Package is the package abstraction that the parsers return
type Package struct {
	Version                 string `json:"version,omitempty"`
	Name                    string `json:"name"`
	Path                    string `json:"path,omitempty"`
	LocalPath               string `json:"dir"`
	Supplier                Supplier
	PackageURL              string `json:"purl"`
	Checksum                Checksum
	PackageHomePage         string `json:"homePage"`
	PackageDownloadLocation string `json:"downloadLocation"`
	LicenseConcluded        string `json:"licenseConcluded"`
	LicenseDeclared         string `json:"licenseDeclared"`
	CommentsLicense         string `json:"licenseComments"`
	OtherLicense            []license.License
	Copyright               string `json:"copyright"`
	PackageComment          string `json:"comment"`
	Root                    bool
	Packages                map[string]*Package
}

// TypeContact ...
type SupplierType string

const (
	Person       SupplierType = "Person"
	Organization SupplierType = "Organization"
)

// Supplier abstracts the supplier of the package
type Supplier struct {
	Type            SupplierType
	Name            string
	Email           string
	FuncGetSupplier func() string `json:"-"`
}

func (s *Supplier) emailIsEmpty() bool {
	email := strings.ToLower(s.Email)
	return (len(s.Email) == 0) ||
		(strings.Compare(email, "none") == 0) ||
		(strings.Compare(email, "unknown") == 0)
}

// Get default supplier based on Name value or let each plugin build its own logic
func (s *Supplier) Get() string {
	if s.FuncGetSupplier != nil {
		return s.FuncGetSupplier()
	}

	if s.Name == "" {
		return ""
	}

	if s.Type == "" {
		s.Type = Organization
	}

	pkgSupplier := fmt.Sprintf("%s: %s", s.Type, s.Name)
	if !s.emailIsEmpty() {
		pkgSupplier += fmt.Sprintf(" (%s)", s.Email)
	}

	return pkgSupplier
}

type Checksum struct {
	Algorithm HashAlgorithm
	Content   []byte
	Value     string
}

func (c *Checksum) String() string {
	if c.Value == "" {
		c.Value = c.Compute(c.Content)
	}
	return c.Value
}

func (c *Checksum) Compute(content []byte) string {
	var h hash.Hash
	switch c.Algorithm {
	case HashAlgoSHA256:
		h = sha256.New()
	case HashAlgoSHA512:
		h = sha512.New()
	default:
		h = sha1.New()
	}
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

// HashAlgorithm ...
type HashAlgorithm string

const (
	HashAlgoSHA1        HashAlgorithm = "SHA1"
	HashAlgoSHA224      HashAlgorithm = "SHA224"
	HashAlgoSHA256      HashAlgorithm = "SHA256"
	HashAlgoSHA384      HashAlgorithm = "SHA384"
	HashAlgoSHA512      HashAlgorithm = "SHA512"
	HashAlgoMD2         HashAlgorithm = "MD2"
	HashAlgoMD4         HashAlgorithm = "MD4"
	HashAlgoMD5         HashAlgorithm = "MD5"
	HashAlgoMD6         HashAlgorithm = "MD6"
	HashAlgoUnsupported HashAlgorithm = "unsupported"
)

// GetHashAlgorithm takes a string and returns a HashAlgorithm type
func GetHashAlgorithm(h string) HashAlgorithm {
	switch u := strings.ToUpper(h); u {
	case "SHA1":
		return HashAlgoSHA1
	case "SHA224":
		return HashAlgoSHA224
	case "SHA256":
		return HashAlgoSHA256
	case "SHA384":
		return HashAlgoSHA384
	case "SHA512":
		return HashAlgoSHA512
	case "MD2":
		return HashAlgoMD2
	case "MD4":
		return HashAlgoMD2
	case "MD5":
		return HashAlgoMD5
	case "MD6":
		return HashAlgoMD6
	default:
		return HashAlgoUnsupported
	}
}
