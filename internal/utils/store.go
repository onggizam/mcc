package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"gopkg.in/yaml.v3"
)

type Item struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	Server    string    `json:"server"`
	Active    bool      `json:"active"`
}

func Add(multiDir, name, src string, force bool) error {
	src = expandPath(src)
	if err := validateKubeconfigFile(src); err != nil {
		return fmt.Errorf("invalid kubeconfig: %w", err)
	}
	dst := filepath.Join(multiDir, name)
	if !force {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("entry '%s' already exists (use --force to overwrite)", name)
		}
	}
	if err := copyFile(src, dst, 0o600); err != nil {
		return err
	}
	return nil
}

func Delete(multiDir, name string) error {
	path := filepath.Join(multiDir, name)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("entry '%s' not found", name)
		}
		return err
	}
	return nil
}

func extractServer(b []byte) string {
	var m map[string]any
	if err := yaml.Unmarshal(b, &m); err != nil {
		return ""
	}
	clusters, ok := m["clusters"].([]any)
	if !ok || len(clusters) == 0 {
		return ""
	}
	first, ok := clusters[0].(map[string]any)
	if !ok {
		return ""
	}
	cluster, ok := first["cluster"].(map[string]any)
	if !ok {
		return ""
	}
	if s, ok := cluster["server"].(string); ok {
		return s
	}
	return ""
}

func List(multiDir string, targetKubeconfig string) ([]Item, error) {
	entries, err := os.ReadDir(multiDir)
	if err != nil {
		return nil, err
	}

	activeHash := ""
	if b, err := os.ReadFile(targetKubeconfig); err == nil {
		activeHash = fmt.Sprintf("%x", sha256.Sum256(b))
	}

	items := make([]Item, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(multiDir, e.Name())
		b, _ := os.ReadFile(fullPath)
		server := extractServer(b)

		hash := fmt.Sprintf("%x", sha256.Sum256(b))

		items = append(items, Item{
			Name:      e.Name(),
			Path:      fullPath,
			CreatedAt: fi.ModTime(),
			Server:    server,
			Active:    hash == activeHash,
		})
	}
	return items, nil
}

func RenderList(items []Item, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(items)
	}

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(
			tw.Rendition{Symbols: tw.NewSymbols(tw.StyleASCII)},
		)),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.Off},
				Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
			},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)

	table.Header([]string{"NO", "NAME", "CREATED", "SERVER", "ACTIVE"})

	for i, it := range items {
		active := ""
		if it.Active {
			active = "Now!"
		}
		row := []string{
			fmt.Sprintf("%d", i+1),
			it.Name,
			it.CreatedAt.Format("2006-01-02 15:04:05"),
			it.Server,
			active,
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

func copyFile(src, dst string, perm os.FileMode) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	tmp := dst + ".tmp"
	df, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(df, sf); err != nil {
		df.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := df.Sync(); err != nil {
		df.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := df.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}
