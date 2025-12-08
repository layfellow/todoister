// Copyright 2025 Marco Bravo Mejías. All rights reserved.
// Use of this source code is governed by a GPL v3 license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

const (
	CacheFileName   = "todoist.pb"
	VersionFileName = "version"
)

// SchemaVersion holds the application version for cache compatibility checking.
// This should be set by the cmd package during initialization.
var SchemaVersion string

// GetVersionFilePath returns the path to the version file.
func GetVersionFilePath() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userCacheDir, Prog, VersionFileName), nil
}

// getSchemaFromVersion extracts major.minor from a version string like "0.4.0"
func getSchemaFromVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return version
}

// GetCachePath returns the path to the cache file.
// Returns an error if the user cache directory cannot be determined.
func GetCachePath() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userCacheDir, Prog, CacheFileName), nil
}

// EnsureCacheDir creates the cache directory if it doesn't exist.
// Returns an error if the directory cannot be created.
func EnsureCacheDir() error {
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}
	cacheDir := filepath.Dir(cachePath)
	return os.MkdirAll(cacheDir, 0755)
}

// LoadCache reads and deserializes the Protobuf cache file.
// Returns nil if the cache doesn't exist, is corrupted, or schema version mismatch.
// This allows the caller to proceed with a full sync.
func LoadCache() (*CachedTodoistData, error) {
	// First check version compatibility if version file exists
	versionPath, err := GetVersionFilePath()
	if err == nil {
		if data, err := os.ReadFile(versionPath); err == nil {
			// Parse schema=x.y from file
			content := strings.TrimSpace(string(data))
			if strings.HasPrefix(content, "schema=") {
				cachedSchema := strings.TrimPrefix(content, "schema=")
				currentSchema := getSchemaFromVersion(SchemaVersion)
				if cachedSchema != currentSchema {
					Warn(fmt.Sprintf("Cache schema mismatch (cached: %s, current: %s), will perform full sync", cachedSchema, currentSchema), nil)
					// Overwrite version file with current schema
					newContent := fmt.Sprintf("schema=%s\n", currentSchema)
					if err := os.WriteFile(versionPath, []byte(newContent), 0644); err != nil {
						Warn("Failed to update version file", err)
					}
					return nil, nil // Version mismatch, return nil to trigger full sync
				}
			}
		}
		// If version file doesn't exist or can't be read, continue loading cache
	}

	cachePath, err := GetCachePath()
	if err != nil {
		return nil, nil // Can't get cache path, return nil to trigger full sync
	}

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, nil // Cache doesn't exist, return nil to trigger full sync
	}

	// Read cache file
	data, err := os.ReadFile(cachePath)
	if err != nil {
		Warn("Failed to read cache file, will perform full sync", err)
		return nil, nil // Can't read cache, return nil to trigger full sync
	}

	// Deserialize protobuf
	cached := &CachedTodoistData{}
	if err := proto.Unmarshal(data, cached); err != nil {
		Warn("Failed to unmarshal cache file, will perform full sync", err)
		return nil, nil // Corrupted cache, return nil to trigger full sync
	}

	return cached, nil
}

// SaveCache serializes and writes the Protobuf cache file.
// Also writes the version file if it doesn't exist.
// Returns an error if the cache cannot be saved.
func SaveCache(data *CachedTodoistData) error {
	// Ensure cache directory exists
	if err := EnsureCacheDir(); err != nil {
		return err
	}

	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// Serialize protobuf
	bytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile(cachePath, bytes, 0644); err != nil {
		return err
	}

	// Write version file if it doesn't exist
	versionPath, err := GetVersionFilePath()
	if err == nil {
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			schema := getSchemaFromVersion(SchemaVersion)
			content := fmt.Sprintf("schema=%s\n", schema)
			if err := os.WriteFile(versionPath, []byte(content), 0644); err != nil {
				Warn("Failed to write version file", err)
			}
		}
	}

	return nil
}
