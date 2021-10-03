package lint_install

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/irvinlim/golang-ci/pkg/golangcilint"
	"github.com/irvinlim/golang-ci/pkg/utils/stringutils"
)

var Command = &cobra.Command{
	Use:   "lint-install VERSION",
	Short: "Installs golangci-lint",
	Long: stringutils.MakeLines(
		"Installs a specific version of golangci-lint, downloading to $GOPATH/bin as necessary.\n",
		"Example usage:",
		"  # Installs the latest version at $GOPATH/bin/golangci-lint.",
		"  golang-ci lint-install latest\n",
		"  # Installs a specific version of golangci-lint at $GOPATH/bin/golangci-lint@v1.31.0.",
		"  golang-ci lint-install v1.31.0",
	),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		var destPath string

		// Get path to download to.
		if version == "latest" {
			dest, err := golangcilint.GetDefaultPath()
			if err != nil {
				return errors.Wrapf(err, "cannot get golangci-lint path")
			}
			destPath = dest
		} else {
			dest, err := golangcilint.GetPathForVersion(version)
			if err != nil {
				return errors.Wrapf(err, "cannot get golangci-lint path")
			}
			destPath = dest
		}

		// Install version as necessary.
		if stat, err := os.Stat(destPath); err != nil || stat.IsDir() {
			log.Printf("[lint-install] Installing %v at %v...\n", version, destPath)
			if err := golangcilint.InstallGolangCILintVersion(version, destPath, "", ""); err != nil {
				return errors.Wrapf(err, "cannot install golangci-lint at %v", destPath)
			}
			log.Printf("[lint-install] Installed golangci-lint at path: %v\n", destPath)
		} else {
			log.Printf("[lint-install] File already exists: %v\n", destPath)
		}

		return nil
	},
}
