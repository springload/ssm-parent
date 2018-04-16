package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified command",
	Run: func(cobraCmd *cobra.Command, args []string) {
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
		for key, value := range megamap {
			os.Setenv(key, fmt.Sprintf("%v", value))
		}

		command, err := exec.LookPath(args[0])
		ctx := log.WithFields(log.Fields{"command": command})
		if err != nil {
			ctx.WithError(err).Fatal("Cant find the command")
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c)

		cmd := exec.Command(command, args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Start(); err != nil {
			ctx.WithError(err).Fatal("Can't run the command")
		}

		go func() {
			for sig := range c {
				cmd.Process.Signal(sig)
			}
		}()

		if err := cmd.Wait(); err != nil {
			ctx.WithError(err).Fatal("The command exited with an error")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
