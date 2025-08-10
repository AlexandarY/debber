package main

import (
	"fmt"
	"os"

	"github.com/AlexandarY/debber/internal/debber"
	"github.com/spf13/cobra"
)

var debianFile string

var rootCmd = &cobra.Command{
	Use:   "debber",
	Short: "Simple debian package generator",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new debian.toml template",
	Long:  `Generate a new debian.toml file with all required fields for a package build`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "Generating %s\n", debianFile)
		err := debber.CreateNewDebFile(debianFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a debian/ directory from debian.toml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "Creating the debian/ dir and content\n")
		content, err := debber.ParseFile(debianFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		fmt.Fprintln(os.Stdout, content)
		debDir, err := debber.CreateDebianDirectory("tmp/", content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		err = debDir.CreateControl()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&debianFile, "name", "n", "debian.toml", "Specify the name of the debian config file")

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
	}
}
