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
	"errors"
	"net/http"
)

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

func GenericAutoDecodedDo[T any](do Doer, reqFunc func() (*Request, error)) (T, *Response, error) {
	var data T

	// build request
	req, err := reqFunc()
	if err != nil {
		return data, nil, err
	}

	// send request
	resp, err := do.Do(req)
	if err != nil {
		return data, nil, err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return data, resp, responseErr(resp)
	}

	// decode response
	decData := new(T)
	err = resp.Decode(decData, AutoDecoder)
	if err != nil {
		return data, resp, err
	}
	data = *decData
	return data, resp, nil
}

func responseErr(resp *Response) error {
	if 400 <= resp.StatusCode {
		st := http.StatusText(resp.StatusCode)
		return errors.New(st)
	}
	return UnexpectedStatusCodeErr
}
