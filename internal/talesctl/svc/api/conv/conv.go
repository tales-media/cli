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

package conv

import (
	"fmt"
	"strings"
	"time"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	extapiv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
	"shio.solutions/tales.media/cli/pkg/opencast/apis/meta/base"
)

func Map[F, T any](ocList []F, mapper func(F) T) []T {
	list := make([]T, 0, len(ocList))
	for _, item := range ocList {
		list = append(list, mapper(item))
	}
	return list
}

func SortDirectionToOCSortDirection(sortDirection api.SortDirection) extapiclientv1.SortDirection {
	switch sortDirection {
	case api.Ascending:
		return extapiclientv1.Ascending
	case api.Descending:
		return extapiclientv1.Descending
	}
	panic(fmt.Sprintf("BUG: unknown SortDirection '%d'", sortDirection))
}

func OCAPIToAPIInfo(ocAPI extapiv1.API, ocVersion extapiv1.APIVersion) api.APIInfo {
	return api.APIInfo{
		URL:               ocAPI.URL,
		DefaultVersion:    ocAPI.Version,
		SupportedVersions: ocVersion.Versions,
	}
}

func OCMeToMe(ocMe extapiv1.Me, roles []string) api.Me {
	return api.Me{
		Username: ocMe.Username,
		Name:     ocMe.Name,
		Email:    ocMe.Email,
		UserRole: ocMe.UserRole,
		Provider: ocMe.Provider,
		Roles:    roles,
	}
}

func OCOrganizationToOrganization(ocOrganization extapiv1.Organization, ocProperties base.Properties) api.Organization {
	return api.Organization{
		ID:            ocOrganization.ID,
		Name:          ocOrganization.Name,
		AdminRole:     ocOrganization.AdminRole,
		AnonymousRole: ocOrganization.AnonymousRole,
		Properties:    ocProperties,
	}
}

func OCAgentToAgent(ocAgent extapiv1.Agent) api.Agent {
	a := api.Agent{
		Name:   ocAgent.AgentID,
		URL:    ocAgent.URL,
		Status: OCAgentStatusToAgentStatus(ocAgent.Status),
		Inputs: ocAgent.Inputs,
	}
	if !ocAgent.Update.IsZero() {
		a.LastUpdate = new(ocAgent.Update.Time)
	}
	return a
}

