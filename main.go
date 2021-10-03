package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/irvinlim/golang-ci/pkg/cmd/lint"
	lintInstall "github.com/irvinlim/golang-ci/pkg/cmd/lint-install"
	lintInstallAll "github.com/irvinlim/golang-ci/pkg/cmd/lint-install-all"
)

func main() {
	Run()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var rootCmd = &cobra.Command{
	Use:   "golang-ci",
	Short: "CI tool for Golang",
}

func init() {
	rootCmd.AddCommand(lint.Command)
	rootCmd.AddCommand(lintInstall.Command)
	rootCmd.AddCommand(lintInstallAll.Command)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error encountered: %v\n", err)
		if err, ok := err.(stackTracer); ok {
			log.Println("Traceback:")
			log.Printf("%+v\n", err.StackTrace())
		}
		os.Exit(1)
	}
}
