package client

import (
	"os"
	"path/filepath"
)

type Client struct {
	collectionPath string
}

func NewClient(collectionPath string) (*Client, error) {
	err := os.MkdirAll(collectionPath, 0755)
	if err != nil {
		return nil, err
	}
	c := &Client{
		collectionPath: collectionPath,
	}
	return c, nil
}

func (c *Client) GetAbsolutePath(relativePaths ...string) string {
	relativePaths = append([]string{c.collectionPath}, relativePaths...)
	return filepath.Join(relativePaths...)
}

func (c *Client) GetRelativePath(absolutePath string) (string, error) {
	return filepath.Rel(c.collectionPath, absolutePath)
}
