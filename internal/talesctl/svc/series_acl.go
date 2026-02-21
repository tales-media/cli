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

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

type SeriesACL interface {
	Get(context.Context, SeriesACLGetRequest) ([]api.ACE, error)
}

type SeriesACLGetRequest struct {
	SeriesID string
}

type opencastSeriesACL struct {
	extAPI extapiclientv1.Client
}

var _ SeriesACL = &opencastSeriesACL{}

func NewOpencastSeriesACL(extAPI extapiclientv1.Client) SeriesACL {
	return &opencastSeriesACL{
		extAPI: extAPI,
	}
}

func (svc *opencastSeriesACL) Get(ctx context.Context, req SeriesACLGetRequest) ([]api.ACE, error) {
	ocACL, _, err := svc.extAPI.GetSeriesACL(
		ctx,
		req.SeriesID,
	)
	if err != nil {
		return nil, err
	}
	acl := conv.Map(ocACL, conv.OCACEToACE)
	return acl, nil
}

type talesSeriesACL = opencastSeriesACL

var _ SeriesACL = &talesSeriesACL{}

var NewTalesSeriesACL = NewOpencastSeriesACL
