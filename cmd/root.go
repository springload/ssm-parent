package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/springload/ssm-parent/ssm/transformations"
)

var (
	config              string
	transformationsList []transformations.Transformation
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

func initSettings() {
	if config != "" {
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())
		} else {
			log.WithError(err).Fatal("Had some errors while parsing the config")
		}
	}
	// parse to an array first
	var transformationsInterfaceArray = []interface{}{}
	if err := viper.UnmarshalKey("transformations", &transformationsInterfaceArray); err != nil {
		log.WithError(err).Fatal("can't decode config")
	}
	// unmarshal to the tiny struct first to see what the action is
	for n, t := range transformationsInterfaceArray {
		var hint = struct{ Action string }{}

		if err := mapstructure.Decode(t, &hint); err != nil {
			log.WithError(err).Fatal("can't decode config")
		}
		switch hint.Action {

		case "delete":
			tr := new(transformations.DeleteTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)
		case "rename":
			tr := new(transformations.RenameTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)
		case "template":
			tr := new(transformations.TemplateTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)

		default:
			log.Warnf("Got unparsed action: %s", hint.Action)
		}
	}

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	cobra.OnInitialize(initSettings)
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path to the config file (optional). Allows to set transformations")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Turn on debug logging")
	rootCmd.PersistentFlags().BoolP("expand", "e", false, "Expand all arguments and values using shell-style syntax")
	rootCmd.PersistentFlags().BoolP("expand-names", "", false, "Expand SSM names using shell-style syntax. The '--expand' does the same, but this flag is more selective")
	rootCmd.PersistentFlags().BoolP("expand-paths", "", false, "Expand SSM paths using shell-style syntax. The '--expand' does the same, but this flag is more selective")
	rootCmd.PersistentFlags().StringSliceP("expand-values", "", []string{}, "Expand SSM values using shell-style syntax. The '--expand' does the same, but this flag is more selective. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringSliceP("path", "p", []string{}, "Path to a SSM parameter. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Name of the SSM parameter to retrieve. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringSliceP("plain-path", "", []string{}, "Path to a SSM parameter. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringSliceP("plain-name", "", []string{}, "Name of the SSM parameter to retrieve. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().BoolP("recursive", "r", false, "Walk through the provided SSM paths recursively.")
	rootCmd.PersistentFlags().BoolP("strict", "s", false, "Strict mode. Fail if found less parameters than number of names.")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("expand", rootCmd.PersistentFlags().Lookup("expand"))
	viper.BindPFlag("expand-names", rootCmd.PersistentFlags().Lookup("expand-names"))
	viper.BindPFlag("expand-paths", rootCmd.PersistentFlags().Lookup("expand-paths"))
	viper.BindPFlag("expand-values", rootCmd.PersistentFlags().Lookup("expand-values"))
	viper.BindPFlag("paths", rootCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("names", rootCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("plain-paths", rootCmd.PersistentFlags().Lookup("plain-path"))
	viper.BindPFlag("plain-names", rootCmd.PersistentFlags().Lookup("plain-name"))
	viper.BindPFlag("recursive", rootCmd.PersistentFlags().Lookup("recursive"))
	viper.BindPFlag("strict", rootCmd.PersistentFlags().Lookup("strict"))
}
