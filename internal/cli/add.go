package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/onggizam/mcc/internal/utils"
)

var (
	addName  string
	addFile  string
	addForce bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a kubeconfig to the multi store",
		RunE: func(cmd *cobra.Command, args []string) error {
			src := addFile
			if src == "" {
				src = utils.ResolveKubeconfig("")
			}
			if err := utils.Add(utils.ResolveMultiDir(multiDirFlag), addName, src, addForce); err != nil {
				return err
			}
			fmt.Printf("Added kubeconfig \"%s\" to multi store\n", addName)
			return nil
		},
	}
	cmd.Flags().StringVarP(&addName, "name", "n", "", "name to store kubeconfig as (required)")
	cmd.Flags().StringVarP(&addFile, "file", "f", "", "source kubeconfig file (default: ~/.kube/config)")
	cmd.Flags().BoolVar(&addForce, "force", false, "overwrite if name already exists")
	_ = cmd.MarkFlagRequired("name")
	rootCmd.AddCommand(cmd)
}
