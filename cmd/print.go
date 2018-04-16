// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

		megamap := make(map[string]interface{})
		parameters, err := ssm.GetParameters(names, paths, strict, recursive)
		if err != nil {
			log.WithError(err).Fatal("Can't get parameters")
		}
		for _, parameter := range parameters {
			err = mergo.Merge(&megamap, &parameter, mergo.WithOverride)
			if err != nil {
				log.WithError(err).Fatal("Can't merge maps")
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
