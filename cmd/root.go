package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	paths     []string
	names     []string
	recursive bool
	strict    bool
	Version   = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parent",
	Short: "Docker entrypoint that get parameters from AWS SSM Parameter Store",
	Long: `Parent is a docker entrypoint.
	
It gets specified parameters (possibly secret) from AWS SSM Parameter Store,
then exports them to the underlying process.
`,
	Version: Version,
	//	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&paths, "path", "p", []string{}, "Path to a SSM parameter. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&names, "name", "n", []string{}, "Name of the SSM parameter to retrieve. Can be specified multiple times.")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Walk through the provided SSM paths recursively.")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Strict mode. Fail if found less parameters than number of names.")
}
