package cmd

import (
	"encoding/json"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the specified parameters.",
	Run: func(cmd *cobra.Command, args []string) {
		parameters, err := ssm.GetParameters(
			viper.GetStringSlice("names"),
			viper.GetStringSlice("paths"),
			viper.GetStringSlice("plain-names"),
			viper.GetStringSlice("plain-paths"),
			transformationsList,
			viper.GetBool("expand"),
			viper.GetBool("strict"),
			viper.GetBool("recursive"),
		)
		if err != nil {
			log.WithError(err).Fatal("Can't marshal json")
		}
		marshalled, err := json.MarshalIndent(parameters, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("Can't marshal json")
		}
		log.Info(string(marshalled))
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
