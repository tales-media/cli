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
	"errors"
	"slices"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
)

type SeriesCatalog interface {
	List(context.Context, SeriesCatalogListRequest) ([]api.Catalog, error)
	Get(context.Context, SeriesCatalogGetRequest) (api.Catalog, error)
}

type SeriesCatalogListRequest struct {
	SeriesID string
}

type SeriesCatalogGetRequest struct {
	SeriesID string
	Flavor   api.Flavor
}

type opencastSeriesCatalog struct {
	extAPI extapiclientv1.Client
}

var _ SeriesCatalog = &opencastSeriesCatalog{}

func NewOpencastSeriesCatalog(extAPI extapiclientv1.Client) SeriesCatalog {
	return &opencastSeriesCatalog{
		extAPI: extAPI,
	}
}

func (svc *opencastSeriesCatalog) List(ctx context.Context, req SeriesCatalogListRequest) ([]api.Catalog, error) {
	ocMetadata, _, err := svc.extAPI.ListSeriesMetadata(
		ctx,
		req.SeriesID,
	)
	if err != nil {
		return nil, err
	}
	catalogs := conv.Map(ocMetadata, conv.OCCatalogToCatalog)
	return catalogs, nil
}

func (svc *opencastSeriesCatalog) Get(ctx context.Context, req SeriesCatalogGetRequest) (api.Catalog, error) {
	ocMetadata, _, err := svc.extAPI.ListSeriesMetadata(
		ctx,
		req.SeriesID,
	)
	if err != nil {
		return api.Catalog{}, err
	}

	flavor := conv.FlavorToOCFlavor(req.Flavor)
	i := slices.IndexFunc(ocMetadata, func(c extapiv1.Catalog) bool { return c.Flavor == flavor })
	if i < 0 {
		return api.Catalog{}, errors.New("Not Found")
	}
	ocCatalog := &ocMetadata[i]

	return conv.OCCatalogToCatalog(*ocCatalog), nil
}

type talesSeriesCatalog = opencastSeriesCatalog

var _ SeriesCatalog = &talesSeriesCatalog{}

var NewTalesSeriesCatalog = NewOpencastSeriesCatalog
