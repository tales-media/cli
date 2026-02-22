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

	iso639_3 "github.com/barbashov/iso639-3"

	"shio.solutions/tales.media/cli/pkg/opencast/apis/meta/base"
)

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

type Property struct {
	Key   string `human:"Key" json:"key" yaml:"key"`
	Value string `human:"Value" json:"value" yaml:"value"`
}

type Flavor string

const (
	// Opencast
	DublinCoreEpisodeFlavor = Flavor("dublincore/episode")
	DublinCoreSeriesFlavor  = Flavor("dublincore/series")
	SecurityEpisodeFlavor   = Flavor("security/xacml+episode")
	SecuritySeriesFlavor    = Flavor("security/xacml+series")
	SMILCuttingFlavor       = Flavor("smil/cutting")
	MPEG7SegmentsFlavor     = Flavor("mpeg-7/segments")
)

type ACE struct {
	Role   string `human:"Role" json:"role" yaml:"role"`
	Action Action `human:"Action" json:"action" yaml:"action"`
	Allow  bool   `human:"Allow,wideonly" json:"allow" yaml:"allow"`
}

type Action string

const (
	// Opencast
	ReadAction  = Action("read")
	WriteAction = Action("write")

	// tales.media
	ReadMetadataAction  = Action("read-metadata")
	AnnotateAction      = Action("annotate")
	AnnotateAdminAction = Action("annotate-admin")
)

type Catalog struct {
	Flavor Flavor
	Fields map[string]Field
}

type Field struct {
	Type  FieldType
	Value FieldValue
}

type FieldType string

const (
	BooleanFieldType  = FieldType("boolean")
	NumberFieldType   = FieldType("number")
	TextFieldType     = FieldType("text")
	TextListFieldType = FieldType("text_list")
	DateTimeFieldType = FieldType("date_time")
)

type FieldValue any

type (
	BooleanFieldValue  bool
	NumberFieldValue   int64
	TextFieldValue     string
	TextListFieldValue []string
	DateTimeFieldValue = base.DateTime
)

var (
	_ FieldValue = BooleanFieldValue(true)
	_ FieldValue = NumberFieldValue(0)
	_ FieldValue = TextFieldValue("")
	_ FieldValue = TextListFieldValue{}
	_ FieldValue = DateTimeFieldValue{}
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

// TODO: SignedURL

type Group struct {
	ID             string   `human:"ID" json:"id" yaml:"id"`
	OrganizationID string   `human:"Organization ID" json:"organizationId" yaml:"organizationId"`
	Name           string   `human:"Name" json:"name" yaml:"name"`
	Description    string   `human:"Description,wideonly" json:"description" yaml:"description"`
	GroupRole      string   `human:"GroupRole" json:"groupRole" yaml:"groupRole"`
	Roles          []string `human:"Roles,wideonly" json:"roles" yaml:"roles"`
	Members        []string `human:"Members,wideonly" json:"members" yaml:"members"`
}

// TODO: StatisticProvider
// TODO: StatisticProviderType
// TODO: StatisticResourceType
// TODO: StatisticProviderParameter
// TODO: StatisticProviderParameterType
// TODO: StatisticQuery
// TODO: StatisticQueryResult
// TODO: StatisticQueryResultTimeSeriesData

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

type Event struct {
	// General

	ID              string         `human:"ID" json:"id" yaml:"id"`
	Status          EventStatus    `human:"Status" json:"eventStatus" yaml:"eventStatus"`
	WorkflowStatus  WorkflowStatus `human:"Workflow Status" json:"workflowStatus" yaml:"workflowStatus"`
	SnapshotVersion *int64         `human:"Snapshot Version,wideonly" json:"snapshotVersion" yaml:"snapshotVersion"`

	// Metadata

	Title        *string        `human:"Title" json:"title" yaml:"title"`
	Description  *string        `human:"Description,wideonly" json:"description" yaml:"description"`
	Series       *string        `human:"Series" json:"series" yaml:"series"`
	SeriesID     *string        `human:"Series ID,wideonly" json:"seriesId" yaml:"seriesId"`
	CreationDate *time.Time     `human:"Creation Date,wideonly" json:"creationDate" yaml:"creationDate"`
	StartDate    *time.Time     `human:"Start Date" json:"startDate" yaml:"startDate"`
	Duration     *time.Duration `human:"Duration" json:"duration" yaml:"duration"`
	EndDate      *time.Time     `human:"End Date" json:"endDate" yaml:"endDate"`
	Creator      *string        `human:"Creator" json:"creator" yaml:"creator"`
	Contributors []string       `human:"Contributors,wideonly" json:"contributors" yaml:"contributors"`
	Presenters   []string       `human:"Presenters,wideonly" json:"presenters" yaml:"presenters"`
	Location     *string        `human:"Location,wideonly" json:"location" yaml:"location"`
	Language     *Language      `human:"Language,wideonly" json:"language" yaml:"language"`
	RightsHolder *string        `human:"Rights Holder,wideonly" json:"rightsHolder" yaml:"rightsHolder"`
	License      *string        `human:"License,wideonly" json:"license" yaml:"license"`
	Subjects     []string       `human:"Subjects,wideonly" json:"subjects" yaml:"subjects"`
	Source       *string        `human:"Source,wideonly" json:"source" yaml:"source"`

	// Scheduling

	ScheduledStartDate *time.Time `human:"Scheduled Start Date,wideonly" json:"scheduledStartDate" yaml:"scheduledStartDate"`
	ScheduledEndDate   *time.Time `human:"Scheduled End Date,wideonly" json:"scheduledEndDate" yaml:"scheduledEndDate"`
	AgentName          *string    `human:"Agent Name,wideonly" json:"agentName" yaml:"agentName"`
	AgentInputs        []string   `human:"Agent Inputs,wideonly" json:"agentInputs" yaml:"agentInputs"`
}

type Language iso639_3.Language

var (
	_ encoding.TextMarshaler   = Language{}
	_ encoding.TextUnmarshaler = &Language{}
)

func (l Language) String() string {
	return l.Name
}

func (l Language) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Language) UnmarshalText(text []byte) error {
	return l.UnmarshalTextString(string(text))
}

