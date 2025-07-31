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
	"strconv"
	"strings"

	oc "github.com/tales-media/cli/pkg/opencast/client"
)

type WithPagination struct {
	Limit  int
	Offset int
}

func (opt WithPagination) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("limit", strconv.Itoa(opt.Limit)),
		oc.WithQuery("offset", strconv.Itoa(opt.Offset)),
	)
}

type WithFilter map[FilterKey]string

type FilterKey string

func (opt WithFilter) Apply(r *oc.Request) error {
	if len(opt) == 0 {
		return nil
	}

	sb := strings.Builder{}

	// grow string builder once
	n := len(opt) - 1 // comma
	for k, v := range opt {
		n = n +
			len(k) + // key
			1 + // colon
			len(v) // value
	}
	sb.Grow(n)

	// build filter
	for k, v := range opt {
		if sb.Len() > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(string(k))
		sb.WriteString(":")
		sb.WriteString(v)
	}

	return r.ApplyOptions(oc.WithQuery("filter", sb.String()))
}

type WithSort []Sort

func (opt WithSort) Apply(r *oc.Request) error {
	if len(opt) == 0 {
		return nil
	}

	sb := strings.Builder{}

	// grow string builder once
	n := len(opt) - 1 // comma
	for _, s := range opt {
		n = n +
			len(s.By) + // by
			1 + // colon
			len(s.Direction) // direction
	}

	// build filter
	for i, s := range opt {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(string(s.By))
		sb.WriteString(":")
		sb.WriteString(string(s.Direction))
	}

	return r.ApplyOptions(oc.WithQuery("sort", sb.String()))
}

type Sort struct {
	By        SortKey
	Direction SortDirection
}

type SortKey string

type SortDirection string

const (
	Ascending  = SortDirection("ASC")
	Descending = SortDirection("DESC")
)

func init() {
	_ = WithSort{{By: "name", Direction: Ascending}}
}
