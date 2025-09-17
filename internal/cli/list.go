package cli

import (
	"github.com/spf13/cobra"
	"github.com/onggizam/mcc/internal/utils"
)

var listJSON bool

func init() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stored kubeconfigs",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := utils.List(utils.ResolveMultiDir(multiDirFlag), utils.ResolveKubeconfig(""))
			if err != nil {
				return err
			}
			return utils.RenderList(items, listJSON)
		},
	}
	cmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(cmd)
}
