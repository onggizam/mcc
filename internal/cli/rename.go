package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func resolveMultiDir(cmd *cobra.Command) (string, error) {
	if f := cmd.Flags().Lookup("multi-dir"); f != nil {
		if v, err := cmd.Flags().GetString("multi-dir"); err == nil && v != "" {
			return v, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home dir: %w", err)
	}
	return filepath.Join(home, ".kube", "multi"), nil
}

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a stored kubeconfig",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName := args[0]
		newName := args[1]

		multiDir, err := resolveMultiDir(cmd)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(multiDir, 0o755); err != nil {
			return fmt.Errorf("failed to ensure multi dir: %w", err)
		}

		oldPath := filepath.Join(multiDir, oldName)
		newPath := filepath.Join(multiDir, newName)

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			return fmt.Errorf("cluster config %q not found", oldName)
		}

		if _, err := os.Stat(newPath); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("cluster config %q already exists (use --force to overwrite)", newName)
			}
			if err := os.Remove(newPath); err != nil {
				return fmt.Errorf("failed to overwrite existing config %q: %v", newName, err)
			}
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to rename %q to %q: %v", oldName, newName, err)
		}

		fmt.Printf("Cluster config renamed: %s → %s\n", oldName, newName)
		return nil
	},
}

func init() {
	renameCmd.Flags().Bool("force", false, "Overwrite if the new name already exists")
	rootCmd.AddCommand(renameCmd)
}
