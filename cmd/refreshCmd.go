package cmd

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type refreshCmd struct {
	out io.Writer
	ns  string
}

// NewRefreshCommand creates the command for rendering the Kubernetes server version.
func NewRefreshCommand(streams genericclioptions.IOStreams) *cobra.Command {
	rCmd := &refreshCmd{
		out: streams.Out,
		ns:  getNamespace(genericclioptions.NewConfigFlags(true)),
	}

	cmd := &cobra.Command{
		Use:          "refresh",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("this command does not accept arguments")
			}
			return rCmd.run()
		},
	}

	return cmd
}

func (ir *refreshCmd) run() error {
	client, err := NewK8sClient()
	if err != nil {
		return err
	}

	l, err := client.LsIngress(ir.ns)

	for _, ingress := range l {
		client.DeleteIngress(&ingress)
		ingress.ResourceVersion = ""
		client.CreateIngress(&ingress)
	}
	return nil
}

// getNamespace takes a set of kubectl flag values and returns the namespace we should be operating in
func getNamespace(flags *genericclioptions.ConfigFlags) string {
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()
	if err != nil || len(namespace) == 0 {
		namespace = "default"
	}
	return namespace
}
