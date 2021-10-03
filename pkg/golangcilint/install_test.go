package golangcilint_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/irvinlim/golang-ci/pkg/golangcilint"
)

func TestInstallGolangCILintVersion(t *testing.T) {
	version := "v1.31.0"

	destDir, err := ioutil.TempDir("", t.Name()+"-dest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(destDir) }()

	destPath := path.Join(destDir, golangcilint.GetNameForVersion(version))

	// Download version.
	if err := golangcilint.InstallGolangCILintVersion(version, destPath, "", ""); err != nil {
		t.Errorf("cannot install: %v", err)
		return
	}
	if _, err := os.Stat(destPath); err != nil {
		t.Errorf("cannot stat %v", destPath)
		return
	}

	// Execute file, check that it works.
	cmd := exec.Command(destPath, "--version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = stdout.Close() }()
	var stdoutString string
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			t.Errorf("cannot read stdout: %v", err)
		}
		stdoutString = string(bytes)
	}()

	if err := cmd.Run(); err != nil {
		t.Errorf("cannot run %v --version: %v", destPath, err)
		return
	}
	wg.Wait()

	// Check if contains version string, without leading "v"
	if !strings.Contains(stdoutString, version[1:]) {
		t.Errorf(`expected "%v" to contain "%v"`, stdoutString, version)
	}
}
