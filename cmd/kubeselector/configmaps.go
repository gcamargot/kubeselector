package kubeselector

import (
	"context"
	"fmt"
	"time"

	"github.com/gcamargot/kubeselector/pkg/k8s"
	"github.com/gcamargot/kubeselector/pkg/ui"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	deleteConfigMap   bool
	describeConfigMap bool
)

var configMapsCmd = &cobra.Command{
	Use:   "configmaps",
	Short: "List and manage configmaps",
	Long:  `List configmaps in the current namespace and perform actions like delete or describe`,
	RunE:  runConfigMaps,
}

func init() {
	rootCmd.AddCommand(configMapsCmd)
	configMapsCmd.Flags().BoolVar(&deleteConfigMap, "delete", false, "Delete the selected configmap")
	configMapsCmd.Flags().BoolVar(&describeConfigMap, "describe", false, "Describe the selected configmap")
}

func runConfigMaps(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Kubernetes client
	client, err := k8s.NewClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	fmt.Printf("Using namespace: %s\n", client.GetNamespace())

	// Get configmaps
	configMaps, err := client.GetConfigMaps(ctx)
	if err != nil {
		return fmt.Errorf("failed to get configmaps: %w", err)
	}

	if len(configMaps) == 0 {
		fmt.Println("No configmaps found in namespace")
		return nil
	}

	// Filter configmaps if filter is specified
	configMapNames := make([]string, 0, len(configMaps))
	descriptions := make([]string, 0, len(configMaps))
	configMapMap := make(map[string]corev1.ConfigMap)

	for _, configMap := range configMaps {
		configMapMap[configMap.Name] = configMap
		configMapNames = append(configMapNames, configMap.Name)
		age := time.Since(configMap.CreationTimestamp.Time).Round(time.Second)
		dataCount := len(configMap.Data)
		descriptions = append(descriptions, fmt.Sprintf("Keys: %d, Age: %s", dataCount, age))
	}

	if filter != "" {
		configMapNames = ui.FilterItems(configMapNames, filter)
		if len(configMapNames) == 0 {
			fmt.Printf("No configmaps found matching filter '%s'\n", filter)
			return nil
		}
	}

	// Select a configmap
	selectedConfigMap, err := ui.SelectFromList("Select a configmap:", configMapNames, descriptions)
	if err != nil {
		return fmt.Errorf("failed to select configmap: %w", err)
	}

	configMap := configMapMap[selectedConfigMap]

	// Perform action based on flags
	if deleteConfigMap {
		fmt.Printf("Deleting configmap '%s'...\n", selectedConfigMap)
		if err := client.DeleteConfigMap(ctx, selectedConfigMap); err != nil {
			return err
		}
		fmt.Printf("ConfigMap '%s' deleted successfully\n", selectedConfigMap)
	} else if describeConfigMap {
		fmt.Printf("\nConfigMap: %s\n", configMap.Name)
		fmt.Printf("Namespace: %s\n", configMap.Namespace)
		fmt.Printf("Created: %s\n", configMap.CreationTimestamp.Time)
		fmt.Printf("\nData keys:\n")
		for key := range configMap.Data {
			fmt.Printf("  - %s\n", key)
		}
		if len(configMap.Labels) > 0 {
			fmt.Printf("\nLabels:\n")
			for key, value := range configMap.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	} else {
		// Just show basic info
		fmt.Printf("\nSelected configmap: %s\n", selectedConfigMap)
		fmt.Printf("Keys: %d\n", len(configMap.Data))
		fmt.Printf("Use --delete or --describe flags to perform actions\n")
	}

	return nil
}
