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

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	extapiv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11"
	oc "github.com/tales-media/cli/pkg/opencast/client"
)

type WithWorkflowOptions struct {
	WithOperations    bool
	WithConfiguration bool
}

var _ oc.RequestOpts = WithWorkflowOptions{}

func (opt WithWorkflowOptions) Apply(r *oc.Request) error {
	return r.ApplyOptions(
		oc.WithQuery("withoperations", strconv.FormatBool(opt.WithOperations)),
		oc.WithQuery("withconfiguration", strconv.FormatBool(opt.WithConfiguration)),
	)
}

func (c *client) CreateWorkflow(ctx context.Context, eventID, workflowDefinitionID string, configuration extapiv1.Properties, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) {
			return c.CreateWorkflowRequest(ctx, eventID, workflowDefinitionID, configuration, opts...)
		},
	)
}

func (c *client) CreateWorkflowRequest(ctx context.Context, eventID, workflowDefinitionID string, configuration extapiv1.Properties, opts ...oc.RequestOpts) (*oc.Request, error) {
	body := oc.NewMultipartBody()
	if err := body.WriteStringField("event_identifier", eventID); err != nil {
		return nil, err
	}
	if err := body.WriteStringField("workflow_definition_identifier", workflowDefinitionID); err != nil {
		return nil, err
	}
	if len(configuration) > 0 {
		cfgJSON, err := json.Marshal(configuration)
		if err != nil {
			return nil, err
		}
		if err := body.WriteField("configuration", cfgJSON); err != nil {
			return nil, err
		}
	}
	return oc.NewRequest(
		ctx,
		http.MethodPost,
		extapiv1.WorkflowInstancesServiceType,
		"/api/workflows/",
		body,
		opts...,
	)
}

func (c *client) GetWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) { return c.GetWorkflowRequest(ctx, id, opts...) },
	)
}

func (c *client) GetWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		extapiv1.WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}

func (c *client) UpdateWorkflow(ctx context.Context, id string, state *extapiv1.WorkflowState, configuration extapiv1.Properties, opts ...oc.RequestOpts) (*extapiv1.WorkflowInstance, *oc.Response, error) {
	return oc.GenericAutoDecodedDo[*extapiv1.WorkflowInstance](
		c,
		func() (*oc.Request, error) { return c.UpdateWorkflowRequest(ctx, id, state, configuration, opts...) },
	)
}

func (c *client) UpdateWorkflowRequest(ctx context.Context, id string, state *extapiv1.WorkflowState, configuration extapiv1.Properties, opts ...oc.RequestOpts) (*oc.Request, error) {
	body := oc.NewMultipartBody()
	if state != nil {
		if err := body.WriteStringField("state", string(*state)); err != nil {
			return nil, err
		}
	}
	if len(configuration) > 0 {
		cfgJSON, err := json.Marshal(configuration)
		if err != nil {
			return nil, err
		}
		if err := body.WriteField("configuration", cfgJSON); err != nil {
			return nil, err
		}
	}
	return oc.NewRequest(
		ctx,
		http.MethodPut,
		extapiv1.WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		body,
		opts...,
	)
}

func (c *client) DeleteWorkflow(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Response, error) {
	return oc.GenericDo(
		c,
		func() (*oc.Request, error) { return c.DeleteWorkflowRequest(ctx, id, opts...) },
	)
}

func (c *client) DeleteWorkflowRequest(ctx context.Context, id string, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodDelete,
		extapiv1.WorkflowInstancesServiceType,
		"/api/workflows/"+url.PathEscape(id),
		oc.NoBody,
		opts...,
	)
}
