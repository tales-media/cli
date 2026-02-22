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

type PlaylistEntry interface {
	List(context.Context, PlaylistEntryListRequest) ([]api.PlaylistEntry, error)
}

type PlaylistEntryListRequest struct {
	PlaylistID string
}

type opencastPlaylistEntry struct {
	extAPI extapiclientv1.Client
}

var _ PlaylistEntry = &opencastPlaylistEntry{}

func NewOpencastPlaylistEntry(extAPI extapiclientv1.Client) PlaylistEntry {
	return &opencastPlaylistEntry{
		extAPI: extAPI,
	}
}

func (svc *opencastPlaylistEntry) List(ctx context.Context, req PlaylistEntryListRequest) ([]api.PlaylistEntry, error) {
	ocPlaylist, _, err := svc.extAPI.GetPlaylist(
		ctx,
		req.PlaylistID,
	)
	if err != nil {
		return nil, err
	}
	entries := conv.Map(ocPlaylist.Entries, conv.OCPlaylistEntryToPlaylistEntry)
	return entries, nil
}

type talesPlaylistEntry = opencastPlaylistEntry

var _ PlaylistEntry = &talesPlaylistEntry{}

var NewTalesPlaylistEntry = NewOpencastPlaylistEntry
