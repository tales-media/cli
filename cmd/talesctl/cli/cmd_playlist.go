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

package cli

import (
	"github.com/spf13/cobra"
	"shio.solutions/tales.media/cli/internal/talesctl/svc"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

// TODO: create
// TODO: delete
// TODO: update

func playlistCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "playlist",
		Short:                 "Manage Playlists",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ResourcesGroup,
		ManagementGroup,
	)
	cmd.AddCommand(
		// resources
		playlistACLCommand(cfg),
		playlistEntryCommand(cfg),

		// management
		playlistListCommand(cfg),
		playlistGetCommand(cfg),
	)
	return cmd
}

func playlistListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list",
		"List Playlists",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Playlist
				req svc.PlaylistListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastPlaylist(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesPlaylist(extAPI) },
			})()

			req.SortBy = getSortByFlag(cmd.Flags())
			req.SortDirection = getSortDirectionFlag(cmd.Flags())

			return s.List(cmd.Context(), req)
		},
	)
	addSortByFlag(&mapValue[string]{
		Default: "updated",
		Map: map[string]string{
			"updated": "updated",
		},
	}, cmd.Flags())
	addSortDirectionFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func playlistGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get an Playlist",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Playlist
				req svc.PlaylistGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastPlaylist(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesPlaylist(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
