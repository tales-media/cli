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

type EventTrack interface {
	List(context.Context, EventTrackListRequest) ([]api.TrackElement, error)
}

type EventTrackListRequest struct {
	EventID string
}

type opencastEventTrack struct {
	extAPI extapiclientv1.Client
}

var _ EventTrack = &opencastEventTrack{}

func NewOpencastEventTrack(extAPI extapiclientv1.Client) EventTrack {
	return &opencastEventTrack{
		extAPI: extAPI,
	}
}

func (svc *opencastEventTrack) List(ctx context.Context, req EventTrackListRequest) ([]api.TrackElement, error) {
	ocMedia, _, err := svc.extAPI.ListEventMedia(
		ctx,
		req.EventID,
	)
	if err != nil {
		return nil, err
	}
	tracks := conv.Map(ocMedia, conv.OCMediaTrackElementToTrackElement)
	return tracks, nil
}

type talesEventTrack = opencastEventTrack

var _ EventTrack = &talesEventTrack{}

var NewTalesEventTrack = NewOpencastEventTrack
