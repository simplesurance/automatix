package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// PathExists returns true if the path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

// IsDir returns true if the path is a directory
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

// Mkdir creates recursively directories
func Mkdir(path string) error {
	return os.MkdirAll(path, os.FileMode(0755))
}

// IsFile returns true if path is a file.
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !fi.IsDir()
}

// Sha256Hash returns the SHA256 hash of a file
func Sha256Hash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "opening file failed")
	}

	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "reading file failed")
	}
	sum := h.Sum(nil)

	return fmt.Sprintf("%x", sum), nil
}
