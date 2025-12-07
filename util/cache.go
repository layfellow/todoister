// Copyright 2025 layfellow. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import (
	"os"
	"path/filepath"

	"google.golang.org/protobuf/proto"
)

const (
	CacheFileName = "todoist.pb"
)

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
// Returns nil if the cache doesn't exist or is corrupted (no error).
// This allows the caller to proceed with a full sync.
func LoadCache() (*CachedTodoistData, error) {
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

	return nil
}
