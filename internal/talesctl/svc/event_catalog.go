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

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
	extapiv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

type EventCatalog interface {
	List(context.Context, EventCatalogListRequest) ([]api.Catalog, error)
	Get(context.Context, EventCatalogGetRequest) (api.Catalog, error)
}

type EventCatalogListRequest struct {
	EventID string
}

type EventCatalogGetRequest struct {
	EventID string
	Flavor  api.Flavor
}

type opencastEventCatalog struct {
	extAPI extapiclientv1.Client
}

var _ EventCatalog = &opencastEventCatalog{}

func NewOpencastEventCatalog(extAPI extapiclientv1.Client) EventCatalog {
	return &opencastEventCatalog{
		extAPI: extAPI,
	}
}

func (svc *opencastEventCatalog) List(ctx context.Context, req EventCatalogListRequest) ([]api.Catalog, error) {
	ocMetadata, _, err := svc.extAPI.ListEventMetadata(
		ctx,
		req.EventID,
	)
	if err != nil {
		return nil, err
	}
	catalogs := conv.Map(ocMetadata, conv.OCCatalogToCatalog)
	return catalogs, nil
}

func (svc *opencastEventCatalog) Get(ctx context.Context, req EventCatalogGetRequest) (api.Catalog, error) {
	ocMetadata, _, err := svc.extAPI.ListEventMetadata(
		ctx,
		req.EventID,
	)
	if err != nil {
		return api.Catalog{}, err
	}

	var ocCatalog *extapiv1.Catalog
	flavor := conv.FlavorToOCFlavor(req.Flavor)
	for _, c := range ocMetadata {
		if c.Flavor == flavor {
			ocCatalog = &c
			break
		}
	}
	if ocCatalog == nil {
		return api.Catalog{}, errors.New("Not Found")
	}

	return conv.OCCatalogToCatalog(*ocCatalog), nil
}

type talesEventCatalog = opencastEventCatalog

var _ EventCatalog = &talesEventCatalog{}

var NewTalesEventCatalog = NewOpencastEventCatalog
