package cli

import (
	"fmt"
	"os"

	"github.com/onggizam/mcc/internal/utils"
	"github.com/onggizam/mcc/pkg/version"
	"github.com/spf13/cobra"
)

var (
	multiDirFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "mcc",
	Short: "Multi Cluster Changer for kubeconfig",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		md := utils.ResolveMultiDir(multiDirFlag)
		return os.MkdirAll(md, 0o700)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&multiDirFlag, "multi-dir", "", "multi store dir (default: ~/.kube/multi)")

	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mcc current version: %s\n", version.Version)
	},
}

// Execute entrypoint
func Execute() { _ = rootCmd.Execute() }
