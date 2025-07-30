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

import "errors"

var (
	UnexpectedStatusCodeErr = errors.New("UnexpectedStatusCode")
)

func applyOpts[T any](obj T, opts []func(T) error) error {
	for _, f := range opts {
		if err := f(obj); err != nil {
			return err
		}
	}
	return nil
}

func GenericAutoDecodedCall[T any](reqFunc func() (*Request, error), doFunc func(*Request) (*Response, error)) (*T, *Response, error) {
	req, err := reqFunc()
	if err != nil {
		return nil, nil, err
	}
	resp, err := doFunc(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, resp, UnexpectedStatusCodeErr
	}
	data := new(T)
	err = resp.Decode(data, AutoDecoder)
	if err != nil {
		return nil, resp, err
	}
	return data, resp, nil
}
