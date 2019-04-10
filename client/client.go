package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/kohkimakimoto/hq/structs"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	address    string
	httpClient *http.Client
}

func New(address string) *Client {
	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	return &Client{
		address:    address,
		httpClient: http.DefaultClient,
	}
}

var (
	DefaultUserAgent = fmt.Sprintf("%s-Client/%s", hq.DisplayName, hq.Version)
)

func (c *Client) do(method, url string, payload interface{}) (*http.Response, error) {
	var payloadBytes []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		payloadBytes = b
	}

	req, err := http.NewRequest(method, c.address+url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", DefaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, http.StatusText(resp.StatusCode))
		}

		ret := &structs.ErrorResponse{}
		if err := json.Unmarshal(body, ret); err != nil {
			return nil, errors.Wrap(err, http.StatusText(resp.StatusCode))
		}

		return nil, errors.New(ret.Error)
	}

	return resp, nil
}

func (c *Client) Info() (*structs.Info, error) {
	resp, err := c.do("POST", "/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &structs.Info{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Client) CreateJob(payload *structs.CreateJobRequest) (*structs.Job, error) {
	resp, err := c.do("POST", "/job", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &structs.Job{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
