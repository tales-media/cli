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

	"github.com/tales-media/cli/internal/talesctl/svc/api"
	"github.com/tales-media/cli/internal/talesctl/svc/api/conv"
	extapiv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11/client"
	oc "github.com/tales-media/cli/pkg/opencast/client"
)

type Agent interface {
	List(context.Context) ([]api.Agent, error)
	Get(context.Context, string) (api.Agent, error)
}

type opencastAgent struct {
	extAPI extapiclientv1.Client
}

var _ Agent = &opencastAgent{}

func NewOpencastAgent(extAPI extapiclientv1.Client) Agent {
	return &opencastAgent{
		extAPI: extAPI,
	}
}

func (svc *opencastAgent) List(ctx context.Context) ([]api.Agent, error) {
	var ocAgents []extapiv1.Agent
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			return svc.extAPI.ListAgentRequest(
				ctx,
				extapiclientv1.WithLimit(PageSize),
				extapiclientv1.WithOffset(i*PageSize),
			)
		},
		oc.CollectAllPages(&ocAgents),
	)
	if err != nil {
		return nil, err
	}
	agents := conv.Map(ocAgents, conv.OCAgentToAgent)
	return agents, nil
}

func (svc *opencastAgent) Get(ctx context.Context, name string) (api.Agent, error) {
	ocAgent, _, err := svc.extAPI.GetAgent(ctx, name)
	if err != nil {
		return api.Agent{}, err
	}
	return conv.OCAgentToAgent(*ocAgent), nil
}

type talesAgent = opencastAgent

var _ Agent = &talesAgent{}

var NewTalesAgent = NewOpencastAgent
