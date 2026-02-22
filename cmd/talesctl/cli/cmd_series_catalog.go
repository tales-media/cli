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
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

// TODO: update
// TODO: delete

func seriesCatalogCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "catalog",
		Short:                 "Manage Series Catalogs",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		seriesCatalogListCommand(cfg),
		seriesCatalogGetCommand(cfg),
	)
	return cmd
}

func seriesCatalogListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list [series id]",
		"List Series Catalogs",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.SeriesCatalog
				req svc.SeriesCatalogListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastSeriesCatalog(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesSeriesCatalog(extAPI) },
			})()

			req.SeriesID = args[0]

			return s.List(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func seriesCatalogGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [series id] [catalog flavor]",
		"Get a Series Catalog",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.SeriesCatalog
				req svc.SeriesCatalogGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastSeriesCatalog(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesSeriesCatalog(extAPI) },
			})()

			req.SeriesID = args[0]
			req.Flavor = api.Flavor(args[1])

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(2)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
