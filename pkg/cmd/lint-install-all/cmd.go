package lint_install_all

import (
	"github.com/spf13/cobra"

	"github.com/irvinlim/golang-ci/pkg/golangcilint"
	"github.com/irvinlim/golang-ci/pkg/utils/stringutils"
)

var (
	minVersion string
)

func init() {
	Command.Flags().StringVar(&minVersion, "min-version", "",
		"Minimum version of golangci-lint that should be installed.")
}

var Command = &cobra.Command{
	Use:   "lint-install-all",
	Short: "Installs all versions of golangci-lint",
	Long: stringutils.MakeLines(
		"Installs all versions of golangci-lint, downloading to $GOPATH/bin.",
	),
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prepare $GOPATH/bin.
		binPath, err := golangcilint.GetGobin()
		if err != nil {
			return err
		}

		// Download to $GOPATH/bin.
		if _, err := golangcilint.DownloadAllGolangCILintVersions(binPath, minVersion); err != nil {
			return err
		}

		return nil
	},
}
