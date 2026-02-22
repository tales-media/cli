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

package svc

import (
	"context"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"
	oc "shio.solutions/tales.media/opencast-client-go/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
)

type WorkflowDefinition interface {
	List(context.Context, WorkflowDefinitionListRequest) ([]api.WorkflowDefinition, error)
	Get(context.Context, WorkflowDefinitionGetRequest) (api.WorkflowDefinition, error)
}

type WorkflowDefinitionListRequest struct {
	FilterByTag   string
	SortBy        string // TODO: add conversion function
	SortDirection api.SortDirection
}

type WorkflowDefinitionGetRequest struct {
	ID string
}

type opencastWorkflowDefinition struct {
	extAPI extapiclientv1.Client
}

var _ WorkflowDefinition = &opencastWorkflowDefinition{}

func NewOpencastWorkflowDefinition(extAPI extapiclientv1.Client) WorkflowDefinition {
	return &opencastWorkflowDefinition{
		extAPI: extAPI,
	}
}

func (svc *opencastWorkflowDefinition) List(ctx context.Context, req WorkflowDefinitionListRequest) ([]api.WorkflowDefinition, error) {
	commonReqOpts := []oc.RequestOpts{
		extapiclientv1.WithWorkflowDefinitionOptions{
			WithOperations:             true,
			WithConfigurationPanel:     true,
			WithConfigurationPanelJSON: true,
		},
	}
	if req.FilterByTag != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.WorkflowDefinitionTagFilterKey: req.FilterByTag},
		)
	}
	if req.SortBy != "" {
		commonReqOpts = append(commonReqOpts, extapiclientv1.WithSort{{
			By:        extapiclientv1.SortKey(req.SortBy),
			Direction: conv.SortDirectionToOCSortDirection(req.SortDirection),
		}})
	}

	var ocWorkflowDefinitions []extapiv1.WorkflowDefinition
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			ocReq, err := svc.extAPI.ListWorkflowDefinitionRequest(
				ctx,
				extapiclientv1.WithPagination{
					Limit:  PageSize,
					Offset: i * PageSize,
				},
			)
			if err != nil {
				return nil, err
			}
			if err = ocReq.ApplyOptions(commonReqOpts...); err != nil {
				return nil, err
			}
			return ocReq, nil
		},
		oc.CollectAllPages(&ocWorkflowDefinitions),
	)
	if err != nil {
		return nil, err
	}

	workflowDefinitions := conv.Map(ocWorkflowDefinitions, conv.OCWorkflowDefinitionToWorkflowDefinition)
	return workflowDefinitions, nil
}

func (svc *opencastWorkflowDefinition) Get(ctx context.Context, req WorkflowDefinitionGetRequest) (api.WorkflowDefinition, error) {
	ocWorkflowDefinition, _, err := svc.extAPI.GetWorkflowDefinition(
		ctx,
		req.ID,
		extapiclientv1.WithWorkflowDefinitionOptions{
			WithOperations:             true,
			WithConfigurationPanel:     true,
			WithConfigurationPanelJSON: true,
		},
	)
	if err != nil {
		return api.WorkflowDefinition{}, err
	}
	return conv.OCWorkflowDefinitionToWorkflowDefinition(*ocWorkflowDefinition), nil
}

type talesWorkflowDefinition = opencastWorkflowDefinition

var _ WorkflowDefinition = &talesWorkflowDefinition{}

var NewTalesWorkflowDefinition = NewOpencastWorkflowDefinition
