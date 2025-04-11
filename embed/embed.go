package embed

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"strings"
)

//go:embed config.env
var configFS embed.FS

// LoadEmbeddedConfig loads the embedded config into environment variables
func LoadEmbeddedConfig() {
	// Read the embedded config file
	data, err := configFS.ReadFile("config.env")
	if err != nil {
		log.Printf("Warning: Could not read embedded config: %v", err)
		return
	}

	// Split by newlines and process each line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set by actual environment variable
		if os.Getenv(key) == "" {
			err := os.Setenv(key, value)
			if err != nil {
				return
			}
		}
	}
}

// GetEmbeddedFS returns the embedded file system
func GetEmbeddedFS() fs.FS {
	return configFS
}
