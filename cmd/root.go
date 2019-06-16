package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	paths      []string
	names      []string
	plainPaths []string
	plainNames []string
	recursive  bool
	strict     bool
	expand     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssm-parent",
	Short: "Docker entrypoint that get parameters from AWS SSM Parameter Store",
	Long: `SSM-Parent is a docker entrypoint.
	
It gets specified parameters (possibly secret) from AWS SSM Parameter Store,
then exports them to the underlying process. Or creates a .env file to be consumed by an application.

It reads parameters in the following order: path->name->plain-path->plain-name.
So that every rightmost parameter overrides the previous one.
`,
}

// Execute is the entrypoint for cmd/ module
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&paths, "path", "p", []string{}, "Path to a SSM parameter. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&names, "name", "n", []string{}, "Name of the SSM parameter to retrieve. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&plainPaths, "plain-path", "", []string{}, "Path to a SSM parameter. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&plainNames, "plain-name", "", []string{}, "Name of the SSM parameter to retrieve. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Walk through the provided SSM paths recursively.")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Strict mode. Fail if found less parameters than number of names.")
	rootCmd.PersistentFlags().BoolVarP(&expand, "expand", "e", false, "Expand arguments and values using shell-style syntax")
}
