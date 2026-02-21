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
	"fmt"
	"strings"
	"time"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	extapiv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
	"shio.solutions/tales.media/cli/pkg/opencast/apis/meta/base"
)

func Map[F, T any](ocList []F, mapper func(F) T) []T {
	list := make([]T, 0, len(ocList))
	for _, item := range ocList {
		list = append(list, mapper(item))
	}
	return list
}

func deref[T any](ptr *T, d ...T) T {
	if ptr != nil {
		return *ptr
	}
	if len(d) > 0 {
		return d[0]
	}
	var zero T
	return zero
}

func SortDirectionToOCSortDirection(sortDirection api.SortDirection) extapiclientv1.SortDirection {
	switch sortDirection {
	case api.Ascending:
		return extapiclientv1.Ascending
	case api.Descending:
		return extapiclientv1.Descending
	}
	panic(fmt.Sprintf("BUG: unknown SortDirection '%d'", sortDirection))
}

func OCPropertiesToMap(ocProperties base.Properties) map[string]string {
	return map[string]string(ocProperties)
}

func OCPropertiesToProperties(ocProperties base.Properties) []api.Property {
	properties := make([]api.Property, 0, len(ocProperties))
	for k, v := range ocProperties {
		properties = append(properties, api.Property{
			Key:   k,
			Value: v,
		})
	}
	return properties
}

func FlavorToOCFlavor(flavor api.Flavor) base.Flavor {
	// We are not strictly enforcing the Flavor name as this can be configured in Opencast.
	return base.Flavor(flavor)
}

func OCFlavorToFlavor(ocFlavor base.Flavor) api.Flavor {
	// We are not strictly enforcing the Flavor name as this can be configured in Opencast.
	return api.Flavor(ocFlavor)
}

func OCACEToACE(ocACE extapiv1.ACE) api.ACE {
	return api.ACE{
		Role:   ocACE.Role,
		Action: OCActionToAction(ocACE.Action),
	}
}

func OCActionToAction(ocAction base.Action) api.Action {
	// We are not strictly enforcing the action name as this can be configured in Opencast.
	return api.Action(ocAction)
}

func OCCatalogToCatalog(ocCatalog extapiv1.Catalog) api.Catalog {
	c := api.Catalog{
		Flavor: OCFlavorToFlavor(ocCatalog.Flavor),
		Fields: make(map[string]api.Field, len(ocCatalog.Fields)),
	}
	for _, ocField := range ocCatalog.Fields {
		switch ocField.Type {
		case extapiv1.BooleanFieldType:
			c.Fields[ocField.ID] = OCBooleanFieldToField(ocField)
		case extapiv1.DateFieldType:
			c.Fields[ocField.ID] = OCDateFieldToField(ocField)
		case extapiv1.MixedTextFieldType:
			c.Fields[ocField.ID] = OCMixedTextFieldToField(ocField)
		case extapiv1.IterableTextFieldType:
			c.Fields[ocField.ID] = OCIterableTextFieldToField(ocField)
		case extapiv1.NumberFieldType:
			c.Fields[ocField.ID] = OCNumberFieldToField(ocField)
		case extapiv1.OrderedTextFieldType:
			c.Fields[ocField.ID] = OCOrderedTextFieldToField(ocField)
		case extapiv1.TextFieldType:
			c.Fields[ocField.ID] = OCTextFieldToField(ocField)
		case extapiv1.TextLongFieldType:
			c.Fields[ocField.ID] = OCTextLongFieldToField(ocField)
		case extapiv1.TimeFieldType:
			c.Fields[ocField.ID] = OCTimeFieldToField(ocField)
		default:
			panic(fmt.Sprintf("BUG: unknown Opencast catalog field type '%s'", ocField.Type))
		}
	}
	return c
}

func OCBooleanFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.BooleanFieldValue)
	value := api.BooleanFieldValue(ocValue)
	return api.Field{
		Type:  api.BooleanFieldType,
		Value: value,
	}
}

func OCDateFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.DateTimeFieldValue)
	value := api.DateTimeFieldValue(ocValue)
	return api.Field{
		Type:  api.DateTimeFieldType,
		Value: value,
	}
}

func OCMixedTextFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.MixedTextFieldValue)
	value := api.TextListFieldValue(ocValue)
	return api.Field{
		Type:  api.TextListFieldType,
		Value: value,
	}
}

func OCIterableTextFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.IterableTextFieldValue)
	value := api.TextListFieldValue(ocValue)
	return api.Field{
		Type:  api.TextListFieldType,
		Value: value,
	}
}

func OCNumberFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.NumberFieldValue)
	value := api.NumberFieldValue(ocValue)
	return api.Field{
		Type:  api.NumberFieldType,
		Value: value,
	}
}

func OCOrderedTextFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.OrderedTextFieldValue)
	value := api.TextFieldValue(ocValue)
	return api.Field{
		Type:  api.TextFieldType,
		Value: value,
	}
}

func OCTextFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.TextFieldValue)
	value := api.TextFieldValue(ocValue)
	return api.Field{
		Type:  api.TextFieldType,
		Value: value,
	}
}

func OCTextLongFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.TextLongFieldValue)
	value := api.TextFieldValue(ocValue)
	return api.Field{
		Type:  api.TextFieldType,
		Value: value,
	}
}

func OCTimeFieldToField(ocField extapiv1.Field) api.Field {
	ocValue := ocField.Value.(extapiv1.TimeFieldValue)
	value := api.DateTimeFieldValue(ocValue)
	return api.Field{
		Type:  api.DateTimeFieldType,
		Value: value,
	}
}

func OCAPIToAPIInfo(ocAPI extapiv1.API, ocVersion extapiv1.APIVersion) api.APIInfo {
	return api.APIInfo{
		URL:               ocAPI.URL,
		DefaultVersion:    ocAPI.Version,
		SupportedVersions: ocVersion.Versions,
	}
}

func OCMeToMe(ocMe extapiv1.Me, roles []string) api.Me {
	return api.Me{
		Username: ocMe.Username,
		Name:     ocMe.Name,
		Email:    ocMe.Email,
		UserRole: ocMe.UserRole,
		Provider: ocMe.Provider,
		Roles:    roles,
	}
}

func OCOrganizationToOrganization(ocOrganization extapiv1.Organization, ocProperties base.Properties) api.Organization {
	return api.Organization{
		ID:            ocOrganization.ID,
		Name:          ocOrganization.Name,
		AdminRole:     ocOrganization.AdminRole,
		AnonymousRole: ocOrganization.AnonymousRole,
		Properties:    ocProperties,
	}
}

// TODO: SignedURL

func OCGroupToGroup(ocGroup extapiv1.Group) api.Group {
	return api.Group{
		ID:             ocGroup.Identifier,
		OrganizationID: ocGroup.Organization,
		Name:           ocGroup.Name,
		Description:    ocGroup.Description,
		GroupRole:      ocGroup.Role,
		Roles:          strings.Split(ocGroup.Roles, ","),
		Members:        strings.Split(ocGroup.Members, ","),
	}
}

// TODO: StatisticProvider
// TODO: StatisticProviderType
// TODO: StatisticResourceType
// TODO: StatisticProviderParameter
// TODO: StatisticProviderParameterType
// TODO: StatisticQuery
// TODO: StatisticQueryResult
// TODO: StatisticQueryResultTimeSeriesData

func OCAgentToAgent(ocAgent extapiv1.Agent) api.Agent {
	a := api.Agent{
		Name:   ocAgent.AgentID,
		URL:    ocAgent.URL,
		Status: OCAgentStatusToAgentStatus(ocAgent.Status),
		Inputs: ocAgent.Inputs,
	}
	if !ocAgent.Update.IsZero() {
		a.LastUpdate = new(ocAgent.Update.Time)
	}
	return a
}

