package cmd

import (
	"os"

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

		// we don't want to use godotenv as it creates files with too open permissions
		content, err := godotenv.Marshal(megamap)
		if err != nil {
			log.WithError(err).Fatal("Can't marshal the env to a string")
		}

		file, err := os.OpenFile(args[0], os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.WithError(err).Fatal("Can't create the file")
		}

		_, err = file.WriteString(content)
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
