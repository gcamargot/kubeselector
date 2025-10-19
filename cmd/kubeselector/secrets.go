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
	deleteSecret   bool
	describeSecret bool
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List and manage secrets",
	Long:  `List secrets in the current namespace and perform actions like delete or describe`,
	RunE:  runSecrets,
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.Flags().BoolVar(&deleteSecret, "delete", false, "Delete the selected secret")
	secretsCmd.Flags().BoolVar(&describeSecret, "describe", false, "Describe the selected secret")
}

func runSecrets(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Kubernetes client
	client, err := k8s.NewClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	fmt.Printf("Using namespace: %s\n", client.GetNamespace())

	// Get secrets
	secrets, err := client.GetSecrets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found in namespace")
		return nil
	}

	// Filter secrets if filter is specified
	secretNames := make([]string, 0, len(secrets))
	descriptions := make([]string, 0, len(secrets))
	secretMap := make(map[string]corev1.Secret)

	for _, secret := range secrets {
		secretMap[secret.Name] = secret
		secretNames = append(secretNames, secret.Name)
		age := time.Since(secret.CreationTimestamp.Time).Round(time.Second)
		descriptions = append(descriptions, fmt.Sprintf("Type: %s, Age: %s", secret.Type, age))
	}

	if filter != "" {
		secretNames = ui.FilterItems(secretNames, filter)
		if len(secretNames) == 0 {
			fmt.Printf("No secrets found matching filter '%s'\n", filter)
			return nil
		}
	}

	// Select a secret
	selectedSecret, err := ui.SelectFromList("Select a secret:", secretNames, descriptions)
	if err != nil {
		return fmt.Errorf("failed to select secret: %w", err)
	}

	secret := secretMap[selectedSecret]

	// Perform action based on flags
	if deleteSecret {
		fmt.Printf("Deleting secret '%s'...\n", selectedSecret)
		if err := client.DeleteSecret(ctx, selectedSecret); err != nil {
			return err
		}
		fmt.Printf("Secret '%s' deleted successfully\n", selectedSecret)
	} else if describeSecret {
		fmt.Printf("\nSecret: %s\n", secret.Name)
		fmt.Printf("Namespace: %s\n", secret.Namespace)
		fmt.Printf("Type: %s\n", secret.Type)
		fmt.Printf("Created: %s\n", secret.CreationTimestamp.Time)
		fmt.Printf("\nData keys:\n")
		for key := range secret.Data {
			fmt.Printf("  - %s\n", key)
		}
		if len(secret.Labels) > 0 {
			fmt.Printf("\nLabels:\n")
			for key, value := range secret.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	} else {
		// Just show basic info
		fmt.Printf("\nSelected secret: %s\n", selectedSecret)
		fmt.Printf("Type: %s\n", secret.Type)
		fmt.Printf("Use --delete or --describe flags to perform actions\n")
	}

	return nil
}
