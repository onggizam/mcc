package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/onggizam/mcc/internal/utils"
)

var chBackup bool
var chKubeconfig string

func init() {
	cmd := &cobra.Command{
		Use:   "ch <name>",
		Short: "Switch ~/.kube/config to a stored entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := utils.Switch(
				utils.ResolveMultiDir(multiDirFlag),
				name,
				utils.ResolveKubeconfig(chKubeconfig),
				chBackup,
			); err != nil {
				return err
			}
			if chBackup {
				fmt.Printf("Switched cluster to \"%s\" (backup created)\n", name)
			} else {
				fmt.Printf("Switched cluster to \"%s\"\n", name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&chBackup, "backup", false, "create a backup of the current kubeconfig before switching")
	cmd.Flags().StringVar(&chKubeconfig, "kubeconfig", "", "target kubeconfig (default: ~/.kube/config)")

	rootCmd.AddCommand(cmd)
}