func OCAgentStatusToAgentStatus(ocAgentStatus extapiv1.AgentStatus) api.AgentStatus {
	if ocAgentStatus == "" {
		return ""
	}
	switch ocAgentStatus {
	case extapiv1.UndefinedAgentStatus:
		return ""
	case extapiv1.UnknownAgentStatus:
		return api.UnknownAgentStatus
	case extapiv1.IdleAgentStatus:
		return api.IdleAgentStatus
	case extapiv1.CapturingAgentStatus:
		return api.CapturingAgentStatus
	case extapiv1.UploadingAgentStatus:
		return api.UploadingAgentStatus
	case extapiv1.ShuttingDownAgentStatus:
		return api.ShuttingDownAgentStatus
	case extapiv1.OfflineAgentStatus:
		return api.OfflineAgentStatus
	case extapiv1.ErrorAgentStatus:
		return api.ErrorAgentStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast AgentStatus '%s'", ocAgentStatus))
}

func OCEventToEvent(ocEvent extapiv1.Event) api.Event {
	e := api.Event{
		ID:              ocEvent.Identifier,
		Status:          OCEventStatusToEventStatus(ocEvent.Status),
		WorkflowStatus:  OCProcessingStateToWorkflowStatus(ocEvent.ProcessingState),
		SnapshotVersion: (*int64)(ocEvent.ArchiveVersion),
	}
	if ocEvent.Title != "" {
		e.Title = new(ocEvent.Title)
	}
	if ocEvent.Description != "" {
		e.Description = new(ocEvent.Description)
	}
	if ocEvent.Description != "" {
		e.Description = new(ocEvent.Description)
	}
	if ocEvent.Series != "" {
		e.Series = new(ocEvent.Series)
	}
	if ocEvent.IsPartOf != "" {
		e.SeriesID = new(ocEvent.IsPartOf)
	}
	if !ocEvent.Created.IsZero() {
		e.CreationDate = new(ocEvent.Created.Time)
	}
	if !ocEvent.Start.IsZero() {
		e.StartDate = new(ocEvent.Start.Time)
	}
	if ocEvent.Duration != nil {
		e.Duration = new(time.Duration(*ocEvent.Duration) * time.Millisecond)
	}
	if e.StartDate != nil && e.Duration != nil {
		e.EndDate = new(e.StartDate.Add(*e.Duration))
	}
	if ocEvent.Creator != "" {
		e.Creator = new(ocEvent.Creator)
	}
	if len(ocEvent.Contributor) > 0 {
		e.Contributors = ocEvent.Contributor
	}
	if len(ocEvent.Presenter) > 0 {
		e.Presenters = ocEvent.Presenter
	}
	if ocEvent.Location != "" {
		e.Location = new(ocEvent.Location)
	}
	if ocEvent.Language != "" {
		var l api.Language
		err := l.UnmarshalTextString(ocEvent.Language)
		// TODO: should we handle the error?
		if err == nil {
			e.Language = &l
		}
	}
	if ocEvent.RightsHolder != "" {
		e.RightsHolder = new(ocEvent.RightsHolder)
	}
	if ocEvent.License != "" {
		e.License = new(ocEvent.License)
	}
	if len(ocEvent.Subjects) > 0 {
		e.Subjects = ocEvent.Subjects
	}
	if ocEvent.Source != "" {
		e.Source = new(ocEvent.Source)
	}
	if !ocEvent.Scheduling.Start.IsZero() {
		e.ScheduledStartDate = new(ocEvent.Scheduling.Start.Time)
	}
	if !ocEvent.Scheduling.End.IsZero() {
		e.ScheduledEndDate = new(ocEvent.Scheduling.End.Time)
	}
	if ocEvent.Scheduling.AgentID != "" {
		e.AgentName = new(ocEvent.Scheduling.AgentID)
	}
	if len(ocEvent.Scheduling.Inputs) > 0 {
		e.AgentInputs = ocEvent.Scheduling.Inputs
	}
	return e
}

func OCEventStatusToEventStatus(ocEventStatus extapiv1.EventStatus) api.EventStatus {
	if ocEventStatus == "" {
		return ""
	}
	switch ocEventStatus {
	case extapiv1.IngestingEventStatus:
		return api.IngestingEventStatus
	case extapiv1.PausedEventStatus:
		return api.PausedEventStatus
	case extapiv1.PendingEventStatus:
		return api.PendingEventStatus
	case extapiv1.ProcessedEventStatus:
		return api.ProcessedEventStatus
	case extapiv1.ProcessingEventStatus:
		return api.ProcessingEventStatus
	case extapiv1.ProcessingCancelledEventStatus:
		return api.ProcessingCancelledEventStatus
	case extapiv1.ProcessingFailureEventStatus:
		return api.ProcessingFailureEventStatus
	case extapiv1.RecordingEventStatus:
		return api.RecordingEventStatus
	case extapiv1.RecordingFailureEventStatus:
		return api.RecordingFailureEventStatus
	case extapiv1.ScheduledEventStatus:
		return api.ScheduledEventStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast EventStatus '%s'", ocEventStatus))
}

func EventStatusToOCEventStatus(eventStatus api.EventStatus) extapiv1.EventStatus {
	if eventStatus == "" {
		return ""
	}
	switch eventStatus {
	case api.IngestingEventStatus:
		return extapiv1.IngestingEventStatus
	case api.PausedEventStatus:
		return extapiv1.PausedEventStatus
	case api.PendingEventStatus:
		return extapiv1.PendingEventStatus
	case api.ProcessedEventStatus:
		return extapiv1.ProcessedEventStatus
	case api.ProcessingEventStatus:
		return extapiv1.ProcessingEventStatus
	case api.ProcessingCancelledEventStatus:
		return extapiv1.ProcessingCancelledEventStatus
	case api.ProcessingFailureEventStatus:
		return extapiv1.ProcessingFailureEventStatus
	case api.RecordingEventStatus:
		return extapiv1.RecordingEventStatus
	case api.RecordingFailureEventStatus:
		return extapiv1.RecordingFailureEventStatus
	case api.ScheduledEventStatus:
		return extapiv1.ScheduledEventStatus
	}
	panic(fmt.Sprintf("BUG: unknown EventStatus '%s'", eventStatus))
}

func OCProcessingStateToWorkflowStatus(ocProcessingState extapiv1.ProcessingState) api.WorkflowStatus {
	if ocProcessingState == "" {
		return ""
	}
	switch ocProcessingState {
	case extapiv1.UndefinedProcessingState:
		return ""
	case extapiv1.InstantiatedProcessingState:
		return api.InstantiatedWorkflowStatus
	case extapiv1.RunningProcessingState:
		return api.RunningWorkflowStatus
	case extapiv1.StoppedProcessingState:
		return api.StoppedWorkflowStatus
	case extapiv1.PausedProcessingState:
		return api.PausedWorkflowStatus
	case extapiv1.SucceededProcessingState:
		return api.SucceededWorkflowStatus
	case extapiv1.FailedProcessingState:
		return api.FailedWorkflowStatus
	case extapiv1.FailingProcessingState:
		return api.FailingWorkflowStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast ProcessingState '%s'", ocProcessingState))
}

func OCPublicationToPublication(ocPublication extapiv1.Publication) api.Publication {
	return api.Publication{
		ID:          ocPublication.ID,
		Channel:     ocPublication.Channel,
		MediaType:   ocPublication.MediaType,
		URL:         ocPublication.URL,
		Tracks:      Map(ocPublication.Media, OCTrackElementToTrackElement),
		Catalogs:    Map(ocPublication.Metadata, OCCatalogElementToCatalogElement),
		Attachments: Map(ocPublication.Attachments, OCAttachmentElementToAttachmentElement),
	}
}

func OCTrackElementToTrackElement(ocTrackElement extapiv1.TrackElement) api.TrackElement {
	t := api.TrackElement{
		ID:                ocTrackElement.ID,
		MediaType:         ocTrackElement.MediaType,
		URL:               ocTrackElement.URL,
		Flavor:            OCFlavorToFlavor(ocTrackElement.Flavor),
		Size:              int64(ocTrackElement.Size),
		Checksum:          ocTrackElement.Checksum,
		Tags:              ocTrackElement.Tags,
		Width:             (*int64)(ocTrackElement.Width),
		Height:            (*int64)(ocTrackElement.Height),
		BitRate:           (*float64)(ocTrackElement.BitRate),
		FrameRate:         (*float64)(ocTrackElement.FrameRate),
		FrameCount:        (*int64)(ocTrackElement.FrameCount),
		IsLive:            ocTrackElement.IsLive,
		HasAudioStreams:   ocTrackElement.HasAudio,
		HasVideoStreams:   ocTrackElement.HasVideo,
		IsMainHLSPlaylist: ocTrackElement.IsMasterPlaylist,
		Description:       ocTrackElement.Description,
	}
	if ocTrackElement.Duration != nil {
		t.Duration = new(time.Duration(*ocTrackElement.Duration) * time.Millisecond)
	}
	return t
}

func OCMediaTrackElementToTrackElement(ocMediaTrackElement extapiv1.MediaTrackElement) api.TrackElement {
	t := api.TrackElement{
		ID:                deref(ocMediaTrackElement.Identifier),
		MediaType:         deref(ocMediaTrackElement.MimeType),
		URL:               deref(ocMediaTrackElement.URI),
		Flavor:            OCFlavorToFlavor(deref(ocMediaTrackElement.Flavor)),
		Size:              int64(ocMediaTrackElement.Size),
		Checksum:          deref(ocMediaTrackElement.Checksum),
		Tags:              ocMediaTrackElement.Tags,
		IsLive:            ocMediaTrackElement.IsLive,
		HasAudioStreams:   ocMediaTrackElement.HasAudio,
		HasVideoStreams:   ocMediaTrackElement.HasVideo,
		IsMainHLSPlaylist: ocMediaTrackElement.IsMasterPlaylist,
		Description:       deref(ocMediaTrackElement.Description),
	}
	if ocMediaTrackElement.Duration != nil {
		t.Duration = new(time.Duration(*ocMediaTrackElement.Duration) * time.Millisecond)
	}
	for streamID, stream := range ocMediaTrackElement.Streams {
		switch {
		case strings.HasPrefix(streamID, "audio"):
			t.AudioStreams = append(t.AudioStreams, OCMediaTrackElementStreamToAudioStream(stream))
		case strings.HasPrefix(streamID, "video"):
			t.VideoStreams = append(t.VideoStreams, OCMediaTrackElementStreamToVideoStream(stream))
		case strings.HasPrefix(streamID, "subtitle"):
			t.SubtitleStreams = append(t.SubtitleStreams, OCMediaTrackElementStreamToSubtitleStream(stream))
		}
	}
	return t
}

func OCMediaTrackElementStreamToAudioStream(ocMediaTrackElementStream extapiv1.MediaTrackElementStream) api.AudioStream {
	return api.AudioStream{
		ID:                   deref(ocMediaTrackElementStream.Identifier),
		BitRate:              (*float64)(ocMediaTrackElementStream.BitRate),
		CaptureDevice:        ocMediaTrackElementStream.CaptureDevice,
		CaptureDeviceVendor:  ocMediaTrackElementStream.CaptureDeviceVendor,
		CaptureDeviceVersion: ocMediaTrackElementStream.CaptureDeviceVersion,
		EncoderLibraryVendor: ocMediaTrackElementStream.EncoderLibraryVendor,
		Format:               ocMediaTrackElementStream.Format,
		FormatVersion:        ocMediaTrackElementStream.FormatVersion,
		FrameCount:           (*int64)(ocMediaTrackElementStream.FrameCount),
		BitDepth:             (*int64)(ocMediaTrackElementStream.BitDepth),
		Channels:             (*int64)(ocMediaTrackElementStream.Channels),
		SamplingRate:         (*int64)(ocMediaTrackElementStream.SamplingRate),
		PeakLevel:            (*float64)(ocMediaTrackElementStream.PkLevDB),
		RootMeanSquareLevel:  (*float64)(ocMediaTrackElementStream.RmsLevDB),
		RootMeanSquarePeak:   (*float64)(ocMediaTrackElementStream.RmsPkDB),
	}
}

func OCMediaTrackElementStreamToVideoStream(ocMediaTrackElementStream extapiv1.MediaTrackElementStream) api.VideoStream {
	return api.VideoStream{
		ID:                   deref(ocMediaTrackElementStream.Identifier),
		BitRate:              (*float64)(ocMediaTrackElementStream.BitRate),
		CaptureDevice:        ocMediaTrackElementStream.CaptureDevice,
		CaptureDeviceVendor:  ocMediaTrackElementStream.CaptureDeviceVendor,
		CaptureDeviceVersion: ocMediaTrackElementStream.CaptureDeviceVersion,
		EncoderLibraryVendor: ocMediaTrackElementStream.EncoderLibraryVendor,
		Format:               ocMediaTrackElementStream.Format,
		FormatVersion:        ocMediaTrackElementStream.FormatVersion,
		FrameCount:           (*int64)(ocMediaTrackElementStream.FrameCount),
		Height:               (*int64)(ocMediaTrackElementStream.FrameHeight),
		Width:                (*int64)(ocMediaTrackElementStream.FrameWidth),
		FrameRate:            (*float64)(ocMediaTrackElementStream.FrameRate),
		ScanOrder:            new(OCScanOrderToScanOrder(deref(ocMediaTrackElementStream.ScanOrder))),
		ScanType:             new(OCScanTypeToScanType(deref(ocMediaTrackElementStream.ScanType))),
	}
}

func OCMediaTrackElementStreamToSubtitleStream(ocMediaTrackElementStream extapiv1.MediaTrackElementStream) api.SubtitleStream {
	return api.SubtitleStream{
		ID:                   deref(ocMediaTrackElementStream.Identifier),
		BitRate:              (*float64)(ocMediaTrackElementStream.BitRate),
		CaptureDevice:        ocMediaTrackElementStream.CaptureDevice,
		CaptureDeviceVendor:  ocMediaTrackElementStream.CaptureDeviceVendor,
		CaptureDeviceVersion: ocMediaTrackElementStream.CaptureDeviceVersion,
		EncoderLibraryVendor: ocMediaTrackElementStream.EncoderLibraryVendor,
		Format:               ocMediaTrackElementStream.Format,
		FormatVersion:        ocMediaTrackElementStream.FormatVersion,
		FrameCount:           (*int64)(ocMediaTrackElementStream.FrameCount),
	}
}

func OCScanOrderToScanOrder(ocScanOrder extapiv1.ScanOrder) api.ScanOrder {
	if ocScanOrder == "" {
		return ""
	}
	switch ocScanOrder {
	case extapiv1.UndefinedScanOrder:
		return ""
	case extapiv1.TopFieldFirstScanOrder:
		return api.TopFieldFirstScanOrder
	case extapiv1.BottomFieldFirstScanOrder:
		return api.BottomFieldFirstScanOrder
	}
	panic(fmt.Sprintf("BUG: unknown Opencast ScanOrder '%s'", ocScanOrder))
}

func OCScanTypeToScanType(ocScanType extapiv1.ScanType) api.ScanType {
	if ocScanType == "" {
		return ""
	}
	switch ocScanType {
	case extapiv1.UndefinedScanType:
		return ""
	case extapiv1.InterlacedScanType:
		return api.InterlacedScanType
	case extapiv1.ProgressiveScanType:
		return api.ProgressiveScanType
	}
	panic(fmt.Sprintf("BUG: unknown Opencast ScanType '%s'", ocScanType))
}

func OCCatalogElementToCatalogElement(ocCatalogElement extapiv1.CatalogElement) api.CatalogElement {
	return api.CatalogElement{
		ID:        ocCatalogElement.ID,
		MediaType: ocCatalogElement.MediaType,
		URL:       ocCatalogElement.URL,
		Flavor:    OCFlavorToFlavor(ocCatalogElement.Flavor),
		Size:      int64(ocCatalogElement.Size),
		Checksum:  ocCatalogElement.Checksum,
		Tags:      ocCatalogElement.Tags,
	}
}

func OCAttachmentElementToAttachmentElement(ocAttachmentElement extapiv1.AttachmentElement) api.AttachmentElement {
	return api.AttachmentElement{
		ID:        ocAttachmentElement.ID,
		MediaType: ocAttachmentElement.MediaType,
		URL:       ocAttachmentElement.URL,
		Flavor:    OCFlavorToFlavor(ocAttachmentElement.Flavor),
		Size:      int64(ocAttachmentElement.Size),
		Checksum:  ocAttachmentElement.Checksum,
		Tags:      ocAttachmentElement.Tags,
		Ref:       ocAttachmentElement.Ref,
	}
}

func OCSeriesToSeries(ocSeries extapiv1.Series) api.Series {
	s := api.Series{
		ID: ocSeries.Identifier,
	}
	if ocSeries.Title != "" {
		s.Title = new(ocSeries.Title)
	}
	if ocSeries.Description != "" {
		s.Description = new(ocSeries.Description)
	}
	if !ocSeries.Created.IsZero() {
		s.CreationDate = new(ocSeries.Created.Time)
	}
	if ocSeries.Creator != "" {
		s.Creator = new(ocSeries.Creator)
	}
	if len(ocSeries.Contributors) > 0 {
		s.Contributors = ocSeries.Contributors
	}
	if len(ocSeries.Organizers) > 0 {
		s.Organizers = ocSeries.Organizers
	}
	if len(ocSeries.Publishers) > 0 {
		s.Publishers = ocSeries.Publishers
	}
	if ocSeries.Language != "" {
		var l api.Language
		err := l.UnmarshalTextString(ocSeries.Language)
		// TODO: should we handle the error?
		if err == nil {
			s.Language = &l
		}
	}
	if ocSeries.RightsHolder != "" {
		s.RightsHolder = new(ocSeries.RightsHolder)
	}
	if ocSeries.License != "" {
		s.License = new(ocSeries.License)
	}
	if len(ocSeries.Subjects) > 0 {
		s.Subjects = ocSeries.Subjects
	}
	return s
}

func OCWorkflowInstanceToWorkflow(ocWorkflowInstance extapiv1.WorkflowInstance) api.Workflow {
	return api.Workflow{
		ID:                   int64(ocWorkflowInstance.Identifier),
		Title:                ocWorkflowInstance.Title,
		Description:          ocWorkflowInstance.Description,
		WorkflowDefinitionID: ocWorkflowInstance.WorkflowDefinitionIdentifier,
		EventID:              ocWorkflowInstance.EventIdentifier,
		Creator:              ocWorkflowInstance.Creator,
		Status:               OCWorkflowStateToWorkflowStatus(ocWorkflowInstance.State),
		Operations:           Map(ocWorkflowInstance.Operations, OCOperationInstanceToOperation),
		Configuration:        ocWorkflowInstance.Configuration,
	}
}

func OCWorkflowStateToWorkflowStatus(ocWorkflowState extapiv1.WorkflowState) api.WorkflowStatus {
	if ocWorkflowState == "" {
		return ""
	}
	switch ocWorkflowState {
	case extapiv1.InstantiatedWorkflowState:
		return api.InstantiatedWorkflowStatus
	case extapiv1.RunningWorkflowState:
		return api.RunningWorkflowStatus
	case extapiv1.StoppedWorkflowState:
		return api.StoppedWorkflowStatus
	case extapiv1.PausedWorkflowState:
		return api.PausedWorkflowStatus
	case extapiv1.SucceededWorkflowState:
		return api.SucceededWorkflowStatus
	case extapiv1.FailedWorkflowState:
		return api.FailedWorkflowStatus
	case extapiv1.FailingWorkflowState:
		return api.FailingWorkflowStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowState '%s'", ocWorkflowState))
}

func WorkflowStatusToOCWorkflowState(workflowState api.WorkflowStatus) extapiv1.WorkflowState {
	if workflowState == "" {
		return ""
	}
	switch workflowState {
	case api.InstantiatedWorkflowStatus:
		return extapiv1.InstantiatedWorkflowState
	case api.RunningWorkflowStatus:
		return extapiv1.RunningWorkflowState
	case api.StoppedWorkflowStatus:
		return extapiv1.StoppedWorkflowState
	case api.PausedWorkflowStatus:
		return extapiv1.PausedWorkflowState
	case api.SucceededWorkflowStatus:
		return extapiv1.SucceededWorkflowState
	case api.FailedWorkflowStatus:
		return extapiv1.FailedWorkflowState
	case api.FailingWorkflowStatus:
		return extapiv1.FailingWorkflowState
	}
	panic(fmt.Sprintf("BUG: unknown WorkflowState '%s'", workflowState))
}

func OCOperationInstanceToOperation(ocOperationInstance extapiv1.OperationInstance) api.Operation {
	o := api.Operation{
		ID:                   (*int64)(ocOperationInstance.Identifier),
		Operation:            ocOperationInstance.Operation,
		Description:          ocOperationInstance.Description,
		Status:               OCWorkflowOperationStateToWorkflowOperationStatus(ocOperationInstance.State),
		TimeInQueue:          int64(ocOperationInstance.TimeInQueue),
		Host:                 ocOperationInstance.Host,
		If:                   ocOperationInstance.If,
		FailWorkflowOnError:  ocOperationInstance.FailWorkflowOnError,
		ErrorHandlerWorkflow: ocOperationInstance.ErrorHandlerWorkflow,
		RetryStrategy:        OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocOperationInstance.RetryStrategy),
		MaxAttempts:          int64(ocOperationInstance.MaxAttempts),
		FailedAttempts:       int64(ocOperationInstance.FailedAttempts),
		Configuration:        ocOperationInstance.Configuration,
	}
	if !ocOperationInstance.Start.IsZero() {
		o.Start = new(ocOperationInstance.Start.Time)
	}
	if !ocOperationInstance.Completion.IsZero() {
		o.Completion = new(ocOperationInstance.Completion.Time)
	}
	return o
}

func OCWorkflowOperationStateToWorkflowOperationStatus(ocWorkflowOperationState extapiv1.WorkflowOperationState) api.WorkflowOperationStatus {
	if ocWorkflowOperationState == "" {
		return ""
	}
	switch ocWorkflowOperationState {
	case extapiv1.InstantiatedWorkflowOperationState:
		return api.InstantiatedWorkflowOperationStatus
	case extapiv1.RunningWorkflowOperationState:
		return api.RunningWorkflowOperationStatus
	case extapiv1.PausedWorkflowOperationState:
		return api.PausedWorkflowOperationStatus
	case extapiv1.SucceededWorkflowOperationState:
		return api.SucceededWorkflowOperationStatus
	case extapiv1.FailedWorkflowOperationState:
		return api.FailedWorkflowOperationStatus
	case extapiv1.SkippedWorkflowOperationState:
		return api.SkippedWorkflowOperationStatus
	case extapiv1.RetryWorkflowOperationState:
		return api.RetryWorkflowOperationStatus
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowOperationState '%s'", ocWorkflowOperationState))
}

func OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocWorkflowRetryStrategy extapiv1.WorkflowRetryStrategy) api.WorkflowRetryStrategy {
	if ocWorkflowRetryStrategy == "" {
		return ""
	}
	switch ocWorkflowRetryStrategy {
	case extapiv1.UndefinedWorkflowRetryStrategy:
		return ""
	case extapiv1.NoneWorkflowRetryStrategy:
		return api.NoneWorkflowRetryStrategy
	case extapiv1.RetryWorkflowRetryStrategy:
		return api.RetryWorkflowRetryStrategy
	case extapiv1.HoldWorkflowRetryStrategy:
		return api.HoldWorkflowRetryStrategy
	}
	panic(fmt.Sprintf("BUG: unknown Opencast WorkflowRetryStrategy '%s'", ocWorkflowRetryStrategy))
}

func OCWorkflowDefinitionToWorkflowDefinition(ocWorkflowDefinition extapiv1.WorkflowDefinition) api.WorkflowDefinition {
	return api.WorkflowDefinition{
		ID:                     ocWorkflowDefinition.Identifier,
		Title:                  ocWorkflowDefinition.Title,
		Description:            ocWorkflowDefinition.Description,
		Tags:                   ocWorkflowDefinition.Tags,
		ConfigurationPanel:     ocWorkflowDefinition.ConfigurationPanel,
		ConfigurationPanelJSON: ocWorkflowDefinition.ConfigurationPanelJSON,
		Operations:             Map(ocWorkflowDefinition.Operations, OCOperationDefinitionToOperationDefinition),
	}
}

func OCOperationDefinitionToOperationDefinition(ocOperationDefinition extapiv1.OperationDefinition) api.OperationDefinition {
	return api.OperationDefinition{
		Operation:            ocOperationDefinition.Operation,
		Description:          ocOperationDefinition.Description,
		Configuration:        ocOperationDefinition.Configuration,
		If:                   ocOperationDefinition.If,
		Unless:               ocOperationDefinition.Unless,
		FailWorkflowOnError:  ocOperationDefinition.FailWorkflowOnError,
		ErrorHandlerWorkflow: ocOperationDefinition.ErrorHandlerWorkflow,
		RetryStrategy:        OCWorkflowRetryStrategyToWorkflowRetryStrategy(ocOperationDefinition.RetryStrategy),
		MaxAttempts:          int64(ocOperationDefinition.MaxAttempts),
	}
}
