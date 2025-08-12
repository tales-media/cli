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
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/utils/ptr"

	"github.com/tales-media/cli/internal/talesctl/svc"
	"github.com/tales-media/cli/internal/talesctl/svc/api"
	extapiclientv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11/client"
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

			req.Configuration = make(map[string]string)
			cfgFlags, err := cmd.Flags().GetStringSlice(ConfigFlag)
			if err != nil {
				return nil, err
			}
			for _, c := range cfgFlags {
				k, v, ok := strings.Cut(c, "=")
				if !ok {
					return nil, errors.New("invalid config flag: use key=value syntax")
				}
				req.Configuration[k] = v
			}

			return s.Create(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(2)
	cmd.Flags().StringSliceP(ConfigFlag, ConfigFlagShort, nil, "set workflow configuration property in the form \"key=value\" (can be specified multiple times or a comma separated list)")
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

			req.Configuration = make(map[string]string)
			cfgFlags, err := cmd.Flags().GetStringSlice(ConfigFlag)
			if err != nil {
				return nil, err
			}
			for _, c := range cfgFlags {
				k, v, ok := strings.Cut(c, "=")
				if !ok {
					return nil, errors.New("invalid config flag: use key=value syntax")
				}
				req.Configuration[k] = v
			}

			workflowState := cmd.Flag(StateFlag).Value.(*mapValue[WorkflowState])
			switch workflowState.Value() {
			case RunningWorkflowState:
				req.State = ptr.To(api.RunningWorkflowState)
			case PausedWorkflowState:
				req.State = ptr.To(api.PausedWorkflowState)
			case StoppedWorkflowState:
				req.State = ptr.To(api.StoppedWorkflowState)
			}

			return s.Update(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	workflowStateValue := WorkflowStateValue()
	cmd.Flags().Var(workflowStateValue, StateFlag, workflowStateValue.Usage("set a new workflow state"))
	cmd.Flags().StringSliceP(ConfigFlag, ConfigFlagShort, nil, "set workflow configuration property in the form \"key=value\" (can be specified multiple times or a comma separated list)")
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
	return cmd
}
