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

package formatter

import (
	"errors"
	"io"
	"reflect"
	"strings"
)

type Formatter interface {
	List(w io.Writer, list any) error
	Object(w io.Writer, obj any) error
	Error(w io.Writer, err error) error
}

type errorObj struct {
	Error string `json:"error" yaml:"error"`
}

type genericList struct {
	Count int `json:"count" yaml:"count"`
	Items any `json:"items" yaml:"items"`
}

func NewGenericList(list any) (*genericList, error) {
	lv := reflect.ValueOf(list)
	if lv.Type().Kind() != reflect.Slice {
		return nil, errors.New("input is not a slice")
	}
	gl := &genericList{
		Count: lv.Len(),
		Items: list,
	}
	return gl, nil
}

func ParseStructTag(sf reflect.StructField, key string) (string, StructTagOptions) {
	st := sf.Tag.Get(key)
	if len(st) == 0 {
		return sf.Name, ""
	}
	n, opt, _ := strings.Cut(st, ",")
	return n, StructTagOptions(opt)
}

type StructTagOptions string

func (o StructTagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}
