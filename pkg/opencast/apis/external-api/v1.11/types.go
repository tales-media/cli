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

import (
	"encoding/json"
	"time"

	"github.com/tales-media/cli/pkg/opencast/apis/meta/strobj"
)

const Version = "v1.11.0"

const (
	ServiceType                    = "org.opencastproject.external"
	AgentsServiceType              = "org.opencastproject.external.agents"
	EventsServiceType              = "org.opencastproject.external.events"
	GroupsServiceType              = "org.opencastproject.external.groups"
	ListProvidersServiceType       = "org.opencastproject.external.listproviders"
	PlaylistsServiceType           = "org.opencastproject.external.playlists"
	SecurityServiceType            = "org.opencastproject.external.security"
	SeriesServiceType              = "org.opencastproject.external" // TODO: fix this in Opencast
	StatisticsServiceType          = "org.opencastproject.external.statistics"
	WorkflowDefinitionsServiceType = "org.opencastproject.external.workflows.definitions"
	WorkflowInstancesServiceType   = "org.opencastproject.external.workflows.instances"
)

type Properties map[string]string

type DateTime time.Time

func (dt DateTime) IsZero() bool {
	return time.Time(dt).IsZero()
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte{'"', '"'}, nil
	}
	return time.Time(dt).MarshalJSON()
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 2 && data[0] == '"' && data[1] == '"' {
		// Opencast represents null dates as blank string. Skip unmarshal and let dt remain as zero value.
		return nil
	}
	return (*time.Time)(dt).UnmarshalJSON(data)
}

type Flavor string

const (
	DublinCoreEpisodeFlavor = Flavor("dublincore/episode")
	DublinCoreSeriesFlavor  = Flavor("dublincore/series")
)

// Access Control List (ACL)
type ACL []ACE

// Access Control Entry (ACE)
type ACE struct {
	Allow  bool   `json:"allow,omitempty"`
	Action Action `json:"action,omitempty"`
	Role   string `json:"role,omitempty"`
}

type Action string

const (
	ReadAction  = Action("read")
	WriteAction = Action("write")
)

type Catalog struct {
	Label  string  `json:"label,omitempty"`
	Flavor Flavor  `json:"flavor,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}

type Value struct {
	ID       string          `json:"id,omitempty"`
	RawValue json.RawMessage `json:"value,omitempty"` // TODO: this can be many types depending on the "type" field
}

type Field struct {
	ID              string                                    `json:"id,omitempty"`
	Label           string                                    `json:"label,omitempty"`
	RawValue        json.RawMessage                           `json:"value,omitempty"` // TODO: this can be many types depending on the "type" field
	Type            FieldType                                 `json:"type,omitempty"`
	ReadOnly        bool                                      `json:"readOnly,omitempty"`
	Required        bool                                      `json:"required,omitempty"`
	Collection      *strobj.StringOrObject[map[string]string] `json:"collection,omitempty"`
	Translatable    *bool                                     `json:"translatable,omitempty"`
	Delimiter       *string                                   `json:"delimiter,omitempty"`
	DifferentValues *bool                                     `json:"differentValues,omitempty"`
}

type FieldType string

const (
	BooleanFieldType     = FieldType("boolean")
	DateFieldType        = FieldType("date")
	MixedTextFieldType   = FieldType("mixed_text")
	NumberFieldType      = FieldType("number")
	OrderedTextFieldType = FieldType("ordered_text")
	TextFieldType        = FieldType("text")
	TextLongFieldType    = FieldType("text_long")
	TimeFieldType        = FieldType("time")
)

type API struct {
	Version string `json:"version,omitempty"`
	URL     string `json:"url,omitempty"`
}

type APIVersion struct {
	Default  string   `json:"default,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

type Organization struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	AdminRole     string `json:"adminRole,omitempty"`
	AnonymousRole string `json:"anonymousRole,omitempty"`
}

type Me struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	UserRole string `json:"userrole,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type StringList []string

type SignedURL struct {
	Error      string   `json:"error,omitempty"`
	URL        string   `json:"url,omitempty"`
	ValidUntil DateTime `json:"valid-until,omitempty"`
}

const (
	SignURLError           = "Error while signing url"
	URLCannotBeSignedError = "Given URL cannot be signed"
)

type Group struct {
	Identifier   string `json:"identifier,omitempty"`
	Organization string `json:"organization,omitempty"`
	Role         string `json:"role,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Roles        string `json:"roles,omitempty"`
	Members      string `json:"members,omitempty"`
}

