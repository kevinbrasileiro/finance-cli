package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locating home directory: %w", err)
	}

	dir := filepath.Join(home, ".fin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating data directory %q: %w", dir, err)
	}
	return dir, nil
}

func nextID(model string) (int64, error) {
	dir, err := dataDir()
	if err != nil {
		return 0, err
	}
	path := filepath.Join(dir, model+".seq")

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
		data = nil
	}

	var current int64
	if len(data) > 0 {
		current, err = strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parsing counter %q: %w", path, err)
		}
	}

	next := current + 1
	if err := os.WriteFile(path, []byte(strconv.FormatInt(next, 10)), 0o644); err != nil {
		return 0, err
	}

	return next, nil
}

const dateLayout = "2006-01-02"
