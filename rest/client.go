package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/xlog"
)

// REST Client
type Client struct {
	apiAddr string
}

const DefaultAPIAddress = "https://rest.mediacoin.pro/"

func NewClient(apiAddr string) *Client {
	if apiAddr == "" {
		apiAddr = DefaultAPIAddress
	}
	if !strings.HasPrefix(apiAddr, "http") {
		apiAddr = "http://" + apiAddr
	}
	return &Client{
		apiAddr: apiAddr,
	}
}

func (c *Client) httpGet(path string, res interface{}) (err error) {
	return c.httpRequest("GET", path, nil, res)
}

func (c *Client) httpPut(path string, req interface{}) (err error) {
	return c.httpRequest("PUT", path, req, nil)
}

func (c *Client) httpRequest(method, path string, reqObj, resObj interface{}) (err error) {
	if xlog.TraceIsOn() {
		xlog.Trace.Printf("rest> http-req: %s %s ...", method, c.apiAddr+path)
	}
	req, err := http.NewRequest(method, c.apiAddr+path, nil)
	if reqObj != nil {
		req.Body = bin.NewBuffer(nil, reqObj)
	}
	if err != nil {
		xlog.Error.Printf("rest> http-req-error: %v", err)
		return
	}
	req.Header.Set("Accept", "binary")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		xlog.Error.Printf("rest> http-resp-error: %v", err)
		return
	}
	defer resp.Body.Close()

	if xlog.TraceIsOn() {
		xlog.Trace.Printf("rest> http-resp-code: %d", resp.StatusCode)
	}
	if resp.StatusCode != 200 {
		res, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(res))
	}
	if resObj != nil {
		err = bin.Read(resp.Body, resObj)
	}
	return
}
