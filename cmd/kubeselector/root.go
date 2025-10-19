package kubeselector

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	namespace string
	filter    string
)

var rootCmd = &cobra.Command{
	Use:   "kubeselector",
	Short: "An interactive command-line tool for Kubernetes resource management",
	Long: `kubeselector is an interactive CLI tool that simplifies Kubernetes resource management.
Quickly find Pods, Secrets, ConfigMaps by name or label.
Automatically lists resources in the current namespace, or switch with a flag.
One-key actions: delete, describe, or fetch logs directly.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (uses current context namespace if not specified)")
	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "Filter resources by name")
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
