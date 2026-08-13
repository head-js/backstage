package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadContext(context string) (string, error) {
	if len(context) < 64 && strings.HasSuffix(context, ".md") {
		if !strings.HasPrefix(context, ".") {
			return "", InvalidFormatException("context must be a relative path starting with '.'")
		}
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		path := filepath.Join(pwd, context)
		if _, err := os.Stat(path); err == nil {
			content, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("failed to read context file: %w", err)
			}
			return string(content), nil
		}
	}
	return context, nil
}
