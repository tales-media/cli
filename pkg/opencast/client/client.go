/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"net/http"
	"time"
)

type Doer interface {
	Do(*Request) (*Response, error)
}

type Client interface {
	Doer
}

type client struct {
	sm      ServiceMapper
	http    http.Client
	reqOpts []RequestOpts
}

var _ Client = &client{}

func New(sm ServiceMapper, opts ...ClientOpts) (Client, error) {
	c := &client{
		sm: sm,
	}
	if err := applyOpts(c, opts); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) ApplyOptions(opts ...ClientOpts) error {
	return applyOpts(c, opts)
}

func (c *client) Do(req *Request) (*Response, error) {
	if err := applyOpts(req, c.reqOpts); err != nil {
		return nil, err
	}

	httpReq, err := req.HTTPRequest(c.sm)
	if err != nil {
		return nil, err
	}

	reqStart := time.Now()
	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	resp := newResponse(httpResp)
	resp.Meta.Duration = time.Since(reqStart)

	return resp, nil
}

type ClientOpts = func(*client) error

func WithHTTPClient(h http.Client) ClientOpts {
	return func(c *client) error {
		c.http = h
		return nil
	}
}

func WithRequestOptions(opts ...RequestOpts) ClientOpts {
	return func(c *client) error {
		c.reqOpts = opts
		return nil
	}
}

// TODO: add WithBackoff
