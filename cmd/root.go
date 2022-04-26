package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

		style, err := parseStyle(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = runCmd(pkgName, name, style)
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
		return nil, errors.New("[ERROR] Package name should be set! Please use `--pkg-name` option")
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

func parseStyle(cmd *cobra.Command) (*string, error) {
	style, err := cmd.Flags().GetString("style")
	if err != nil {
		return nil, err
	}

	if style != "default" && style != "standard" {
		return nil, errors.New("[ERROR] Style name should be `default` or `standard`")
	}

	return &style, nil
}

func runCmd(pkgName *string, dirName *string, style *string) error {
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

	if *style == "standard" {
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
	resp, err := http.Get("https://raw.github.com/github/gitignore/994f99fc353f523dfe5633067805a1dd4a53040f/Go.gitignore")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gitignore, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(".gitignore")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(string(gitignore))
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
	rootCmd.Flags().StringP("pkg-name", "p", "", "package")
	rootCmd.Flags().StringP("name", "n", "", "directory name")
	rootCmd.Flags().StringP("style", "s", "default", "workspace architecture")
}
