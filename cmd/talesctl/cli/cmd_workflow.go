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
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

// TODO: add list

func workflowCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "workflow",
		Short:                 "Manage Workflows",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		workflowCreateCommand(cfg),
		workflowGetCommand(cfg),
		workflowUpdateCommand(cfg),
		workflowDeleteCommand(cfg),
	)
	return cmd
}

func workflowCreateCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"create [eventID] [workflowDefinitionID]",
		"Create a Workflow",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Workflow
				req svc.WorkflowCreateRequest
				err error
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflow(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflow(extAPI) },
			})()

			req.EventID = args[0]
			req.WorkflowDefinitionID = args[1]
			req.Configuration, err = getWorkflowPropertiesFlag(cmd.Flags())
			if err != nil {
				return nil, err
			}

			return s.Create(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(2)
	addWorkflowPropertiesFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func workflowGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get a Workflow",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Workflow
				req svc.WorkflowGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflow(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflow(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func workflowUpdateCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"update [id]",
		"Update a Workflow",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Workflow
				req svc.WorkflowUpdateRequest
				err error
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflow(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflow(extAPI) },
			})()

			req.ID = args[0]

			req.Configuration, err = getWorkflowPropertiesFlag(cmd.Flags())
			if err != nil {
				return nil, err
			}

			switch getWorkflowStatusFlag(cmd.Flags()) {
			case RunningWorkflowStatus:
				req.Status = new(api.RunningWorkflowStatus)
			case PausedWorkflowStatus:
				req.Status = new(api.PausedWorkflowStatus)
			case StoppedWorkflowStatus:
				req.Status = new(api.StoppedWorkflowStatus)
			}

			return s.Update(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	addWorkflowPropertiesFlag(cmd.Flags())
	addWorkflowStatusFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func workflowDeleteCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"delete [id]",
		"Delete a Workflow",
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Workflow
				req svc.WorkflowDeleteRequest
				err error
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastWorkflow(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesWorkflow(extAPI) },
			})()

			req.ID = args[0]

			err = s.Delete(cmd.Context(), req)
			return nil, err
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