func OCAgentStatusToAgentStatus(ocAgentStatus extapiv1.AgentStatus) api.AgentStatus {
	if ocAgentStatus == "" {
		return ""
	}
	switch ocAgentStatus {
	case extapiv1.UnknownAgentStatus:
		return api.UnknownAgentStatus
	case extapiv1.IdleAgentStatus:
		return api.IdleAgentStatus
	case extapiv1.CapturingAgentStatus:
		return api.CapturingAgentStatus
	case extapiv1.UploadingAgentStatus:
		return api.UploadingAgentStatus
	case extapiv1.ShuttingDownAgentStatus:
		return api.ShuttingDownAgentStatus
	case extapiv1.OfflineAgentStatus:
		return api.OfflineAgentStatus
	case extapiv1.ErrorAgentStatus:
		return api.ErrorAgentStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast AgentStatus '%s'", ocAgentStatus))
}

func OCWorkflowInstanceToWorkflow(ocWorkflowInstance extapiv1.WorkflowInstance) api.Workflow {
	return api.Workflow{
		ID:                   int64(ocWorkflowInstance.Identifier),
		Title:                ocWorkflowInstance.Title,
		Description:          ocWorkflowInstance.Description,
		WorkflowDefinitionID: ocWorkflowInstance.WorkflowDefinitionIdentifier,
		EventID:              ocWorkflowInstance.EventIdentifier,
		Creator:              ocWorkflowInstance.Creator,
		Status:               OCWorkflowStateToWorkflowStatus(ocWorkflowInstance.State),
		Operations:           Map(ocWorkflowInstance.Operations, OCOperationInstanceToOperation),
		Configuration:        ocWorkflowInstance.Configuration,
	}
}

func OCWorkflowStateToWorkflowStatus(ocWorkflowState extapiv1.WorkflowState) api.WorkflowStatus {
	if ocWorkflowState == "" {
		return ""
	}
	switch ocWorkflowState {
	case extapiv1.InstantiatedWorkflowState:
		return api.InstantiatedWorkflowStatus
	case extapiv1.RunningWorkflowState:
		return api.RunningWorkflowStatus
	case extapiv1.StoppedWorkflowState:
		return api.StoppedWorkflowStatus
	case extapiv1.PausedWorkflowState:
		return api.PausedWorkflowStatus
	case extapiv1.SucceededWorkflowState:
		return api.SucceededWorkflowStatus
	case extapiv1.FailedWorkflowState:
		return api.FailedWorkflowStatus
	case extapiv1.FailingWorkflowState:
		return api.FailingWorkflowStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowState '%s'", ocWorkflowState))
}

func WorkflowStatusToOCWorkflowState(workflowState api.WorkflowStatus) extapiv1.WorkflowState {
	if workflowState == "" {
		return ""
	}
	switch workflowState {
	case api.InstantiatedWorkflowStatus:
		return extapiv1.InstantiatedWorkflowState
	case api.RunningWorkflowStatus:
		return extapiv1.RunningWorkflowState
	case api.StoppedWorkflowStatus:
		return extapiv1.StoppedWorkflowState
	case api.PausedWorkflowStatus:
		return extapiv1.PausedWorkflowState
	case api.SucceededWorkflowStatus:
		return extapiv1.SucceededWorkflowState
	case api.FailedWorkflowStatus:
		return extapiv1.FailedWorkflowState
	case api.FailingWorkflowStatus:
		return extapiv1.FailingWorkflowState
	}
	panic(fmt.Sprintf("BUG: unknown WorkflowState '%s'", workflowState))
}

func OCOperationInstanceToOperation(ocOperationInstance extapiv1.OperationInstance) api.Operation {
	o := api.Operation{
		ID:                   (*int64)(ocOperationInstance.Identifier),
		Operation:            ocOperationInstance.Operation,
		Description:          ocOperationInstance.Description,
		Status:               OCWorkflowOperationStateToWorkflowOperationStatus(ocOperationInstance.State),
		TimeInQueue:          int64(ocOperationInstance.TimeInQueue),
		Host:                 ocOperationInstance.Host,
		If:                   ocOperationInstance.If,
		FailWorkflowOnError:  ocOperationInstance.FailWorkflowOnError,
		ErrorHandlerWorkflow: ocOperationInstance.ErrorHandlerWorkflow,
		RetryStrategy:        OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocOperationInstance.RetryStrategy),
		MaxAttempts:          int64(ocOperationInstance.MaxAttempts),
		FailedAttempts:       int64(ocOperationInstance.FailedAttempts),
		Configuration:        ocOperationInstance.Configuration,
	}
	if !ocOperationInstance.Start.IsZero() {
		o.Start = new(ocOperationInstance.Start.Time)
	}
	if !ocOperationInstance.Completion.IsZero() {
		o.Completion = new(ocOperationInstance.Completion.Time)
	}
	return o
}

func OCWorkflowOperationStateToWorkflowOperationStatus(ocWorkflowOperationState extapiv1.WorkflowOperationState) api.WorkflowOperationStatus {
	if ocWorkflowOperationState == "" {
		return ""
	}
	switch ocWorkflowOperationState {
	case extapiv1.InstantiatedWorkflowOperationState:
		return api.InstantiatedWorkflowOperationStatus
	case extapiv1.RunningWorkflowOperationState:
		return api.RunningWorkflowOperationStatus
	case extapiv1.PausedWorkflowOperationState:
		return api.PausedWorkflowOperationStatus
	case extapiv1.SucceededWorkflowOperationState:
		return api.SucceededWorkflowOperationStatus
	case extapiv1.FailedWorkflowOperationState:
		return api.FailedWorkflowOperationStatus
	case extapiv1.SkippedWorkflowOperationState:
		return api.SkippedWorkflowOperationStatus
	case extapiv1.RetryWorkflowOperationState:
		return api.RetryWorkflowOperationStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowOperationState '%s'", ocWorkflowOperationState))
}

func OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocWorkflowRetryStrategy extapiv1.WorkflowRetryStrategy) api.WorkflowRetryStrategy {
	if ocWorkflowRetryStrategy == "" {
		return ""
	}
	switch ocWorkflowRetryStrategy {
	case extapiv1.NoneWorkflowRetryStrategy:
		return api.NoneWorkflowRetryStrategy
	case extapiv1.RetryWorkflowRetryStrategy:
		return api.RetryWorkflowRetryStrategy
	case extapiv1.HoldWorkflowRetryStrategy:
		return api.HoldWorkflowRetryStrategy
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowRetryStrategy '%s'", ocWorkflowRetryStrategy))
}

func OCWorkflowDefinitionToWorkflowDefinition(ocWorkflowDefinition extapiv1.WorkflowDefinition) api.WorkflowDefinition {
	return api.WorkflowDefinition{
		ID:                     ocWorkflowDefinition.Identifier,
		Title:                  ocWorkflowDefinition.Title,
		Description:            ocWorkflowDefinition.Description,
		Tags:                   ocWorkflowDefinition.Tags,
		ConfigurationPanel:     ocWorkflowDefinition.ConfigurationPanel,
		ConfigurationPanelJSON: ocWorkflowDefinition.ConfigurationPanelJSON,
		Operations:             Map(ocWorkflowDefinition.Operations, OCOperationDefinitionToOperationDefinition),
	}
}

func OCOperationDefinitionToOperationDefinition(ocOperationDefinition extapiv1.OperationDefinition) api.OperationDefinition {
	return api.OperationDefinition{
		Operation:            ocOperationDefinition.Operation,
		Description:          ocOperationDefinition.Description,
		Configuration:        ocOperationDefinition.Configuration,
		If:                   ocOperationDefinition.If,
		Unless:               ocOperationDefinition.Unless,
		FailWorkflowOnError:  ocOperationDefinition.FailWorkflowOnError,
		ErrorHandlerWorkflow: ocOperationDefinition.ErrorHandlerWorkflow,
		RetryStrategy:        OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocOperationDefinition.RetryStrategy),
		MaxAttempts:          int64(ocOperationDefinition.MaxAttempts),
	}
}
