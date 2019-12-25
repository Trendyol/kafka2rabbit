package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"runtime"
	"time"
)

type RestClient struct {
	ctx   context.Context
	agent *http.Client
	req   *http.Request
}

func NewClient() *RestClient {
	return &RestClient{
		ctx: context.Background(),
		agent: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
		req: &http.Request{},
	}
}

func (c *RestClient) WithHeader(key, value string) {
	c.req.Header.Add(key, value)
}

func (c *RestClient) Post(url string, body []byte) *RestClient {
	c.req, _ = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	return c
}

func (c *RestClient) Get(url string) *RestClient {
	c.req, _ = http.NewRequest(http.MethodGet, url, nil)
	return c
}

func (c *RestClient) End() (*http.Response, error) {
	resp, err := c.agent.Do(c.req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RestClient) EndStruct(aStruct interface{}) error {
	resp, err := c.agent.Do(c.req)
	if err != nil {
		return fmt.Errorf("dial timeout,err:%v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode >= fasthttp.StatusBadRequest {
		return fmt.Errorf("bad request")
	}

	err = json.NewDecoder(resp.Body).Decode(aStruct)
	if err != nil {
		return fmt.Errorf("unmarshalling error, err:%v", err)
	}
	return nil
}
