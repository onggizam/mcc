package utils

import (
	"path/filepath"
)

// ResolveMultiDir returns the multi directory. If user provided empty, default to ~/.kube/multi
func ResolveMultiDir(userProvided string) string {
	if userProvided != "" {
		return expandPath(userProvided)
	}
	return expandPath("~/.kube/multi")
}

// ResolveKubeconfig returns the target kubeconfig path (default: ~/.kube/config)
func ResolveKubeconfig(userProvided string) string {
	if userProvided != "" {
		return expandPath(userProvided)
	}
	return expandPath("~/.kube/config")
}

func expandPath(p string) string {
	if len(p) > 1 && p[:2] == "~/" {
		h, _ := osUserHomeDir()
		return filepath.Join(h, p[2:])
	}
	return p
}
