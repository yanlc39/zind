package main

import (
	"github.com/spf13/cobra"
	"os"
	"zind/lxc"
)

func InitRunCommand() *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run a Command in a new container",
		Run: func(selfCommand *cobra.Command, args []string) {
			isTTYChosen, isTTYChosenErr := selfCommand.Flags().GetBool("tty")
			if isTTYChosenErr != nil {
				panic(isTTYChosenErr)
			}
			isInteracted, isInteractedErr := selfCommand.Flags().GetBool("interactive")
			if isInteractedErr != nil {
				panic(isInteractedErr)
			}

			containerCommand := lxc.CreateParentProcess(isInteracted, isTTYChosen, args)
			if commandStartErr := containerCommand.Start(); commandStartErr != nil {
				panic(commandStartErr)
			}

			if commandWaitErr := containerCommand.Wait(); commandWaitErr != nil {
				panic(commandWaitErr)
			}

			os.Exit(-1)
		},
	}

	runCommand.Flags().BoolP("interactive", "i", false, "Keep STDIN open even if not attached")
	runCommand.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
	runCommand.Flags().BoolP("detach", "d", false, "Run container in background and print container ID")

	return runCommand
}
