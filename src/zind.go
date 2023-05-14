package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{
		Use: "zind [Command]",
	}
	rootCommand.AddCommand(InitRunCommand())
	rootCommand.AddCommand(InitChild())
	if rootCommandExecuteErr := rootCommand.Execute(); rootCommandExecuteErr != nil {
		panic(rootCommandExecuteErr)
	}
}
