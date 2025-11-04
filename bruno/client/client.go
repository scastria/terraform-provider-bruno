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

func (c *Client) GetAbsolutePath(relativePath string) string {
	return filepath.Join(c.collectionPath, relativePath)
}

//func (c *Client) GetRelativePath(absolutePath string) string {
//	return filepath.Join(c.collectionPath, relativePath)
//}
