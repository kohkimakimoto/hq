package client

import (
	"github.com/kohkimakimoto/hq/structs"
	"net/http"
	"strings"
)

type Client struct {
	address string
}

func New(address string) *Client {
	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	return &Client{
		address: address,
	}
}

func (c *Client) do(method, url string, payload interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Info() error {

	return nil
}

func (c *Client) CreateJob(req *structs.CreateJobRequest) error {
	return nil
}
