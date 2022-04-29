package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func execute(t *testing.T, c *cobra.Command, tmpDirPath string, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := os.Chdir(tmpDirPath)
	if err != nil {
		return "", err
	}

	err = c.Execute()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func TestPackageNameCommand(t *testing.T) {
	pkgName := "github.com/yuk1ty/example1"
	dirName := "example1"

	tmpDirPath := os.TempDir()

	_, err := execute(t, rootCmd, tmpDirPath, "-p", pkgName)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("exists a directory named `example`", func(t *testing.T) {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", tmpDirPath, "example1")); err != nil {
			t.Fatal(err)
		}
	})

	t.Cleanup(func() {
		os.RemoveAll(fmt.Sprintf("%s/%s", tmpDirPath, dirName))
	})
}

func TestNameCommand(t *testing.T) {
	pkgName := "github.com/yuk1ty/example2"
	dirName := "another"

	tmpDirPath := os.TempDir()

	_, err := execute(t, rootCmd, tmpDirPath, "-p", pkgName, "-n", dirName)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("exists a directory named `another`", func(t *testing.T) {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", tmpDirPath, "another")); err != nil {
			t.Fatal(err)
		}
	})

	t.Cleanup(func() {
		os.RemoveAll(fmt.Sprintf("%s/%s", tmpDirPath, dirName))
	})
}
