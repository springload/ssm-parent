package cmd

import (
	"github.com/apex/log"
	"github.com/imdario/mergo"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/springload/ssm-parent/ssm"
)

// dotenvCmd represents the dotenv command
var dotenvCmd = &cobra.Command{
	Use:   "dotenv <filename>",
	Short: "Writes dotenv file",
	Long:  `Gathers parameters from SSM Parameter store, writes .env file and exits`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		megamap := make(map[string]string)
		parameters, err := ssm.GetParameters(names, paths, plainNames, plainPaths, expand, strict, recursive)
		if err != nil {
			log.WithError(err).Fatal("Can't get parameters")
		}
		for _, parameter := range parameters {
			err = mergo.Merge(&megamap, &parameter, mergo.WithOverride)
			if err != nil {
				log.WithError(err).Fatal("Can't merge maps")
			}
		}
		for key, value := range megamap {
			if expand {
				megamap[key] = ssm.ExpandValue(value)
			}
		}

		err = godotenv.Write(megamap, args[0])
		if err != nil {
			log.WithError(err).Fatal("Can't write the dotenv file")
		} else {
			log.WithFields(log.Fields{"filename": args[0]}).Info("Wrote the .env file")

		}
	},
}

func init() {
	rootCmd.AddCommand(dotenvCmd)

}
