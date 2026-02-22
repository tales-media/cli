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

package svc

import (
	"context"

	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
)

type SeriesProperty interface {
	List(context.Context, SeriesPropertyListRequest) ([]api.Property, error)
}

type SeriesPropertyListRequest struct {
	SeriesID string
}

type opencastSeriesProperty struct {
	extAPI extapiclientv1.Client
}

var _ SeriesProperty = &opencastSeriesProperty{}

func NewOpencastSeriesProperty(extAPI extapiclientv1.Client) SeriesProperty {
	return &opencastSeriesProperty{
		extAPI: extAPI,
	}
}

func (svc *opencastSeriesProperty) List(ctx context.Context, req SeriesPropertyListRequest) ([]api.Property, error) {
	ocProperties, _, err := svc.extAPI.GetSeriesProperties(
		ctx,
		req.SeriesID,
	)
	if err != nil {
		return nil, err
	}
	return conv.OCPropertiesToProperties(ocProperties), nil
}

type talesSeriesProperty = opencastSeriesProperty

var _ SeriesProperty = &talesSeriesProperty{}

var NewTalesSeriesProperty = NewOpencastSeriesProperty
