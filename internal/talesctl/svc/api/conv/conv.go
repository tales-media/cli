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

func OCAgentToAgent(ocAgent extapiv1.Agent) api.Agent {
	return api.Agent{
		Name:       ocAgent.AgentID,
		URL:        ocAgent.URL,
		LastUpdate: ocAgent.Update,
		Status:     api.AgentStatus(ocAgent.Status),
		Inputs:     ocAgent.Inputs,
	}
}

func Map[F, T any](ocList []F, mapper func(F) T) []T {
	list := make([]T, 0, len(ocList))
	for _, item := range ocList {
		list = append(list, mapper(item))
	}
	return list
}
