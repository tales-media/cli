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

// TODO: update

func eventACLCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "acl",
		Short:                 "Manage Event ACL",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		eventACLListCommand(cfg),
	)
	return cmd
}

func eventACLListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list [event id]",
		"List Event ACL",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.EventACL
				req svc.EventACLListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastEventACL(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesEventACL(extAPI) },
			})()

			req.EventID = args[0]

			return s.List(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