func (l *Language) UnmarshalTextString(text string) error {
	optL := iso639_3.FromAnyCode(text)
	if optL != nil {
		*l = Language(*optL)
		return nil
	}

	optL = iso639_3.FromName(text)
	if optL != nil {
		*l = Language(*optL)
		return nil
	}

	return fmt.Errorf("Language: unknown langauge '%s'", text)
}

type EventStatus string

const (
	IngestingEventStatus           = EventStatus("ingesting")
	PausedEventStatus              = EventStatus("paused")
	PendingEventStatus             = EventStatus("pending")
	ProcessedEventStatus           = EventStatus("processed")
	ProcessingEventStatus          = EventStatus("processing")
	ProcessingCancelledEventStatus = EventStatus("processing_cancelled")
	ProcessingFailureEventStatus   = EventStatus("processing_failure")
	RecordingEventStatus           = EventStatus("recording")
	RecordingFailureEventStatus    = EventStatus("recording_failure")
	ScheduledEventStatus           = EventStatus("scheduled")
)

type Publication struct {
	ID          string              `human:"ID" json:"id" yaml:"id"`
	Channel     string              `human:"Channel" json:"channel" yaml:"channel"`
	MediaType   string              `human:"Media Type" json:"mediaType" yaml:"mediaType"`
	URL         string              `human:"URL" json:"url" yaml:"url"`
	Tracks      []TrackElement      `human:"Tracks,wideonly" json:"tracks" yaml:"tracks"`
	Catalogs    []CatalogElement    `human:"Catalogs,wideonly" json:"catalogs" yaml:"catalogs"`
	Attachments []AttachmentElement `human:"Attachments,wideonly" json:"attachments" yaml:"attachments"`
}

type TrackElement struct {
	// common
	ID        string   `human:"ID" json:"id" yaml:"id"`
	MediaType string   `human:"Media Type" json:"mediaType" yaml:"mediaType"`
	URL       string   `human:"URL" json:"url" yaml:"url"`
	Flavor    Flavor   `human:"Flavor" json:"flavor" yaml:"flavor"`
	Size      int64    `human:"Size,wideonly" json:"size" yaml:"size"`
	Checksum  string   `human:"Checksum,wideonly" json:"checksum" yaml:"checksum"`
	Tags      []string `human:"Tags,wideonly" json:"tags" yaml:"tags"`

	// track
	Duration          *time.Duration `human:"Duration,wideonly" json:"duration" yaml:"duration"`
	Width             *int64         `human:"Width,wideonly" json:"width" yaml:"width"`
	Height            *int64         `human:"Height,wideonly" json:"height" yaml:"height"`
	BitRate           *float64       `human:"BitRate,wideonly" json:"bitRate" yaml:"bitRate"`
	FrameRate         *float64       `human:"FrameRate,wideonly" json:"frameRate" yaml:"frameRate"`
	FrameCount        *int64         `human:"FrameCount,wideonly" json:"frameCount" yaml:"frameCount"`
	IsLive            bool           `human:"Is Live,wideonly" json:"isLiveStream" yaml:"isLiveStream"`
	HasAudioStreams   bool           `human:"Has Audio Streams,wideonly" json:"hasAudioStreams" yaml:"hasAudioStreams"`
	HasVideoStreams   bool           `human:"Has Video Streams,wideonly" json:"hasVideoStreams" yaml:"hasVideoStreams"`
	IsMainHLSPlaylist bool           `human:"Is Main HLS Playlist,wideonly" json:"isMainHlsPlaylist" yaml:"isMainHlsPlaylist"`
	Description       string         `human:"Description,wideonly" json:"description" yaml:"description"`

	// media track
	AudioStreams    []AudioStream    `human:"Audio Streams,wideonly" json:"audioStreams" yaml:"audioStreams"`
	VideoStreams    []VideoStream    `human:"Video Streams,wideonly" json:"videoStreams" yaml:"videoStreams"`
	SubtitleStreams []SubtitleStream `human:"Subtitle Streams,wideonly" json:"subtitleStreams" yaml:"subtitleStreams"`
}

