package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Client struct {
	client *memcache.Client
}

func New(servers []string) *Client {
	return &Client{client: memcache.New(servers...)}
}

func (c *Client) Get(key string) (*memcache.Item, error) {
	return c.client.Get(key)
}

func (c *Client) Set(item *memcache.Item) error {
	return c.client.Set(item)
}

func (c *Client) Delete(key string) error {
	return c.client.Delete(key)
}
