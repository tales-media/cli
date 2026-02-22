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

func infoCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info",
		Short:                 "Print Information",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		infoApiCommand(cfg),
		infoMeCommand(cfg),
		infoOrganizationCommand(cfg),
	)
	return cmd
}

func infoApiCommand(cfg *Config) *cobra.Command {
	return extAPICommand(
		"api",
		mustSelect(cfg.AliasType, map[AliasType]string{
			OpencastAlias: "Print Opencast API Information",
			TalesAlias:    "Print tales.media API Information",
		}),
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var s svc.Info
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastInfo(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesInfo(extAPI) },
			})()
			return s.API(cmd.Context())
		},
	)
}

func infoMeCommand(cfg *Config) *cobra.Command {
	return extAPICommand(
		"me",
		mustSelect(cfg.AliasType, map[AliasType]string{
			OpencastAlias: "Print Opencast User Information",
			TalesAlias:    "Print tales.media User Information",
		}),
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var s svc.Info
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastInfo(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesInfo(extAPI) },
			})()
			return s.Me(cmd.Context())
		},
	)
}

func infoOrganizationCommand(cfg *Config) *cobra.Command {
	return extAPICommand(
		"organization",
		mustSelect(cfg.AliasType, map[AliasType]string{
			OpencastAlias: "Print Opencast Organization Information",
			TalesAlias:    "Print tales.media Organization Information",
		}),
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var s svc.Info
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastInfo(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesInfo(extAPI) },
			})()
			return s.Organization(cmd.Context())
		},
	)
}
