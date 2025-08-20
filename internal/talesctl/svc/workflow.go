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

	"k8s.io/utils/ptr"

	"github.com/tales-media/cli/internal/talesctl/svc/api"
	"github.com/tales-media/cli/internal/talesctl/svc/api/conv"
	extapiv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

type Workflow interface {
	Create(context.Context, WorkflowCreateRequest) (api.Workflow, error)
	Get(context.Context, WorkflowGetRequest) (api.Workflow, error)
	Update(context.Context, WorkflowUpdateRequest) (api.Workflow, error)
	Delete(context.Context, WorkflowDeleteRequest) error
}

type WorkflowCreateRequest struct {
	EventID              string
	WorkflowDefinitionID string
	Configuration        map[string]string
}

type WorkflowGetRequest struct {
	ID string
}

type WorkflowUpdateRequest struct {
	ID            string
	Status        *api.WorkflowStatus
	Configuration map[string]string
}

type WorkflowDeleteRequest struct {
	ID string
}

type opencastWorkflow struct {
	extAPI extapiclientv1.Client
}

var _ Workflow = &opencastWorkflow{}

func NewOpencastWorkflow(extAPI extapiclientv1.Client) Workflow {
	return &opencastWorkflow{
		extAPI: extAPI,
	}
}

func (svc *opencastWorkflow) Create(ctx context.Context, req WorkflowCreateRequest) (api.Workflow, error) {
	ocWorkflow, _, err := svc.extAPI.CreateWorkflow(
		ctx,
		&extapiclientv1.CreateWorkflowRequestBody{
			EventID:              req.EventID,
			WorkflowDefinitionID: req.WorkflowDefinitionID,
			Configuration:        req.Configuration,
		},
		extapiclientv1.WithWorkflowOptions{
			WithOperations:    true,
			WithConfiguration: true,
		},
	)
	if err != nil {
		return api.Workflow{}, err
	}
	return conv.OCWorkflowInstanceToWorkflow(*ocWorkflow), nil
}

func (svc *opencastWorkflow) Get(ctx context.Context, req WorkflowGetRequest) (api.Workflow, error) {
	ocWorkflow, _, err := svc.extAPI.GetWorkflow(
		ctx,
		req.ID,
		extapiclientv1.WithWorkflowOptions{
			WithOperations:    true,
			WithConfiguration: true,
		},
	)
	if err != nil {
		return api.Workflow{}, err
	}
	return conv.OCWorkflowInstanceToWorkflow(*ocWorkflow), nil
}

func (svc *opencastWorkflow) Update(ctx context.Context, req WorkflowUpdateRequest) (api.Workflow, error) {
	var ocState *extapiv1.WorkflowState
	if req.Status != nil {
		ocState = ptr.To(conv.WorkflowStatusToOCWorkflowState(*req.Status))
	}
	ocWorkflow, _, err := svc.extAPI.UpdateWorkflow(
		ctx,
		req.ID,
		&extapiclientv1.UpdateWorkflowRequestBody{
			State:         ocState,
			Configuration: req.Configuration,
		},
		extapiclientv1.WithWorkflowOptions{
			WithOperations:    true,
			WithConfiguration: true,
		},
	)
	if err != nil {
		return api.Workflow{}, err
	}
	return conv.OCWorkflowInstanceToWorkflow(*ocWorkflow), nil
}

func (svc *opencastWorkflow) Delete(ctx context.Context, req WorkflowDeleteRequest) error {
	_, err := svc.extAPI.DeleteWorkflow(ctx, req.ID)
	return err
}

type talesWorkflow = opencastWorkflow

var _ Workflow = &talesWorkflow{}

var NewTalesWorkflow = NewOpencastWorkflow
