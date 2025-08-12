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

	"github.com/tales-media/cli/internal/talesctl/svc"
	extapiclientv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

func workflowDefinitionCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow-definition",
		Short:                 "Manage Workflow Definitions",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddCommand(
		workflowDefinitionListCommand(cfg),
		workflowDefinitionGetCommand(cfg),
	)
	return cmd
}

func workflowDefinitionListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list",
		"List Workflow Definitions",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.WorkflowDefinition
				req svc.WorkflowDefinitionListRequest
				err error
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflowDefinition(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflowDefinition(extAPI) },
			})()

			if req.FilterByTag, err = cmd.Flags().GetString(TagFlag); err != nil {
				return nil, err
			}

			return s.List(cmd.Context(), req)
		},
	)
	cmd.Flags().String(TagFlag, "", "only include workflow definitions with this tag")
	return cmd
}

func workflowDefinitionGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get a Workflow Definition",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.WorkflowDefinition
				req svc.WorkflowDefinitionGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflowDefinition(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflowDefinition(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}
