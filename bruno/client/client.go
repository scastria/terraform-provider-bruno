package client

type Client struct {
	collectionPath string
}

func NewClient(collectionPath string) (*Client, error) {
	c := &Client{
		collectionPath: collectionPath,
	}
	return c, nil
}
