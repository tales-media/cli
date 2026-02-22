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

func seriesCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "series",
		Short:                 "Manage Series",
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
		seriesACLCommand(cfg),
		seriesCatalogCommand(cfg),
		seriesPropertyCommand(cfg),

		// management
		seriesListCommand(cfg),
		seriesGetCommand(cfg),
	)
	return cmd
}

func seriesListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list",
		"List Series",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Series
				req svc.SeriesListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastSeries(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesSeries(extAPI) },
			})()

			req.FilterByText = getFilterByXStringFlag("text", cmd.Flags())
			req.FilterByID = getFilterByXStringFlag("id", cmd.Flags())
			req.FilterByTitle = getFilterByXStringFlag("title", cmd.Flags())
			req.FilterByDescription = getFilterByXStringFlag("description", cmd.Flags())
			req.FilterByCreationDate = getFilterByXStringFlag("creation-date", cmd.Flags())
			req.FilterByCreator = getFilterByXStringFlag("creator", cmd.Flags())
			req.FilterByContributors = getFilterByXStringFlag("contributors", cmd.Flags())
			req.FilterByOrganizers = getFilterByXStringFlag("organizers", cmd.Flags())
			req.FilterByPublishers = getFilterByXStringFlag("publishers", cmd.Flags())
			req.FilterByLanguage = getFilterByXStringFlag("language", cmd.Flags())
			req.FilterByRightsHolder = getFilterByXStringFlag("rights-holder", cmd.Flags())
			req.FilterByLicense = getFilterByXStringFlag("license", cmd.Flags())
			req.FilterBySubject = getFilterByXStringFlag("subject", cmd.Flags())
			req.FilterByManagedAcl = getFilterByXStringFlag("managed-acl", cmd.Flags())
			req.SortBy = getSortByFlag(cmd.Flags())
			req.SortDirection = getSortDirectionFlag(cmd.Flags())

			return s.List(cmd.Context(), req)
		},
	)
	addFilterByXStringFlag("text", cmd.Flags())
	addFilterByXStringFlag("id", cmd.Flags())
	addFilterByXStringFlag("title", cmd.Flags())
	addFilterByXStringFlag("description", cmd.Flags())
	addFilterByXStringFlag("creation-date", cmd.Flags())
	addFilterByXStringFlag("creator", cmd.Flags())
	addFilterByXStringFlag("contributors", cmd.Flags())
	addFilterByXStringFlag("organizers", cmd.Flags())
	addFilterByXStringFlag("publishers", cmd.Flags())
	addFilterByXStringFlag("language", cmd.Flags())
	addFilterByXStringFlag("rights-holder", cmd.Flags())
	addFilterByXStringFlag("license", cmd.Flags())
	addFilterByXStringFlag("subject", cmd.Flags())
	addFilterByXStringFlag("managed-acl", cmd.Flags())
	addSortByFlag(&mapValue[string]{
		Default: "title",
		Map: map[string]string{
			"title":        "title",
			"creator":      "creator",
			"created":      "created",
			"contributors": "contributors",
		},
	}, cmd.Flags())
	addSortDirectionFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func seriesGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get a Series",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Series
				req svc.SeriesGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastSeries(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesSeries(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
