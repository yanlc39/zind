package main

import (
	"github.com/spf13/cobra"
	"zind/lxc"
)

func InitChild() *cobra.Command {
	childCommand := &cobra.Command{
		Use: "child",
		Run: func(selfCommand *cobra.Command, args []string) {
			if createChildProcessErr := lxc.CreateChildProcess(args); createChildProcessErr != nil {
				panic(createChildProcessErr)
			}
		},
	}

	return childCommand
}
