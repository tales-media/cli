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
	"github.com/tales-media/cli/internal/talesctl/svc/api"
	extapiv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11"
)

func Map[F, T any](ocList []F, mapper func(F) T) []T {
	list := make([]T, 0, len(ocList))
	for _, item := range ocList {
		list = append(list, mapper(item))
	}
	return list
}

func OCAPIToAPIInfo(ocAPI extapiv1.API, ocVersion extapiv1.APIVersion) api.APIInfo {
	return api.APIInfo{
		URL:               ocAPI.URL,
		DefaultVersion:    ocAPI.Version,
		SupportedVersions: ocVersion.Versions,
	}
}

func OCMeToMe(ocMe extapiv1.Me, roles extapiv1.StringList) api.Me {
	return api.Me{
		Username: ocMe.Username,
		Name:     ocMe.Name,
		Email:    ocMe.Email,
		UserRole: ocMe.UserRole,
		Provider: ocMe.Provider,
		Roles:    roles,
	}
}

func OCOrganizationToOrganization(ocOrganization extapiv1.Organization, ocProperties extapiv1.OrganizationProperties) api.Organization {
	return api.Organization{
		ID:            ocOrganization.ID,
		Name:          ocOrganization.Name,
		AdminRole:     ocOrganization.AdminRole,
		AnonymousRole: ocOrganization.AnonymousRole,
		Properties:    ocProperties,
	}
}

func OCAgentToAgent(ocAgent extapiv1.Agent) api.Agent {
	return api.Agent{
		Name:       ocAgent.AgentID,
		URL:        ocAgent.URL,
		LastUpdate: ocAgent.Update,
		Status:     api.AgentStatus(ocAgent.Status),
		Inputs:     ocAgent.Inputs,
	}
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
		RetryStrategy:        api.RetryStrategy(ocOperationDefinition.RetryStrategy),
		MaxAttempts:          ocOperationDefinition.MaxAttempts,
	}
}
