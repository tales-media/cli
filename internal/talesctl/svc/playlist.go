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
	extapiv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
	oc "shio.solutions/tales.media/cli/pkg/opencast/client"
)

type Playlist interface {
	List(context.Context, PlaylistListRequest) ([]api.Playlist, error)
	Get(context.Context, PlaylistGetRequest) (api.Playlist, error)
}

type PlaylistListRequest struct {
	SortBy        string // TODO: add conversion function
	SortDirection api.SortDirection
}

type PlaylistGetRequest struct {
	ID string
}

type opencastPlaylist struct {
	extAPI extapiclientv1.Client
}

var _ Playlist = &opencastPlaylist{}

func NewOpencastPlaylist(extAPI extapiclientv1.Client) Playlist {
	return &opencastPlaylist{
		extAPI: extAPI,
	}
}

func (svc *opencastPlaylist) List(ctx context.Context, req PlaylistListRequest) ([]api.Playlist, error) {
	commonReqOpts := []oc.RequestOpts{}
	if req.SortBy != "" {
		commonReqOpts = append(commonReqOpts, extapiclientv1.WithSort{{
			By:        extapiclientv1.SortKey(req.SortBy),
			Direction: conv.SortDirectionToOCSortDirection(req.SortDirection),
		}})
	}

	var ocPlaylists []extapiv1.Playlist
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			ocReq, err := svc.extAPI.ListPlaylistRequest(
				ctx,
				extapiclientv1.WithPagination{
					Limit:  PageSize,
					Offset: i * PageSize,
				},
			)
			if err != nil {
				return nil, err
			}
			if err = ocReq.ApplyOptions(commonReqOpts...); err != nil {
				return nil, err
			}
			return ocReq, nil
		},
		oc.CollectAllPages(&ocPlaylists),
	)
	if err != nil {
		return nil, err
	}
	Playlists := conv.Map(ocPlaylists, conv.OCPlaylistToPlaylist)
	return Playlists, nil
}

func (svc *opencastPlaylist) Get(ctx context.Context, req PlaylistGetRequest) (api.Playlist, error) {
	ocPlaylist, _, err := svc.extAPI.GetPlaylist(
		ctx,
		req.ID,
	)
	if err != nil {
		return api.Playlist{}, err
	}
	return conv.OCPlaylistToPlaylist(*ocPlaylist), nil
}

type talesPlaylist = opencastPlaylist

var _ Playlist = &talesPlaylist{}

var NewTalesPlaylist = NewOpencastPlaylist
