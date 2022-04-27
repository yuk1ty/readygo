package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "readygo",
	Short: "The easiest way to get started with Go.",
	Long:  `readygo is a tiny CLI tool for creating basic Go project.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgName, err := parsePackageName(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := parseDirectoryName(cmd, pkgName)
		if err != nil {
			fmt.Println(err)
			return
		}

		layout, err := parseLayout(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = runCmd(pkgName, name, layout)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func parsePackageName(cmd *cobra.Command) (*string, error) {
	pkgName, err := cmd.Flags().GetString("pkg-name")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if pkgName == "" {
		return nil, errors.New("[ERROR] Package name should be set! Please use `--pkg-name` option. For more details, please hit --help option")
	}
	return &pkgName, nil
}

func parseDirectoryName(cmd *cobra.Command, pkgName *string) (*string, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}
	if name == "" {
		// e.g. If passed `github.com/yuk1ty/readygo` then extract `readygo`
		splitted := strings.Split(*pkgName, "/")
		name = splitted[len(splitted)-1]
	}
	return &name, nil
}

const (
	Default  = "default"
	Standard = "standard"
)

func parseLayout(cmd *cobra.Command) (*string, error) {
	layout, err := cmd.Flags().GetString("layout")
	if err != nil {
		return nil, err
	}

	if layout != Default && layout != Standard {
		return nil, fmt.Errorf("[ERROR] Layout name should be `%s` or `%s`", Default, Standard)
	}

	return &layout, nil
}

func runCmd(pkgName *string, dirName *string, layout *string) error {
	err := os.Mkdir(*dirName, 0777)
	if err != nil {
		return err
	}

	err = os.Chdir(*dirName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "init", *pkgName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	git := exec.Command("git", "init")
	err = git.Run()
	if err != nil {
		return err
	}

	err = createGitIgnore()
	if err != nil {
		return err
	}

	err = createMainGo()
	if err != nil {
		return err
	}

	if *layout == Standard {
		err = createStandardLayoutDirs()
		if err != nil {
			return err
		}
	}

	return nil
}

func createMainGo() error {
	mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}`

	f, err := os.Create("main.go")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(mainGo)
	if err != nil {
		return err
	}

	return nil
}

func createGitIgnore() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Try to get `gibo`'s boilerplate. If it exists, keep reading the file.
	// If not, try to get remote boilerplate on GitHub.
	var gitignore string
	giboBoilerplate, err := ioutil.ReadFile(filepath.Join(homedir, ".gitignore-boilerplates", "Go.gitignore"))
	if err != nil {
		resp, err := http.Get("https://raw.github.com/github/gitignore/994f99fc353f523dfe5633067805a1dd4a53040f/Go.gitignore")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		gitignore = string(body)
	} else {
		gitignore = string(giboBoilerplate)
	}

	f, err := os.Create(".gitignore")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(gitignore)
	if err != nil {
		return err
	}
	return nil
}

func createStandardLayoutDirs() error {
	var dirs [3]string = [3]string{
		"cmd",
		"internal",
		"pkg",
	}

	for _, name := range dirs {
		err := os.Mkdir(name, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("pkg-name", "p", "", "Define your package name. This is used for go mod init [module path].")
	rootCmd.Flags().StringP("name", "n", "", "Define the directory name of your project. This can be omitted. If you do so, name will be extracted from package name.")
	rootCmd.Flags().StringP("layout", "l", "default", "Define your project layout. You can choose `default` or `standard`. If you omit this option, the value becomes `default`.")
}
