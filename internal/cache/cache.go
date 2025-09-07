package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CachedMessage struct {
	Message   string    `json:"message"`
	Files     []string  `json:"files"`
	DiffHash  string    `json:"diff_hash"`
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
}

type Cache struct {
	cacheDir string
}

func New() *Cache {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "commitgen")
	_ = os.MkdirAll(cacheDir, 0755) // ignore error, cache is optional
	return &Cache{cacheDir: cacheDir}
}

func (c *Cache) GetCacheKey(files []string, patch string) string {
	h := sha256.New()
	for _, file := range files {
		h.Write([]byte(file))
	}
	h.Write([]byte(patch))
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

func (c *Cache) Get(files []string, patch string) (*CachedMessage, error) {
	key := c.GetCacheKey(files, patch)
	cachePath := filepath.Join(c.cacheDir, key+".json")

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cached CachedMessage
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	if time.Since(cached.Timestamp) > 24*time.Hour {
		os.Remove(cachePath)
		return nil, fmt.Errorf("cache expired")
	}

	return &cached, nil
}

func (c *Cache) Set(files []string, patch string, message string, provider string) error {
	key := c.GetCacheKey(files, patch)
	cachePath := filepath.Join(c.cacheDir, key+".json")

	cached := CachedMessage{
		Message:   message,
		Files:     files,
		DiffHash:  key,
		Timestamp: time.Now(),
		Provider:  provider,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (c *Cache) GetLatest() (*CachedMessage, error) {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		return nil, err
	}

	var latest *CachedMessage
	var latestTime time.Time

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			cachePath := filepath.Join(c.cacheDir, entry.Name())
			data, err := os.ReadFile(cachePath)
			if err != nil {
				continue
			}

			var cached CachedMessage
			if err := json.Unmarshal(data, &cached); err != nil {
				continue
			}

			if cached.Timestamp.After(latestTime) {
				latest = &cached
				latestTime = cached.Timestamp
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no cached messages found")
	}

	return latest, nil
}

func (c *Cache) Clear() error {
	return os.RemoveAll(c.cacheDir)
}
