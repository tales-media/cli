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

package v1_11

import "time"

const Version = "v1.11.0"

const (
	ServiceType                    = "org.opencastproject.external"
	AgentsServiceType              = "org.opencastproject.external.agents"
	EventsServiceType              = "org.opencastproject.external.events"
	GroupsServiceType              = "org.opencastproject.external.groups"
	ListprovidersServiceType       = "org.opencastproject.external.listproviders"
	PlaylistsServiceType           = "org.opencastproject.external.playlists"
	SecurityServiceType            = "org.opencastproject.external.security"
	SeriesServiceType              = "org.opencastproject.external.series"
	StatisticsServiceType          = "org.opencastproject.external.statistics"
	WorkflowDefinitionsServiceType = "org.opencastproject.external.workflows.definitions"
	WorkflowInstancesServiceType   = "org.opencastproject.external.workflows.instances"
)

type Properties map[string]string

type DateTime = time.Time

type API struct {
	Version string `json:"version,omitempty"`
	URL     string `json:"url,omitempty"`
}

type APIVersion struct {
	Default  string   `json:"default,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type Organization struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AdminRole     string `json:"adminRole"`
	AnonymousRole string `json:"anonymousRole"`
}

type Me struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserRole string `json:"userrole"`
	Provider string `json:"provider"`
}

type StringList []string

type Agent struct {
	AgentID string      `json:"agent_id"`
	Inputs  []string    `json:"inputs"`
	Update  DateTime    `json:"update"`
	URL     string      `json:"url"`
	Status  AgentStatus `json:"status"`
}

type AgentStatus string

const (
	UnknownAgentStatus      = AgentStatus("unknown")
	IdleAgentStatus         = AgentStatus("idle")
	CapturingAgentStatus    = AgentStatus("capturing")
	UploadingAgentStatus    = AgentStatus("uploading")
	ShuttingDownAgentStatus = AgentStatus("shutting_down")
	OfflineAgentStatus      = AgentStatus("offline")
	ErrorAgentStatus        = AgentStatus("error")
)

type WorkflowDefinition struct {
	Identifier             string                `json:"identifier"`
	Title                  string                `json:"title"`
	Description            string                `json:"description"`
	Tags                   []string              `json:"tags"`
	ConfigurationPanel     string                `json:"configuration_panel"`
	ConfigurationPanelJSON string                `json:"configuration_panel_json"`
	Operations             []OperationDefinition `json:"operations"`
}

type OperationDefinition struct {
	Operation            string            `json:"operation"`
	Description          string            `json:"description"`
	Configuration        map[string]string `json:"configuration"`
	If                   string            `json:"if"`
	Unless               string            `json:"unless"`
	FailWorkflowOnError  bool              `json:"fail_workflow_on_error"`
	ErrorHandlerWorkflow string            `json:"error_handler_workflow"`
	RetryStrategy        RetryStrategy     `json:"retry_strategy"`
	MaxAttempts          int               `json:"max_attempts"`
}

type RetryStrategy string

const (
	NoneRetryStrategy  = RetryStrategy("none")
	RetryRetryStrategy = RetryStrategy("retry")
	HoldRetryStrategy  = RetryStrategy("hold")
)
