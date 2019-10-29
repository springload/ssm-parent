package cmd

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run command",
	Short: "Runs the specified command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		var cmdArgs []string

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
				value = ssm.ExpandValue(value)
			}
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
		if expand {
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
