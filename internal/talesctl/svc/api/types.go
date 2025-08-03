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

package api

import "time"

type APIInfo struct {
	URL               string   `human:"URL" json:"url" yaml:"url"`
	DefaultVersion    string   `human:"Default Version" json:"defaultVersion" yaml:"defaultVersion"`
	SupportedVersions []string `human:"Supported Versions,wideonly" json:"supportedVersions" yaml:"supportedVersions"`
}

type Me struct {
	Username string   `human:"Username" json:"username" yaml:"username"`
	Name     string   `human:"Name" json:"name" yaml:"name"`
	Email    string   `human:"E-Mail" json:"email" yaml:"email"`
	UserRole string   `human:"User Role,wideonly" json:"userRole" yaml:"userRole"`
	Provider string   `human:"Provider,wideonly" json:"provider" yaml:"provider"`
	Roles    []string `human:"Roles,wideonly" json:"roles" yaml:"roles"`
}

type Organization struct {
	ID            string            `human:"ID" json:"id" yaml:"id"`
	Name          string            `human:"Name" json:"name" yaml:"name"`
	AdminRole     string            `human:"Admin Role" json:"adminRole" yaml:"adminRole"`
	AnonymousRole string            `human:"Anonymous Role" json:"anonymousRole" yaml:"anonymousRole"`
	Properties    map[string]string `human:"Properties,wideonly" json:"properties" yaml:"properties"`
}

type Agent struct {
	Name       string      `human:"Name" json:"name" yaml:"name"`
	URL        string      `human:"URL" json:"url" yaml:"url"`
	LastUpdate time.Time   `human:"Last Update" json:"lastUpdate" yaml:"update"`
	Status     AgentStatus `human:"Status" json:"status" yaml:"status"`
	Inputs     []string    `human:"Inputs,wideonly" json:"inputs" yaml:"inputs"`
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
	ID                     string                `human:"ID" json:"id" yaml:"id"`
	Title                  string                `human:"Title" json:"title" yaml:"title"`
	Description            string                `human:"Description,wideonly" json:"description" yaml:"description"`
	Tags                   []string              `human:"Tags,wideonly" json:"tags" yaml:"tags"`
	ConfigurationPanel     string                `human:"Configuration Panel,wideonly" json:"configuration_panel" yaml:"configuration_panel"`
	ConfigurationPanelJSON string                `human:"Configuration Panel JSON,wideonly" json:"configuration_panel_json" yaml:"configuration_panel_json"`
	Operations             []OperationDefinition `human:"Operations,wideonly" json:"operations" yaml:"operations"`
}

type OperationDefinition struct {
	Operation            string            `human:"Operation" json:"operation" yaml:"operation"`
	Description          string            `human:"Description" json:"description" yaml:"description"`
	Configuration        map[string]string `human:"Configuration,wideonly" json:"configuration" yaml:"configuration"`
	If                   string            `human:"If,wideonly" json:"if" yaml:"if"`
	Unless               string            `human:"Unless,wideonly" json:"unless" yaml:"unless"`
	FailWorkflowOnError  bool              `human:"Fail Workflow On Error,wideonly" json:"fail_workflow_on_error" yaml:"fail_workflow_on_error"`
	ErrorHandlerWorkflow string            `human:"Error Handler Workflow,wideonly" json:"error_handler_workflow" yaml:"error_handler_workflow"`
	RetryStrategy        RetryStrategy     `human:"Retry Strategy,wideonly" json:"retry_strategy" yaml:"retry_strategy"`
	MaxAttempts          int               `human:"Max Attempts,wideonly" json:"max_attempts" yaml:"max_attempts"`
}

type RetryStrategy string

const (
	NoneRetryStrategy  = RetryStrategy("none")
	RetryRetryStrategy = RetryStrategy("retry")
	HoldRetryStrategy  = RetryStrategy("hold")
)