type StatisticProvider struct {
	Identifier   string                       `json:"identifier,omitempty"`
	Title        string                       `json:"title,omitempty"`
	Description  string                       `json:"description,omitempty"`
	Type         StatisticProviderType        `json:"type,omitempty"`
	ResourceType StatisticResourceType        `json:"resourceType,omitempty"`
	Parameters   []StatisticProviderParameter `json:"parameters,omitempty"`
}

type StatisticProviderType string

const (
	UnknownStatisticProviderType    = StatisticProviderType("unknown")
	TimeSeriesStatisticProviderType = StatisticProviderType("timeseries")
)

type StatisticResourceType string

const (
	UnknownStatisticResourceType      = StatisticResourceType("unknown")
	EpisodeStatisticResourceType      = StatisticResourceType("episode")
	SeriesStatisticResourceType       = StatisticResourceType("series")
	OrganizationStatisticResourceType = StatisticResourceType("organization")
)

type StatisticProviderParameter struct {
	Name     string                         `json:"name,omitempty"`
	Type     StatisticProviderParameterType `json:"type,omitempty"`
	Optional bool                           `json:"optional,omitempty"`
	Values   []string                       `json:"values,omitempty"`
}

const (
	ResourceIDStatisticProviderParameter     = "resourceId"
	FromStatisticProviderParameter           = "from"
	ToStatisticProviderParameter             = "to"
	DataResolutionStatisticProviderParameter = "dataResolution"
	DetailLevelStatisticProviderParameter    = "detailLevel"
)

type StatisticProviderParameterType string

const (
	DatetimeStatisticProviderParameterType    = StatisticProviderParameterType("datetime")
	EnumerationStatisticProviderParameterType = StatisticProviderParameterType("enumeration")
	StringStatisticProviderParameterType      = StatisticProviderParameterType("string")
)

type StatisticQuery struct {
	Provider   Identifier `json:"provider,omitempty"`
	Parameters Properties `json:"parameters,omitempty"`
}

type StatisticQueryResult struct {
	Provider   StatisticProvider `json:"provider,omitempty"`
	Parameters Properties        `json:"parameters,omitempty"`
	RawData    json.RawMessage   `json:"data,omitempty"`
}

type StatisticQueryResultTimeSeriesData struct {
	Labels []DateTime `json:"labels,omitempty"`
	Values []float64  `json:"values,omitempty"`
	Total  *float64   `json:"total,omitempty"`
}

