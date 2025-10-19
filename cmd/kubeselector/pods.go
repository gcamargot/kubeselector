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
	deletePod   bool
	describePod bool
	logsPod     bool
)

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List and manage pods",
	Long:  `List pods in the current namespace and perform actions like delete, describe, or view logs`,
	RunE:  runPods,
}

func init() {
	rootCmd.AddCommand(podsCmd)
	podsCmd.Flags().BoolVar(&deletePod, "delete", false, "Delete the selected pod")
	podsCmd.Flags().BoolVar(&describePod, "describe", false, "Describe the selected pod")
	podsCmd.Flags().BoolVar(&logsPod, "logs", false, "View logs for the selected pod")
}

func runPods(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Kubernetes client
	client, err := k8s.NewClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	fmt.Printf("Using namespace: %s\n", client.GetNamespace())

	// Get pods
	pods, err := client.GetPods(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	if len(pods) == 0 {
		fmt.Println("No pods found in namespace")
		return nil
	}

	// Filter pods if filter is specified
	podNames := make([]string, 0, len(pods))
	descriptions := make([]string, 0, len(pods))
	podMap := make(map[string]corev1.Pod)

	for _, pod := range pods {
		podMap[pod.Name] = pod
		podNames = append(podNames, pod.Name)
		status := string(pod.Status.Phase)
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second)
		descriptions = append(descriptions, fmt.Sprintf("Status: %s, Age: %s", status, age))
	}

	if filter != "" {
		podNames = ui.FilterItems(podNames, filter)
		if len(podNames) == 0 {
			fmt.Printf("No pods found matching filter '%s'\n", filter)
			return nil
		}
	}

	// Select a pod
	selectedPod, err := ui.SelectFromList("Select a pod:", podNames, descriptions)
	if err != nil {
		return fmt.Errorf("failed to select pod: %w", err)
	}

	pod := podMap[selectedPod]

	// Perform action based on flags
	if deletePod {
		fmt.Printf("Deleting pod '%s'...\n", selectedPod)
		if err := client.DeletePod(ctx, selectedPod); err != nil {
			return err
		}
		fmt.Printf("Pod '%s' deleted successfully\n", selectedPod)
	} else if describePod {
		fmt.Printf("\nPod: %s\n", pod.Name)
		fmt.Printf("Namespace: %s\n", pod.Namespace)
		fmt.Printf("Status: %s\n", pod.Status.Phase)
		fmt.Printf("Node: %s\n", pod.Spec.NodeName)
		fmt.Printf("IP: %s\n", pod.Status.PodIP)
		fmt.Printf("Created: %s\n", pod.CreationTimestamp.Time)
		fmt.Printf("\nContainers:\n")
		for _, container := range pod.Spec.Containers {
			fmt.Printf("  - Name: %s\n", container.Name)
			fmt.Printf("    Image: %s\n", container.Image)
		}
		fmt.Printf("\nLabels:\n")
		for key, value := range pod.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
		if len(pod.Status.Conditions) > 0 {
			fmt.Printf("\nConditions:\n")
			for _, condition := range pod.Status.Conditions {
				fmt.Printf("  %s: %s (Reason: %s)\n", condition.Type, condition.Status, condition.Reason)
			}
		}
	} else if logsPod {
		containerName := ""
		if len(pod.Spec.Containers) > 1 {
			// If multiple containers, let user select
			containerNames := make([]string, len(pod.Spec.Containers))
			for i, c := range pod.Spec.Containers {
				containerNames[i] = c.Name
			}
			containerName, err = ui.SelectFromList("Select a container:", containerNames, nil)
			if err != nil {
				return fmt.Errorf("failed to select container: %w", err)
			}
		}
		fmt.Printf("Fetching logs for pod '%s'...\n", selectedPod)
		logs, err := client.GetPodLogs(ctx, selectedPod, containerName)
		if err != nil {
			return err
		}
		fmt.Printf("\n--- Logs for pod '%s' ---\n", selectedPod)
		fmt.Println(logs)
	} else {
		// Just show basic info
		fmt.Printf("\nSelected pod: %s\n", selectedPod)
		fmt.Printf("Status: %s\n", pod.Status.Phase)
		fmt.Printf("Use --delete, --describe, or --logs flags to perform actions\n")
	}

	return nil
}
