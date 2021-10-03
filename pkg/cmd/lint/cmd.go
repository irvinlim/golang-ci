package lint

import (
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/irvinlim/golang-ci/pkg/golangcilint"
	"github.com/irvinlim/golang-ci/pkg/utils/stringutils"
)

var (
	version string
)

func init() {
	Command.Flags().StringVar(&version, "version", "",
		"Version of golangci-lint to use. If not specified, will use latest version available.")
}

var Command = &cobra.Command{
	Use:   "lint",
	Short: "Runs linters using golangci-lint",
	Long: stringutils.MakeLines(
		"Executes a specific version of golangci-lint, downloading to $GOPATH/bin as necessary.\n",
		"Example usage:",
		"  # Run default version. Note that the `--` is necessary to pass arguments to golangci-lint.",
		"  golang-ci lint -- run -v\n",
		"  # Run a specific version of golangci-lint.",
		"  golang-ci lint --version v1.31.0 -- run -v",
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := version
		var golangCILint string

		// Check for specific version if specified.
		if version != "" {
			// Get path for specific version of golangci-lint.
			dest, err := golangcilint.GetPathForVersion(version)
			if err != nil {
				return errors.Wrapf(err, "cannot get golangci-lint path")
			}

			// Installs specific version of golangci-lint if not already present.
			if stat, err := os.Stat(dest); err != nil || stat.IsDir() {
				log.Printf("[lint] Installing %v at %v...\n", version, dest)
				if err := golangcilint.InstallGolangCILintVersion(version, dest, "", ""); err != nil {
					return errors.Wrapf(err, "cannot install golangci-lint %v", version)
				}
			}

			golangCILint = dest
		} else {
			// Assumes that default version is already installed.
			dest, err := golangcilint.GetDefaultPath()
			if err != nil {
				return errors.Wrapf(err, "cannot get golangci-lint path")
			}
			golangCILint = dest
		}

		log.Printf("[lint] Using golangci-lint at path %v.\n", golangCILint)

		// Passes arguments to golangci-lint and execute.
		golangCILintCmd := exec.Command(golangCILint, args...)
		golangCILintCmd.Stdout = os.Stdout
		golangCILintCmd.Stderr = os.Stderr
		return golangCILintCmd.Run()
	},
}
