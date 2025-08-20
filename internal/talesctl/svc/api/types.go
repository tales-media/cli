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

import (
	"encoding"
	"fmt"
	"time"
)

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

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
	LastUpdate *time.Time  `human:"Last Update" json:"lastUpdate" yaml:"update"`
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

type Workflow struct {
	ID                   int64             `human:"ID" json:"id" yaml:"id"`
	Title                string            `human:"Title" json:"title" yaml:"title"`
	Description          string            `human:"Description,wideonly" json:"description" yaml:"description"`
	WorkflowDefinitionID string            `human:"Workflow Definition ID" json:"workflowDefinitionId" yaml:"workflowDefinitionId"`
	EventID              string            `human:"Event ID" json:"eventId" yaml:"eventId"`
	Creator              string            `human:"Creator" json:"creator" yaml:"creator"`
	Status               WorkflowStatus    `human:"Status" json:"status" yaml:"status"`
	Operations           []Operation       `human:"Operations,wideonly" json:"operations" yaml:"operations"`
	Configuration        map[string]string `human:"Configuration,wideonly" json:"configuration" yaml:"configuration"`
}

type WorkflowStatus string

const (
	InstantiatedWorkflowStatus = WorkflowStatus("instantiated")
	RunningWorkflowStatus      = WorkflowStatus("running")
	StoppedWorkflowStatus      = WorkflowStatus("stopped")
	PausedWorkflowStatus       = WorkflowStatus("paused")
	SucceededWorkflowStatus    = WorkflowStatus("succeeded")
	FailedWorkflowStatus       = WorkflowStatus("failed")
	FailingWorkflowStatus      = WorkflowStatus("failing")
)

type Operation struct {
	ID                   *int64                  `human:"ID" json:"id" yaml:"id"`
	Operation            string                  `human:"Operation" json:"operation" yaml:"operation"`
	Description          string                  `human:"Description" json:"description" yaml:"description"`
	Status               WorkflowOperationStatus `human:"Status" json:"status" yaml:"status"`
	TimeInQueue          int64                   `human:"Time In Queue,wideonly" json:"timeInQueue" yaml:"timeInQueue"`
	Host                 string                  `human:"Host,wideonly" json:"host" yaml:"host"`
	If                   string                  `human:"If,wideonly" json:"if" yaml:"if"`
	FailWorkflowOnError  bool                    `human:"Fail Workflow On Error,wideonly" json:"failWorkflowOnError" yaml:"failWorkflowOnError"`
	ErrorHandlerWorkflow string                  `human:"Error Handler Workflow,wideonly" json:"errorHandlerWorkflow" yaml:"errorHandlerWorkflow"`
	RetryStrategy        WorkflowRetryStrategy   `human:"Retry Strategy,wideonly" json:"retryStrategy" yaml:"retryStrategy"`
	MaxAttempts          int64                   `human:"Max Attempts,wideonly" json:"maxAttempts" yaml:"maxAttempts"`
	FailedAttempts       int64                   `human:"Failed Attempts,wideonly" json:"failedAttempts" yaml:"failedAttempts"`
	Configuration        map[string]string       `human:"Configuration,wideonly" json:"configuration" yaml:"configuration"`
	Start                *time.Time              `human:"Start" json:"start" yaml:"start"`
	Completion           *time.Time              `human:"Completion" json:"completion" yaml:"completion"`
}

type WorkflowOperationStatus string

const (
	InstantiatedWorkflowOperationStatus = WorkflowOperationStatus("instantiated")
	RunningWorkflowOperationStatus      = WorkflowOperationStatus("running")
	PausedWorkflowOperationStatus       = WorkflowOperationStatus("paused")
	SucceededWorkflowOperationStatus    = WorkflowOperationStatus("succeeded")
	FailedWorkflowOperationStatus       = WorkflowOperationStatus("failed")
	SkippedWorkflowOperationStatus      = WorkflowOperationStatus("skipped")
	RetryWorkflowOperationStatus        = WorkflowOperationStatus("retry")
)

type WorkflowRetryStrategy string

const (
	NoneWorkflowRetryStrategy  = WorkflowRetryStrategy("none")
	RetryWorkflowRetryStrategy = WorkflowRetryStrategy("retry")
	HoldWorkflowRetryStrategy  = WorkflowRetryStrategy("hold")
)

type WorkflowDefinition struct {
	ID                     string                `human:"ID" json:"id" yaml:"id"`
	Title                  string                `human:"Title" json:"title" yaml:"title"`
	Description            string                `human:"Description,wideonly" json:"description" yaml:"description"`
	Tags                   []string              `human:"Tags,wideonly" json:"tags" yaml:"tags"`
	ConfigurationPanel     *string               `human:"Configuration Panel,wideonly" json:"configurationPanel" yaml:"configurationPanel"`
	ConfigurationPanelJSON *string               `human:"Configuration Panel JSON,wideonly" json:"configurationPanelJson" yaml:"configurationPanelJson"`
	Operations             []OperationDefinition `human:"Operations,wideonly" json:"operations" yaml:"operations"`
}

type OperationDefinition struct {
	Operation            string                `human:"Operation" json:"operation" yaml:"operation"`
	Description          string                `human:"Description" json:"description" yaml:"description"`
	Configuration        map[string]string     `human:"Configuration,wideonly" json:"configuration" yaml:"configuration"`
	If                   string                `human:"If,wideonly" json:"if" yaml:"if"`
	Unless               string                `human:"Unless,wideonly" json:"unless" yaml:"unless"`
	FailWorkflowOnError  bool                  `human:"Fail Workflow On Error,wideonly" json:"failWorkflowOnError" yaml:"failWorkflowOnError"`
	ErrorHandlerWorkflow string                `human:"Error Handler Workflow,wideonly" json:"errorHandlerWorkflow" yaml:"errorHandlerWorkflow"`
	RetryStrategy        WorkflowRetryStrategy `human:"Retry Strategy,wideonly" json:"retryStrategy" yaml:"retryStrategy"`
	MaxAttempts          int64                 `human:"Max Attempts,wideonly" json:"maxAttempts" yaml:"maxAttempts"`
}
