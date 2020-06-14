package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const appLabel = "kubectl ingress-refresh"

var version string

// SetVersion set the application version for consumption in the output of the command.
func SetVersion(v string) {
	version = v
}

type refreshCmd struct {
	out io.Writer
}

func newRefreshCmd(out io.Writer) *cobra.Command {
	refresh := &refreshCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "penis",
		Short: "print the version number and exit",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return refresh.run()
		},
	}
	return cmd
}

func (v *refreshCmd) run() error {
	_, err := fmt.Fprintf(v.out, "%s %s\n", appLabel, version)
	if err != nil {
		return err
	}
	return nil
}
