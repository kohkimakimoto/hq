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
	"net/url"
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

func (c *Client) post(url string, payload interface{}) (*http.Response, error) {
	var payloadBytes []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		payloadBytes = b
	}

	req, err := http.NewRequest("POST", c.address+url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := c.checkStatusCode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) get(url string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.address+url, nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	if values != nil {
		q := req.URL.Query()
		for k, v := range values {
			for _, vs := range v {
				q.Add(k, vs)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := c.checkStatusCode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", DefaultUserAgent)
}

func (c *Client) checkStatusCode(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, http.StatusText(resp.StatusCode))
		}

		ret := &structs.ErrorResponse{}
		if err := json.Unmarshal(body, ret); err != nil {
			return errors.Wrap(err, http.StatusText(resp.StatusCode))
		}

		return errors.New(ret.Error)
	}

	return nil
}

func (c *Client) Info() (*structs.Info, error) {
	resp, err := c.post("/", nil)
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
	resp, err := c.post("/job", payload)
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

func (c *Client) ListJobs(payload *structs.ListJobsRequest) (*structs.JobList, error) {
	//params["name"] = payload.Name

	//if payload.Begin != "" {
	//	params["begin"] = payload.Begin
	//}
	var values url.Values = url.Values{}

	if payload.Name != ""  {
		values.Add("name", payload.Name)
	}

	if payload.Begin != nil  {
		values.Add("begin", fmt.Sprintf("%d", *payload.Begin))
	}

	if payload.Reverse {
		values.Add("reverse", fmt.Sprintf("%v", payload.Reverse))
	}

	if payload.Limit != 0 {
		values.Add("limit", fmt.Sprintf("%d", payload.Limit))
	}

	resp, err := c.get("/job", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &structs.JobList{
		Jobs: []*structs.Job{},
	}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Client) GetJob(id uint64) (*structs.Job, error) {
	resp, err := c.get(fmt.Sprintf("/job/%d", id), nil)
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

func (c *Client) Stats() (*structs.Stats, error) {
	resp, err := c.get("/stats", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &structs.Stats{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
