package cmd

import (
	"encoding/json"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the specified parameters.",
	Run: func(cmd *cobra.Command, args []string) {

		megamap := make(map[string]string)
		parameters, err := ssm.GetParameters(names, paths, expand, strict, recursive)
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
		marshalled, err := json.MarshalIndent(megamap, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("Can't marshal json")
		}
		log.Info(string(marshalled))
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
