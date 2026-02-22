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

	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc"
)

// TODO: update
// TODO: delete

func playlistEntryCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "entry",
		Short:                 "Manage Playlist Entries",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		playlistEntryListCommand(cfg),
	)
	return cmd
}

func playlistEntryListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list [playlist id]",
		"List Playlist Entries",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.PlaylistEntry
				req svc.PlaylistEntryListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastPlaylistEntry(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesPlaylistEntry(extAPI) },
			})()

			req.PlaylistID = args[0]

			return s.List(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
