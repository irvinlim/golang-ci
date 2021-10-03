package golangcilint

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
)

const (
	installerURL = "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"
)

// DownloadAllGolangCILintVersions will download all known versions of golangci-lint to dest.
// The downloaded binaries will be named as `golangci-lint-$version`, e.g. `golangci-lint-v1.31.0`.
// If minVersion is not empty, any versions equal or lower to minVersion will not be installed.
// Returns list of versions that were installed, or any error encountered if any.
func DownloadAllGolangCILintVersions(dest string, minVersion string) ([]string, error) {
	dir, err := ioutil.TempDir("", "golangci-lint-download")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create temp directory")
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// List all versions.
	versions, err := GetGolangCILintVersions()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot fetch golangci-lint versions")
	}
	log.Printf("[golangci-lint] Found %v versions of golangci-lint.\n", len(versions))

	// Fetch install script.
	installer := path.Join(dir, "install.sh")
	if err := FetchGolangCILintInstaller(installer); err != nil {
		return nil, errors.Wrapf(err, "cannot fetch golangci-lint installer")
	}

	// Run install script for each version.
	installedVersions := make([]string, 0, len(versions))
	for _, version := range versions {
		// Check min version.
		if minVersion != "" {
			if semver.Compare(version, minVersion) < 0 {
				log.Printf("[golangci-lint] Skipping install of %v < minVersion (%v).\n", version, minVersion)
				continue
			}
		}

		// Install specific version to destPath.
		destPath := path.Join(dest, GetNameForVersion(version))
		if err := InstallGolangCILintVersion(version, destPath, installer, dir); err != nil {
			return nil, errors.Wrapf(err, "cannot install %v to %v", version, destPath)
		}

		installedVersions = append(installedVersions, version)
		log.Printf("[golangci-lint] Installed %v to %v.\n", version, dest)
	}

	return installedVersions, nil
}

// InstallGolangCILintVersion will run the installer script against a specific version, saving it to dest.
// To install latest version, specify "latest" as version.
// If installer is empty, it will download to a temporary path and be removed afterwards.
// If tempDir is empty, a temporary directory will be created and be removed afterwards.
func InstallGolangCILintVersion(version, dest, installer, tempDir string) error {
	// Prepare temp directory if not specified.
	if tempDir == "" {
		newTempDir, err := ioutil.TempDir("", version)
		if err != nil {
			return errors.Wrapf(err, "cannot create temp dir")
		}
		tempDir = newTempDir
	}

	// Check if installer already exists.
	var isInstallerExist bool
	if installer == "" {
		installer = path.Join(tempDir, "install.sh")
	}
	if stat, err := os.Stat(installer); err == nil && !stat.IsDir() {
		isInstallerExist = true
	}

	// Download installer if not already present.
	if !isInstallerExist {
		if err := FetchGolangCILintInstaller(installer); err != nil {
			return errors.Wrapf(err, "cannot fetch golangci-lint installer")
		}
		defer func() {
			_ = os.Remove(installer)
			log.Printf("[golangci-lint] Removed %v.\n", installer)
		}()
		log.Printf("[golangci-lint] Downloaded installer.sh to %v.\n", installer)
	}

	// Run installer script to download specific version.
	args := []string{"-b", tempDir}
	if version != "latest" {
		args = append(args, version)
	}
	cmd := exec.Command(installer, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error while running installer for %v", version)
	}

	// Move to dest path.
	if err := os.Rename(path.Join(tempDir, "golangci-lint"), dest); err != nil {
		return errors.Wrapf(err, "cannot move to %v", dest)
	}

	return nil
}

// FetchGolangCILintInstaller will download the golangci-lint installer to a temporary file location.
func FetchGolangCILintInstaller(scriptPath string) error {
	resp, err := http.Get(installerURL)
	if err != nil {
		return errors.Wrapf(err, "cannot fetch install script")
	}
	defer func() { _ = resp.Body.Close() }()
	scriptBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "cannot read install script")
	}
	if err := ioutil.WriteFile(scriptPath, scriptBytes, 0644); err != nil {
		return errors.Wrapf(err, "cannot create install script")
	}
	if err := os.Chmod(scriptPath, 0744); err != nil {
		return errors.Wrapf(err, "cannot chmod install script")
	}
	return nil
}
