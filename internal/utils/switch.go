package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Switch(multiDir, name, targetKubeconfig string, makeBackup bool) error {
	src := filepath.Join(multiDir, name)
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entry '%s' not found", name)
		}
		return err
	}
	// optional backup
	if makeBackup {
		if _, err := os.Stat(targetKubeconfig); err == nil {
			bak := targetKubeconfig + ".bak-" + time.Now().Format("20060102-150405")
			if err := copyFile(targetKubeconfig, bak, 0o600); err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
		}
	}
	// copy to target
	if err := copyFile(src, targetKubeconfig, 0o600); err != nil {
		return err
	}
	return nil
}