type AudioStream struct {
	// common
	ID                   string   `human:"ID" json:"id" yaml:"id"`
	BitRate              *float64 `human:"Bit Rate" json:"bitRate" yaml:"bitRate"`
	CaptureDevice        *string  `human:"Capture Device" json:"captureDevice" yaml:"captureDevice"`
	CaptureDeviceVendor  *string  `human:"Capture Device Vendor" json:"captureDeviceVendor" yaml:"captureDeviceVendor"`
	CaptureDeviceVersion *string  `human:"Capture Device Version" json:"captureDeviceVersion" yaml:"captureDeviceVersion"`
	EncoderLibraryVendor *string  `human:"Encoder Library Vendor" json:"encoderLibraryVendor" yaml:"encoderLibraryVendor"`
	Format               *string  `human:"Format" json:"format" yaml:"format"`
	FormatVersion        *string  `human:"Format Version" json:"formatVersion" yaml:"formatVersion"`
	FrameCount           *int64   `human:"Frame Count" json:"frameCount" yaml:"frameCount"`

	// audio stream
	BitDepth            *int64   `human:"Bit Depth" json:"bitDepth" yaml:"bitDepth"`
	Channels            *int64   `human:"Channels" json:"channels" yaml:"channels"`
	SamplingRate        *int64   `human:"Sampling Rate" json:"samplingRate" yaml:"samplingRate"`
	PeakLevel           *float64 `human:"Peak Level" json:"peakLevel" yaml:"peakLevel"`
	RootMeanSquareLevel *float64 `human:"Root Mean Square Level" json:"rootMeanSquareLevel" yaml:"rootMeanSquareLevel"`
	RootMeanSquarePeak  *float64 `human:"Root Mean Square Peak" json:"rootMeanSquarePeak" yaml:"rootMeanSquarePeak"`
}

type VideoStream struct {
	// common
	ID                   string   `human:"ID" json:"id" yaml:"id"`
	BitRate              *float64 `human:"Bit Rate" json:"bitRate" yaml:"bitRate"`
	CaptureDevice        *string  `human:"Capture Device" json:"captureDevice" yaml:"captureDevice"`
	CaptureDeviceVendor  *string  `human:"Capture Device Vendor" json:"captureDeviceVendor" yaml:"captureDeviceVendor"`
	CaptureDeviceVersion *string  `human:"Capture Device Version" json:"captureDeviceVersion" yaml:"captureDeviceVersion"`
	EncoderLibraryVendor *string  `human:"Encoder Library Vendor" json:"encoderLibraryVendor" yaml:"encoderLibraryVendor"`
	Format               *string  `human:"Format" json:"format" yaml:"format"`
	FormatVersion        *string  `human:"Format Version" json:"formatVersion" yaml:"formatVersion"`
	FrameCount           *int64   `human:"Frame Count" json:"frameCount" yaml:"frameCount"`

	// video stream
	Height    *int64     `human:"Height" json:"height" yaml:"height"`
	Width     *int64     `human:"Width" json:"width" yaml:"width"`
	FrameRate *float64   `human:"Frame Rate" json:"frameRate" yaml:"frameRate"`
	ScanOrder *ScanOrder `human:"Scan Order" json:"scanOrder" yaml:"scanOrder"`
	ScanType  *ScanType  `human:"Scan Type" json:"scanType" yaml:"scanType"`
}

