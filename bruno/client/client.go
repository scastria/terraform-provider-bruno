package client

import (
	"os"
	"path/filepath"
)

type Client struct {
	collectionPath string
}

func NewClient(collectionPath string) (*Client, error) {
	// Check for collection path existence
	_, err := os.Stat(collectionPath)
	if err != nil {
		return nil, err
	}
	c := &Client{
		collectionPath: collectionPath,
	}
	return c, nil
}

func (c *Client) GetPath(path string) string {
	return filepath.Join(c.collectionPath, path)
}
