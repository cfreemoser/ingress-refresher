package main

import (
	"os"

	"irefresher/cmd"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var version = "undefined"

func main() {
	cmd.SetVersion(version)

	refreshCmd := cmd.NewRefreshCommand(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := refreshCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