type SubtitleStream struct {
	// common
	ID                   string   `human:"ID" json:"id" yaml:"id"`
	BitRate              *float64 `human:"Bit Rate" json:"bitRate" yaml:"bitRate"`
	CaptureDevice        *string  `human:"Capture Device" json:"captureDevice" yaml:"captureDevice"`
	CaptureDeviceVendor  *string  `human:"Capture Device Vendor" json:"captureDeviceVendor" yaml:"captureDeviceVendor"`
	CaptureDeviceVersion *string  `human:"Capture Device Version" json:"captureDeviceVersion" yaml:"captureDeviceVersion"`
	EncoderLibraryVendor *string  `human:"Encoder Library Vendor" json:"encoderLibraryVendor" yaml:"encoderLibraryVendor"`
	Format               *string  `human:"Format" json:"format" yaml:"format"`
	FormatVersion        *string  `human:"Format Version" json:"formatVersion" yaml:"formatVersion"`
	FrameCount           *int64   `human:"Frame Count" json:"frameCount" yaml:"frameCount"`

	// subtitle stream
}

type ScanOrder string

const (
	TopFieldFirstScanOrder    = ScanOrder("top_field_first")
	BottomFieldFirstScanOrder = ScanOrder("bottom_field_first")
)

type ScanType string

const (
	InterlacedScanType  = ScanType("interlaced")
	ProgressiveScanType = ScanType("progressive")
)

type CatalogElement struct {
	// common
	ID        string   `human:"ID" json:"id" yaml:"id"`
	MediaType string   `human:"MediaType" json:"mediaType" yaml:"mediaType"`
	URL       string   `human:"URL" json:"url" yaml:"url"`
	Flavor    Flavor   `human:"Flavor" json:"flavor" yaml:"flavor"`
	Size      int64    `human:"Size" json:"size" yaml:"size"`
	Checksum  string   `human:"Checksum" json:"checksum" yaml:"checksum"`
	Tags      []string `human:"Tags,wideonly" json:"tags" yaml:"tags"`

	// catalog
}

type AttachmentElement struct {
	// common
	ID        string   `human:"ID" json:"id" yaml:"id"`
	MediaType string   `human:"MediaType" json:"mediaType" yaml:"mediaType"`
	URL       string   `human:"URL" json:"url" yaml:"url"`
	Flavor    Flavor   `human:"Flavor" json:"flavor" yaml:"flavor"`
	Size      int64    `human:"Size" json:"size" yaml:"size"`
	Checksum  string   `human:"Checksum" json:"checksum" yaml:"checksum"`
	Tags      []string `human:"Tags,wideonly" json:"tags" yaml:"tags"`

	// attachment
	Ref string `human:"Ref" json:"ref" yaml:"ref"`
}

type Series struct {
	// General

	ID string `human:"ID" json:"id" yaml:"id"`

	// Metadata

	Title        *string    `human:"Title" json:"title" yaml:"title"`
	Description  *string    `human:"Description,wideonly" json:"description" yaml:"description"`
	CreationDate *time.Time `human:"CreationDate" json:"creationDate" yaml:"creationDate"`
	Creator      *string    `human:"Creator" json:"creator" yaml:"creator"`
	Contributors []string   `human:"Contributors,wideonly" json:"contributors" yaml:"contributors"`
	Organizers   []string   `human:"Organizers,wideonly" json:"organizers" yaml:"organizers"`
	Publishers   []string   `human:"Publishers,wideonly" json:"publishers" yaml:"publishers"`
	Language     *Language  `human:"Language,wideonly" json:"language" yaml:"language"`
	RightsHolder *string    `human:"RightsHolder,wideonly" json:"rightsHolder" yaml:"rightsHolder"`
	License      *string    `human:"License,wideonly" json:"license" yaml:"license"`
	Subjects     []string   `human:"Subjects,wideonly" json:"subjects" yaml:"subjects"`
}

type Playlist struct {
	// General

	ID string `human:"ID" json:"id" yaml:"id"`

	// Metadata

	Title       *string    `human:"Title" json:"title" yaml:"title"`
	Description *string    `human:"Description,wideonly" json:"description" yaml:"description"`
	Creator     *string    `human:"Creator" json:"creator" yaml:"creator"`
	UpdateDate  *time.Time `human:"UpdateDate,wideonly" json:"updateDate" yaml:"updateDate"`
}

type PlaylistEntry struct {
	ID        int64             `human:"ID" json:"id" yaml:"id"`
	Type      PlaylistEntryType `human:"Type" json:"type" yaml:"type"`
	ContentID string            `human:"ContentID" json:"contentId" yaml:"contentId"`
}

type PlaylistEntryType string

const (
	EventPlaylistEntryType        = PlaylistEntryType("event")
	InaccessiblePlaylistEntryType = PlaylistEntryType("inaccessible")
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
