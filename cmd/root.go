package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

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

		runCmd(pkgName, name, style)
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
		name = *pkgName
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

	f.WriteString(mainGo)

	return nil
}

func createStandardLayoutDirs() error {
	var dirs [12]string = [12]string{
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
