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

// TODO: create
// TODO: delete
// TODO: update

func groupCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "group",
		Short:                 "Manage Groups",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		groupListCommand(cfg),
		groupGetCommand(cfg),
	)
	return cmd
}

func groupListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list",
		"List Groups",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Group
				req svc.GroupListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastGroup(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesGroup(extAPI) },
			})()

			req.FilterByName = getFilterByXStringFlag("name", cmd.Flags())
			req.SortBy = getSortByFlag(cmd.Flags())
			req.SortDirection = getSortDirectionFlag(cmd.Flags())

			return s.List(cmd.Context(), req)
		},
	)
	addFilterByXStringFlag("name", cmd.Flags())
	addSortByFlag(&mapValue[string]{
		Default: "name",
		Map: map[string]string{
			"name":        "name",
			"description": "description",
			"role":        "role",
		},
	}, cmd.Flags())
	addSortDirectionFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func groupGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get a Group",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Group
				req svc.GroupGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastGroup(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesGroup(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
