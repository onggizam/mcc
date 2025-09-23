package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Show current active cluster info and cluster-wide pods summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		overrides := &clientcmd.ConfigOverrides{}
		cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

		restCfg, err := cc.ClientConfig()
		if err != nil {
			return fmt.Errorf("failed to load kubeconfig: %w", err)
		}

		raw, err := cc.RawConfig()
		if err != nil {
			return fmt.Errorf("failed to read raw kubeconfig: %w", err)
		}
		currentCtx := raw.CurrentContext
		if currentCtx == "" {
			fmt.Println("No current context is set in kubeconfig")
			return nil
		}

		ctxMeta := raw.Contexts[currentCtx]
		if ctxMeta == nil {
			fmt.Printf("Context %q not found in kubeconfig\n", currentCtx)
			return nil
		}
		clusterName := ctxMeta.Cluster
		userName := ctxMeta.AuthInfo
		ns := ctxMeta.Namespace
		if ns == "" {
			if n, _, err := cc.Namespace(); err == nil && n != "" {
				ns = n
			} else {
				ns = "default"
			}
		}

		server := ""
		if cl := raw.Clusters[clusterName]; cl != nil {
			server = cl.Server
		}

		clientset, err := kubernetes.NewForConfig(restCfg)
		if err != nil {
			return fmt.Errorf("failed to create k8s client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		podList, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list pods: %w", err)
		}

		total := len(podList.Items)
		var running, pending, succeeded, failed, unknown int
		nsCount := make(map[string]int)

		for _, p := range podList.Items {
			nsCount[p.Namespace]++
			switch p.Status.Phase {
			case "Running":
				running++
			case "Pending":
				pending++
			case "Succeeded":
				succeeded++
			case "Failed":
				failed++
			default:
				unknown++
			}
		}

		fmt.Printf("Current context : %s\n", currentCtx)
		fmt.Printf("Cluster name    : %s\n", clusterName)
		fmt.Printf("Server          : %s\n", server)
		fmt.Printf("User            : %s\n", userName)
		fmt.Printf("Namespace       : %s\n", ns)
		fmt.Println()

		fmt.Println("Cluster-wide Pods Summary (all namespaces)")
		fmt.Printf("Total      : %d\n", total)
		fmt.Printf("Running    : %d\n", running)
		fmt.Printf("Pending    : %d\n", pending)
		fmt.Printf("Succeeded  : %d\n", succeeded)
		fmt.Printf("Failed     : %d\n", failed)
		fmt.Printf("Unknown    : %d\n", unknown)

		if os.Getenv("MCC_NOW_SHOW_NS") == "1" {
			fmt.Println()
			fmt.Println("Pods per namespace:")
			for k, v := range nsCount {
				fmt.Printf("- %s: %d\n", k, v)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)
}
