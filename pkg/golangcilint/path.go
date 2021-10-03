package golangcilint

import (
	"go/build"
	"os"
	"path"

	"github.com/pkg/errors"
)

// GetPathForVersion returns the full path for a specific version of golangci-lint.
func GetPathForVersion(version string) (string, error) {
	return GetGobinPath(GetNameForVersion(version))
}

// GetNameForVersion returns the file name for a specific version of golangci-lint.
func GetNameForVersion(version string) string {
	return "golangci-lint@" + version
}

// GetDefaultPath returns the full path for the default version of golangci-lint.
func GetDefaultPath() (string, error) {
	return GetGobinPath("golangci-lint")
}

// GetGobinPath returns the full path for a filename in $GOPATH/bin.
func GetGobinPath(name string) (string, error) {
	gobinPath, err := GetGobin()
	if err != nil {
		return "", errors.Wrapf(err, "cannot get $GOPATH/bin path")
	}

	return path.Join(gobinPath, name), nil
}

// GetGobin returns the full path of $GOPATH/bin, creating the directory as necessary.
func GetGobin() (string, error) {
	// Get GOPATH.
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Prepare destination path.
	binPath := path.Join(gopath, "bin")
	if err := os.MkdirAll(binPath, 0755); err != nil {
		return "", errors.Wrapf(err, "cannot mkdir %v", binPath)
	}

	return binPath, nil
}