type Agent struct {
	AgentID string      `json:"agent_id,omitempty"`
	Inputs  []string    `json:"inputs,omitempty"`
	Update  DateTime    `json:"update,omitempty"`
	URL     string      `json:"url,omitempty"`
	Status  AgentStatus `json:"status,omitempty"`
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

type Identifier struct {
	Identifier string `json:"identifier,omitempty"`
}

type Event struct {
	ArchiveVersion    *int64          `json:"archive_version,omitempty"`
	Created           DateTime        `json:"created,omitempty"`
	Creator           string          `json:"creator,omitempty"`
	Contributor       []string        `json:"contributor,omitempty"`
	Description       string          `json:"description,omitempty"`
	HasPreviews       bool            `json:"has_previews,omitempty"`
	Identifier        string          `json:"identifier,omitempty"`
	Location          string          `json:"location,omitempty"`
	Presenter         []string        `json:"presenter,omitempty"`
	Language          string          `json:"language,omitempty"`
	RightsHolder      string          `json:"rightsholder,omitempty"`
	License           string          `json:"license,omitempty"`
	IsPartOf          string          `json:"is_part_of,omitempty"`
	Series            string          `json:"series,omitempty"`
	Source            string          `json:"source,omitempty"`
	Status            EventStatus     `json:"status,omitempty"`
	PublicationStatus []string        `json:"publication_status,omitempty"`
	ProcessingState   ProcessingState `json:"processing_state,omitempty"`
	Start             DateTime        `json:"start,omitempty"`
	Duration          *int64          `json:"duration,omitempty"`
	Subjects          []string        `json:"subjects,omitempty"`
	Title             string          `json:"title,omitempty"`
	ACL               ACL             `json:"acl,omitempty"`
	Metadata          []Catalog       `json:"metadata,omitempty"`
	Scheduling        Scheduling      `json:"scheduling,omitempty"`
	Publications      []Publication   `json:"publications,omitempty"`
}

type EventStatus string

const (
	IngestingEventStatus           = EventStatus("EVENTS.EVENTS.STATUS.INGESTING")
	PausedEventStatus              = EventStatus("EVENTS.EVENTS.STATUS.PAUSED")
	PendingEventStatus             = EventStatus("EVENTS.EVENTS.STATUS.PENDING")
	ProcessedEventStatus           = EventStatus("EVENTS.EVENTS.STATUS.PROCESSED")
	ProcessingEventStatus          = EventStatus("EVENTS.EVENTS.STATUS.PROCESSING")
	ProcessingCancelledEventStatus = EventStatus("EVENTS.EVENTS.STATUS.PROCESSING_CANCELLED")
	ProcessingFailureEventStatus   = EventStatus("EVENTS.EVENTS.STATUS.PROCESSING_FAILURE")
	RecordingEventStatus           = EventStatus("EVENTS.EVENTS.STATUS.RECORDING")
	RecordingFailureEventStatus    = EventStatus("EVENTS.EVENTS.STATUS.RECORDING_FAILURE")
	ScheduledEventStatus           = EventStatus("EVENTS.EVENTS.STATUS.SCHEDULED")
)

type ProcessingState string

const (
	InstantiatedProcessingState = ProcessingState("INSTANTIATED")
	RunningProcessingState      = ProcessingState("RUNNING")
	StoppedProcessingState      = ProcessingState("STOPPED")
	PausedProcessingState       = ProcessingState("PAUSED")
	SucceededProcessingState    = ProcessingState("SUCCEEDED")
	FailedProcessingState       = ProcessingState("FAILED")
	FailingProcessingState      = ProcessingState("FAILING")
)

type Processing struct {
	Workflow      string     `json:"workflow,omitempty"`
	Configuration Properties `json:"configuration,omitempty"`
}

type Scheduling struct {
	Start   DateTime `json:"start,omitempty"`
	End     DateTime `json:"end,omitempty"`
	AgentID string   `json:"agent_id,omitempty"`
	Inputs  []string `json:"inputs,omitempty"`
}

type SchedulingRequest struct {
	AgentID  string    `json:"agent_id,omitempty"`
	Inputs   []string  `json:"inputs,omitempty"`
	Start    DateTime  `json:"start,omitempty"`
	End      *DateTime `json:"end,omitempty"`
	Duration *int64    `json:"duration,omitempty"`
	RRule    *RRule    `json:"rrule,omitempty"`
}

type RRule string

type Publication struct {
	ID          string              `json:"id,omitempty"`
	Channel     string              `json:"channel,omitempty"`
	MediaType   string              `json:"mediatype,omitempty"`
	URL         string              `json:"url,omitempty"`
	Media       []TrackElement      `json:"media,omitempty"`
	Attachments []AttachmentElement `json:"attachments,omitempty"`
	Metadata    []CatalogElement    `json:"metadata,omitempty"`
}

const (
	InternalChannel     = "internal"
	EngagePlayerChannel = "engage-player"
	EngageLiveChannel   = "engage-live"
)

type TrackElement struct {
	ID               string   `json:"id,omitempty"`
	MediaType        string   `json:"mediatype,omitempty"`
	URL              string   `json:"url,omitempty"`
	Flavor           Flavor   `json:"flavor,omitempty"`
	Size             int64    `json:"size,omitempty"`
	Checksum         string   `json:"checksum,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	HasAudio         bool     `json:"has_audio,omitempty"`
	HasVideo         bool     `json:"has_video,omitempty"`
	Duration         *int64   `json:"duration,omitempty"`
	Description      string   `json:"description,omitempty"`
	BitRate          *float64 `json:"bitrate,omitempty"`
	FrameRate        *float64 `json:"framerate,omitempty"`
	FrameCount       *int64   `json:"framecount,omitempty"`
	Width            *int     `json:"width,omitempty"`
	Height           *int     `json:"height,omitempty"`
	IsMasterPlaylist bool     `json:"is_master_playlist,omitempty"`
	IsLive           bool     `json:"is_live,omitempty"`
}

type MediaTrackElement struct {
	Checksum           *string                            `json:"checksum,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	Duration           *int64                             `json:"duration,omitempty"`
	ElementDescription *string                            `json:"element-description,omitempty"`
	Flavor             *Flavor                            `json:"flavor,omitempty"`
	Identifier         *string                            `json:"identifier,omitempty"`
	MimeType           *string                            `json:"mimetype,omitempty"`
	Size               int64                              `json:"size,omitempty"`
	HasVideo           bool                               `json:"has_video,omitempty"`
	HasAudio           bool                               `json:"has_audio,omitempty"`
	IsMasterPlaylist   bool                               `json:"is_master_playlist,omitempty"`
	IsLive             bool                               `json:"is_live,omitempty"`
	Streams            map[string]MediaTrackElementStream `json:"streams,omitempty"`
	Tags               []string                           `json:"tags,omitempty"`
	URI                *string                            `json:"uri,omitempty"`
}

type MediaTrackElementStream struct {
	// common fields
	Identifier           *string  `json:"identifier,omitempty"`
	BitRate              *float64 `json:"bitrate,omitempty"`
	CaptureDevice        *string  `json:"capturedevice,omitempty"`
	CaptureDeviceVendor  *string  `json:"capturedevicevendor,omitempty"`
	CaptureDeviceVersion *string  `json:"capturedeviceversion,omitempty"`
	EncoderLibraryVendor *string  `json:"encoderlibraryvendor,omitempty"`
	Format               *string  `json:"format,omitempty"`
	FormatVersion        *string  `json:"formatversion,omitempty"`
	FrameCount           *int64   `json:"framecount,omitempty"`

	// audio fields
	BitDepth     *int     `json:"bitdepth,omitempty"`
	Channels     *int     `json:"channels,omitempty"`
	PkLevDB      *float64 `json:"pklevdb,omitempty"`
	RmsLevDB     *float64 `json:"rmslevdb,omitempty"`
	RmsPkDB      *float64 `json:"rmspkdb,omitempty"`
	SamplingRate *int     `json:"samplingrate,omitempty"`

	// video fields
	FrameHeight *int       `json:"frameheight,omitempty"`
	FrameWidth  *int       `json:"framewidth,omitempty"`
	FrameRate   *float64   `json:"framerate,omitempty"`
	ScanOrder   *ScanOrder `json:"scanorder,omitempty"`
	ScanType    *ScanType  `json:"scantype,omitempty"`
}

type ScanOrder string

const (
	TopFieldFirstScanOrder    = ScanOrder("TopFieldFirst")
	BottomFieldFirstScanOrder = ScanOrder("BottomFieldFirst")
)

type ScanType string

const (
	InterlacedScanType  = ScanType("Interlaced")
	ProgressiveScanType = ScanType("Progressive")
)

type AttachmentElement struct {
	ID        string   `json:"id,omitempty"`
	MediaType string   `json:"mediatype,omitempty"`
	URL       string   `json:"url,omitempty"`
	Flavor    Flavor   `json:"flavor,omitempty"`
	Ref       string   `json:"ref,omitempty"`
	Size      int64    `json:"size,omitempty"`
	Checksum  string   `json:"checksum,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type CatalogElement struct {
	ID        string   `json:"id,omitempty"`
	MediaType string   `json:"mediatype,omitempty"`
	URL       string   `json:"url,omitempty"`
	Flavor    Flavor   `json:"flavor,omitempty"`
	Size      int64    `json:"size,omitempty"`
	Checksum  string   `json:"checksum,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type Series struct {
	Identifier   string   `json:"identifier,omitempty"`
	Title        string   `json:"title,omitempty"`
	Description  string   `json:"description,omitempty"`
	Creator      string   `json:"creator,omitempty"`
	Subjects     []string `json:"subjects,omitempty"`
	Organization string   `json:"organization,omitempty"`
	Created      DateTime `json:"created,omitempty"`
	Contributors []string `json:"contributors,omitempty"`
	Organizers   []string `json:"organizers,omitempty"`
	OptOut       bool     `json:"opt_out,omitempty"` // always false
	Publishers   []string `json:"publishers,omitempty"`
	Language     string   `json:"language,omitempty"`
	License      string   `json:"license,omitempty"`
	RightsHolder string   `json:"rightsholder,omitempty"`
	ACL          ACL      `json:"acl,omitempty"`
}

type Playlist struct {
	ID                   string          `json:"id,omitempty"`
	Entries              []PlaylistEntry `json:"entries,omitempty"`
	Title                string          `json:"title,omitempty"`
	Description          string          `json:"description,omitempty"`
	Creator              string          `json:"creator,omitempty"`
	Updated              string          `json:"updated,omitempty"`
	AccessControlEntries []PlaylistACE   `json:"accessControlEntries,omitempty"`
}

type PlaylistEntry struct {
	ID        string `json:"id,omitempty"`
	ContentID string `json:"contentId,omitempty"`
	Type      string `json:"type,omitempty"`
}

type PlaylistEntryType string

const (
	EventPlaylistEntryType        = PlaylistEntryType("EVENT")
	InaccessiblePlaylistEntryType = PlaylistEntryType("INACCESSIBLE")
)

type PlaylistACE struct {
	ID     string `json:"id,omitempty"`
	Allow  bool   `json:"allow,omitempty"`
	Action Action `json:"action,omitempty"`
	Role   string `json:"role,omitempty"`
}

type WorkflowInstance struct {
	Identifier                   int64               `json:"identifier,omitempty"`
	Title                        string              `json:"title,omitempty"`
	Description                  string              `json:"description,omitempty"`
	WorkflowDefinitionIdentifier string              `json:"workflow_definition_identifier,omitempty"`
	EventIdentifier              string              `json:"event_identifier,omitempty"`
	Creator                      string              `json:"creator,omitempty"`
	State                        WorkflowState       `json:"state,omitempty"`
	Operations                   []OperationInstance `json:"operations,omitempty"`
	Configuration                Properties          `json:"configuration,omitempty"`
}

type WorkflowState string

const (
	InstantiatedWorkflowState = WorkflowState("instantiated")
	RunningWorkflowState      = WorkflowState("running")
	StoppedWorkflowState      = WorkflowState("stopped")
	PausedWorkflowState       = WorkflowState("paused")
	SucceededWorkflowState    = WorkflowState("succeeded")
	FailedWorkflowState       = WorkflowState("failed")
	FailingWorkflowState      = WorkflowState("failing")
)

type OperationInstance struct {
	Identifier           int64                  `json:"identifier,omitempty"`
	Operation            string                 `json:"operation,omitempty"`
	Description          string                 `json:"description,omitempty"`
	State                WorkflowOperationState `json:"state,omitempty"`
	TimeInQueue          int64                  `json:"time_in_queue,omitempty"`
	Host                 string                 `json:"host,omitempty"`
	If                   string                 `json:"if,omitempty"`
	FailWorkflowOnError  bool                   `json:"fail_workflow_on_error,omitempty"`
	ErrorHandlerWorkflow string                 `json:"error_handler_workflow,omitempty"`
	RetryStrategy        WorkflowRetryStrategy  `json:"retry_strategy,omitempty"`
	MaxAttempts          int                    `json:"max_attempts,omitempty"`
	FailedAttempts       int                    `json:"failed_attempts,omitempty"`
	Configuration        Properties             `json:"configuration,omitempty"`
	Start                DateTime               `json:"start,omitempty"`
	Completion           DateTime               `json:"completion,omitempty"`
}

type WorkflowOperationState string

const (
	InstantiatedWorkflowOperationState = WorkflowOperationState("instantiated")
	RunningWorkflowOperationState      = WorkflowOperationState("running")
	PausedWorkflowOperationState       = WorkflowOperationState("paused")
	SucceededWorkflowOperationState    = WorkflowOperationState("succeeded")
	FailedWorkflowOperationState       = WorkflowOperationState("failed")
	SkippedWorkflowOperationState      = WorkflowOperationState("skipped")
	RetryWorkflowOperationState        = WorkflowOperationState("retry")
)

type WorkflowRetryStrategy string

const (
	NoneWorkflowRetryStrategy  = WorkflowRetryStrategy("none")
	RetryWorkflowRetryStrategy = WorkflowRetryStrategy("retry")
	HoldWorkflowRetryStrategy  = WorkflowRetryStrategy("hold")
)

type WorkflowDefinition struct {
	Identifier             string                `json:"identifier,omitempty"`
	Title                  string                `json:"title,omitempty"`
	Description            string                `json:"description,omitempty"`
	Tags                   []string              `json:"tags,omitempty"`
	ConfigurationPanel     *string               `json:"configuration_panel,omitempty"`
	ConfigurationPanelJSON *string               `json:"configuration_panel_json,omitempty"`
	Operations             []OperationDefinition `json:"operations,omitempty"`
}

type OperationDefinition struct {
	Operation            string                `json:"operation,omitempty"`
	Description          string                `json:"description,omitempty"`
	Configuration        Properties            `json:"configuration,omitempty"`
	If                   string                `json:"if,omitempty"`
	Unless               string                `json:"unless,omitempty"`
	FailWorkflowOnError  bool                  `json:"fail_workflow_on_error,omitempty"`
	ErrorHandlerWorkflow string                `json:"error_handler_workflow,omitempty"`
	RetryStrategy        WorkflowRetryStrategy `json:"retry_strategy,omitempty"`
	MaxAttempts          int                   `json:"max_attempts,omitempty"`
}
