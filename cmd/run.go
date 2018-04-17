package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/springload/ssm-parent/ssm"

	"github.com/apex/log"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

var expand bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run command",
	Short: "Runs the specified command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		var cmdArgs []string

		megamap := make(map[string]interface{})
		localNames := names
		localPaths := paths
		if expand {
			localNames = expandArgs(names)
			localPaths = expandArgs(paths)
		}
		parameters, err := ssm.GetParameters(localNames, localPaths, strict, recursive)
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
		if expand {
			cmdArgs = expandArgs(args[1:])
		} else {
			cmdArgs = args[1:]
		}

		cmd := exec.Command(command, cmdArgs...)
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

// expandArgs leverages on shell and echo to expand
// possible args mainly env vars.
// taken from https://github.com/abiosoft/parent
func expandArgs(args []string) []string {
	var expanded []string
	for _, arg := range args {
		e, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo %s", arg)).Output()
		// error is not expected.
		// in the rare case that this errors
		// the original arg is still used.
		if err == nil {
			arg = strings.TrimSpace(string(e))
		}
		expanded = append(expanded, arg)
	}
	return expanded
}

func init() {
	runCmd.Flags().BoolVarP(&expand, "expand", "e", false, "Expand arguments using /bin/sh")
	rootCmd.AddCommand(runCmd)
}
