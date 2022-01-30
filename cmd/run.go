package cmd

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run command",
	Short: "Runs the specified command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		var cmdArgs []string

		parameters, err := ssm.GetParameters(
			viper.GetStringSlice("names"),
			viper.GetStringSlice("paths"),
			viper.GetStringSlice("plain-names"),
			viper.GetStringSlice("plain-paths"),
			transformationsList,
			viper.GetBool("expand"),
			viper.GetBool("strict"),
			viper.GetBool("recursive"),
			viper.GetBool("expand-names"),
			viper.GetBool("expand-paths"),
			viper.GetStringSlice("expand-values"),
		)
		if err != nil {
			log.WithError(err).Fatal("Can't get parameters")
		}
		for key, value := range parameters {
			os.Setenv(key, value)
		}

		command, err := exec.LookPath(args[0])
		ctx := log.WithFields(log.Fields{"command": command})
		if err != nil {
			ctx.WithError(err).Fatal("Cant find the command")
		}
		cmdArgs = append(cmdArgs, args[:1]...)

		c := make(chan os.Signal, 1)
		signal.Notify(c)
		if viper.GetBool("expand") {
			cmdArgs = append(cmdArgs, ssm.ExpandArgs(args[1:])...)
		} else {
			cmdArgs = append(cmdArgs, args[1:]...)
		}
		if err := syscall.Exec(command, cmdArgs, os.Environ()); err != nil {
			ctx.WithError(err).Fatal("Can't run the command")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
