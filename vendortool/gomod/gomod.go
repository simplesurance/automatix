package gomod

import (
	"fmt"
	"path/filepath"

	"github.com/simplesurance/automatix/exec"
	"github.com/simplesurance/automatix/fs"

	"github.com/pkg/errors"
)

var files = []string{
	"go.mod",
	"go.sum",
	"vendor",
}

// NewModule initializes a module object. Path must reference a go.mod file
func NewModule(path string) *Module {
	return &Module{
		dir:     filepath.Dir(path),
		modPath: path,
	}
}

// Module represents a Go module instance
type Module struct {
	dir     string
	modPath string
}

// FullUpdate updates the dependencies in the go.mod file  the lastest minor or
// patch release.
// True is returned if at least one dependency was updated, otherwise false.
func (m *Module) FullUpdate() (bool, error) {
	return m.update(false)
}

// SafeUpdate updates the dependencies in the go.mod file to the lastest minor
// release.
// True is returned if at least one dependency was updated, otherwise false.
func (m *Module) SafeUpdate() (bool, error) {
	return m.update(true)
}

func (m *Module) update(patchOnlyUpdate bool) (bool, error) {
	hashBefore, err := fs.Sha256Hash(m.modPath)
	if err != nil {
		return false, errors.Wrap(err, "hashing go.mod file failed")
	}

	cmd := fmt.Sprintf("go get -u")
	if patchOnlyUpdate {
		cmd += "=patch"
	}

	_, err = exec.Command(m.dir, cmd)
	if err != nil {
		return false, errors.Wrapf(err, "running %q failed", cmd)
	}

	hashAfter, err := fs.Sha256Hash(m.modPath)
	if err != nil {
		return false, errors.Wrap(err, "hashing go.mod file failed")
	}

	return !(hashBefore == hashAfter), nil
}

// VendorDependencies downloads the dependencies into the vendor/ directory. If
// the directory already exist, it is updated to represent the information in
// the go.mod file.
func (m *Module) VendorDependencies() error {
	cmd := "go mod vendor"
	_, err := exec.Command(m.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	return nil
}

// Files returns a list of files and directories that belong to the vendor
// tool.
func (m *Module) Files() []string {
	return files
}
