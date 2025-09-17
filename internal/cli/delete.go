package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/onggizam/mcc/internal/utils"
)

var delYes bool

func init() {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a stored kubeconfig",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !delYes && isTTY() {
				fmt.Printf("Delete '%s'? [y/N]: ", name)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				resp := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if resp != "y" && resp != "yes" {
					fmt.Println("aborted")
					return nil
				}
			}
			return utils.Delete(utils.ResolveMultiDir(multiDirFlag), name)
		},
	}
	cmd.Flags().BoolVarP(&delYes, "yes", "y", false, "do not prompt for confirmation")
	rootCmd.AddCommand(cmd)
}

func isTTY() bool {
	// best-effort: check if stdin is a terminal
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
