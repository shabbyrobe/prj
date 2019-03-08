package prj

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/shabbyrobe/golib/errtools"
)

const DefaultHashAlgorithm = HashSHA512

type Hash struct {
	Algorithm HashAlgorithm
	Value     HashValue
}

func (h Hash) IsEmpty() bool {
	return h.Algorithm == HashNone && len(h.Value) == 0
}

func (h Hash) String() string {
	if h.IsEmpty() {
		return ""
	}
	return fmt.Sprintf("%s:%s", h.Algorithm, h.Value.String())
}

func (h Hash) Equal(c Hash) (equal bool, rerr error) {
	if h.Algorithm == HashNone || c.Algorithm == HashNone {
		return false, fmt.Errorf("prj: hash has no algo: %q != %q", h.String(), c.String())
	}
	if h.Algorithm != c.Algorithm {
		return false, fmt.Errorf("prj: hashes use different algorithm: %q != %q", h.String(), c.String())
	}
	return bytes.Equal(h.Value, c.Value), nil
}

func (v Hash) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", v.String())), nil
}

func (v *Hash) UnmarshalJSON(bts []byte) error {
	var s string
	if err := json.Unmarshal(bts, &s); err != nil {
		return err
	}

	h, err := ParseHash(s)
	if err != nil {
		return err
	}

	*v = h
	return nil
}

func ParseHash(v string) (h Hash, rerr error) {
	parts := strings.SplitN(strings.TrimSpace(v), ":", 2)
	if len(parts) != 2 {
		return h, fmt.Errorf("prj: hash %q format invalid, expected '<algo>:<value>'", v)
	}

	h.Algorithm = HashAlgorithm(parts[0])
	if !h.Algorithm.IsValid() {
		return h, fmt.Errorf("prj: hash %q algorithm invalid %q", v, h.Algorithm)
	}
	bv, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return h, err
	}
	h.Value = HashValue(bv)

	return h, nil
}

type HashValue []byte

func (v HashValue) String() string { return base64.URLEncoding.EncodeToString(v) }

type HashAlgorithm string

const (
	HashNone   HashAlgorithm = ""
	HashSHA512 HashAlgorithm = "sha512"
)

func (ha HashAlgorithm) IsValid() bool {
	return ha == HashSHA512
}

func (ha HashAlgorithm) Sum(hasher hash.Hash, bts []byte) Hash {
	var hash Hash
	hash.Algorithm = ha
	hash.Value = hasher.Sum(bts)
	return hash
}

func (ha HashAlgorithm) HashFile(file string) (fh Hash, rerr error) {
	var f *os.File
	if f, rerr = os.Open(file); rerr != nil {
		return fh, rerr
	}
	defer errtools.DeferClose(&rerr, f)

	return ha.Hash(f)
}

func (ha HashAlgorithm) Hash(rdr io.Reader) (Hash, error) {
	var hash Hash
	hash.Algorithm = ha
	hasher, err := ha.CreateHasher()
	if err != nil {
		return hash, err
	}

	if _, err := io.Copy(hasher, rdr); err != nil {
		return hash, err
	}

	hash.Value = hasher.Sum(nil)
	return hash, nil
}

func (ha HashAlgorithm) CreateHasher() (hash.Hash, error) {
	switch ha {
	case HashSHA512:
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("prj: unsupported hash: %q", ha)
	}
}
