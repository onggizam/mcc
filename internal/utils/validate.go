package utils

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func validateKubeconfigFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var m map[string]any
	if err := yaml.Unmarshal(b, &m); err != nil {
		return err
	}
	// minimal sanity checks
	if _, ok := m["clusters"]; !ok {
		return errors.New("missing 'clusters'")
	}
	if _, ok := m["contexts"]; !ok {
		return errors.New("missing 'contexts'")
	}
	if _, ok := m["users"]; !ok {
		return errors.New("missing 'users'")
	}
	// current-context is optional but nice to have; warn only
	if _, ok := m["current-context"]; !ok {
		fmt.Fprintln(os.Stderr, "warning: no 'current-context' in kubeconfig")
	}
	return nil
}
