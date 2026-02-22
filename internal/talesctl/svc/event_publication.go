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

type EventPublication interface {
	List(context.Context, EventPublicationListRequest) ([]api.Publication, error)
	Get(context.Context, EventPublicationGetRequest) (api.Publication, error)
}

type EventPublicationListRequest struct {
	EventID string
}

type EventPublicationGetRequest struct {
	EventID       string
	PublicationID string
}

type opencastEventPublication struct {
	extAPI extapiclientv1.Client
}

var _ EventPublication = &opencastEventPublication{}

func NewOpencastEventPublication(extAPI extapiclientv1.Client) EventPublication {
	return &opencastEventPublication{
		extAPI: extAPI,
	}
}

func (svc *opencastEventPublication) List(ctx context.Context, req EventPublicationListRequest) ([]api.Publication, error) {
	ocPublications, _, err := svc.extAPI.ListEventPublication(
		ctx,
		req.EventID,
		extapiclientv1.WithEventOptions{
			IncludeInternalPublication: true,
		},
	)
	if err != nil {
		return nil, err
	}
	publications := conv.Map(ocPublications, conv.OCPublicationToPublication)
	return publications, nil
}

func (svc *opencastEventPublication) Get(ctx context.Context, req EventPublicationGetRequest) (api.Publication, error) {
	ocPublications, _, err := svc.extAPI.ListEventPublication(
		ctx,
		req.EventID,
		extapiclientv1.WithEventOptions{
			IncludeInternalPublication: true,
		},
	)
	if err != nil {
		return api.Publication{}, err
	}

	i := slices.IndexFunc(ocPublications, func(pub extapiv1.Publication) bool {
		return pub.ID == req.PublicationID || pub.Channel == req.PublicationID
	})
	if i < 0 {
		return api.Publication{}, errors.New("Not Found")
	}
	ocPublication := &ocPublications[i]

	return conv.OCPublicationToPublication(*ocPublication), nil
}

type talesEventPublication = opencastEventPublication

var _ EventPublication = &talesEventPublication{}

var NewTalesEventPublication = NewOpencastEventPublication
