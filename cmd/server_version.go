package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	v1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type serverVersionCmd struct {
	out io.Writer
}

// NewRefreshCommand creates the command for rendering the Kubernetes server version.
func NewRefreshCommand(streams genericclioptions.IOStreams) *cobra.Command {
	helloWorldCmd := &serverVersionCmd{
		out: streams.Out,
	}

	cmd := &cobra.Command{
		Use:          "server-version",
		Short:        "Prints Kubernetes server version",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return helloWorldCmd.run()
		},
	}

	cmd.AddCommand(newRefreshCmd(streams.Out))
	return cmd
}

func (sv *serverVersionCmd) run() error {
	serverVersion, err := refreshIngress()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(sv.out, "Hello from Kubernetes server with version %s!\n", serverVersion)
	if err != nil {
		return err
	}
	return nil
}

func refreshIngress() (string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	iclient := clientset.NetworkingV1beta1().Ingresses("tlm-dev")

	ingressesList, err := iclient.List(ctx, v1.ListOptions{})
	if err != nil {
		return "", err
	}

	ingresses := ingressesList.Items
	for _, ingress := range ingresses {
		fmt.Println(ingress.Name)
		iclient.Delete(ctx, ingress.Name, v1.DeleteOptions{})
		// Debug
		penisses, err := iclient.List(ctx, v1.ListOptions{})
		if err != nil {
			return "", err
		}
		fmt.Println(len(penisses.Items))
	}

	for _, ingress := range ingresses {
		clearIngress(&ingress)
		_, err := iclient.Create(ctx, &ingress, v1.CreateOptions{})
		if err != nil {
			return "", err
		}
	}

	return "DONE", nil
}

func clearIngress(ingress *v1beta1.Ingress) {
	ingress.ResourceVersion = ""
}
