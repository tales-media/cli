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
	"errors"
	"io"
	"slices"
	"time"

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/base"
	"shio.solutions/tales.media/opencast-client-go/apis/meta/objlist"
	oc "shio.solutions/tales.media/opencast-client-go/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
)

type Event interface {
	Create(context.Context, EventCreateRequest) ([]api.Event, error)
	List(context.Context, EventListRequest) ([]api.Event, error)
	Get(context.Context, EventGetRequest) (api.Event, error)
}

type EventCreateRequest struct {
	// ACL (Opencast)

	ACL []api.ACE

	// ACL (tales.media)

	TalesACLPreset     api.TalesACLPreset
	TalesACLUsersRead  []string
	TalesACLUsersWrite []string

	// Metadata

	MetadataID          string
	MetadataTitle       string
	MetadataDescription string
	MetadataSeriesID    string
	// TODO: CreationDate
	MetadataStartDate time.Time
	MetadataDuration  time.Duration
	// TODO: Creator
	// TODO: Contributors
	// TODO: Presenters
	MetadataLocation string
	// TODO: Language
	MetadataRightsHolder string
	// TODO: License
	MetadataSubjects []string
	MetadataSource   string

	// Workflow

	WorkflowDefinition string
	WorkflowProperties map[string]string

	// Scheduling

	SchedulingAgentName   *string
	SchedulingAgentInputs []string
	SchedulingStartDate   *time.Time
	SchedulingDuration    *time.Duration
	SchedulingRRule       *string

	// Upload

	PresenterStream            io.ReadCloser
	PresenterStreamFilename    *string
	PresentationStream         io.ReadCloser
	PresentationStreamFilename *string
	AudioStream                io.ReadCloser
	AudioStreamFilename        *string
}

type EventListRequest struct {
	FilterByText               string
	FilterByID                 string
	FilterByStatus             api.EventStatus
	FilterByTitle              string
	FilterByDescription        string
	FilterBySeries             string
	FilterBySeriesID           string
	FilterByCreationDate       string
	FilterByStartDate          string
	FilterByContributors       string
	FilterByPresenters         string
	FilterByLocation           string
	FilterByLanguage           string
	FilterByRightsHolder       string
	FilterByLicense            string
	FilterBySubjects           string
	FilterBySource             string
	FilterByScheduledStartDate string
	FilterByAgentName          string
	SortBy                     string // TODO: add conversion function
	SortDirection              api.SortDirection
}

type EventGetRequest struct {
	ID string
}

type opencastEvent struct {
	extAPI extapiclientv1.Client
}

var _ Event = &opencastEvent{}

func NewOpencastEvent(extAPI extapiclientv1.Client) Event {
	return &opencastEvent{
		extAPI: extAPI,
	}
}

func (svc *opencastEvent) Create(ctx context.Context, req EventCreateRequest) ([]api.Event, error) {
	// ACL

	ocACL := extapiv1.ACL(conv.Map(req.ACL, conv.ACEToOCACE))

	// Metadata

	ocDublinCoreEpisodeCatalog := extapiv1.Catalog{
		Flavor: base.DublinCoreEpisodeFlavor,
		Fields: []extapiv1.Field{},
	}
	if req.MetadataID != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.IdentifierFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataID),
		})
	}
	if req.MetadataTitle != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.TitleFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataTitle),
		})
	}
	if req.MetadataDescription != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.DescriptionFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataDescription),
		})
	}
	if req.MetadataSeriesID != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.IsPartOfFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataSeriesID),
		})
	}
	if !req.MetadataStartDate.IsZero() {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.StartDateFieldID,
			Value: extapiv1.DateTimeFieldValue(base.DateTime{Time: req.MetadataStartDate}),
		})
	}
	// TODO: MetadataDuration
	// if req.MetadataDuration != 0 {
	// 	ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
	// 		ID:    extapiv1.DurationFieldID,
	// 		Value: extapiv1.DateTimeFieldValue(base.DateTime{Time: req.MetadataDuration}),
	// 	})
	// }
	if req.MetadataLocation != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.LocationFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataLocation),
		})
	}
	if req.MetadataRightsHolder != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.RightsHolderFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataRightsHolder),
		})
	}
	if len(req.MetadataSubjects) > 0 {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.SubjectsFieldID,
			Value: extapiv1.MixedTextFieldValue(req.MetadataSubjects), // TODO: is this correct?
		})
	}
	if req.MetadataSource != "" {
		ocDublinCoreEpisodeCatalog.Fields = append(ocDublinCoreEpisodeCatalog.Fields, extapiv1.Field{
			ID:    extapiv1.SourceFieldID,
			Value: extapiv1.TextFieldValue(req.MetadataSource),
		})
	}

	// Workflow

	ocProcessing := &extapiv1.Processing{
		Workflow:      req.WorkflowDefinition,
		Configuration: req.WorkflowProperties,
	}

	// Scheduling

	var ocScheduling *extapiv1.SchedulingRequest
	// TODO: implement

	// Event

	ocEvent := &extapiclientv1.CreateEventRequestBody{
		ACL:        ocACL,
		Metadata:   []extapiv1.Catalog{ocDublinCoreEpisodeCatalog},
		Scheduling: ocScheduling,
		Processing: ocProcessing,
	}

	// Upload

	if req.PresenterStreamFilename != nil {
		ocEvent.Scheduling = nil
		ocEvent.PresenterStreamFilename = *req.PresenterStreamFilename
		ocEvent.PresenterStream = req.PresenterStream
	}
	if req.PresentationStreamFilename != nil {
		ocEvent.Scheduling = nil
		ocEvent.PresentationStreamFilename = *req.PresentationStreamFilename
		ocEvent.PresentationStream = req.PresentationStream
	}
	if req.AudioStreamFilename != nil {
		ocEvent.Scheduling = nil
		ocEvent.AudioStreamFilename = *req.AudioStreamFilename
		ocEvent.AudioStream = req.AudioStream
	}

	// create

	resp, _, err := svc.extAPI.CreateEvent(ctx, ocEvent)
	if err != nil {
		return nil, err
	}

	// retrieve created events

	var ids []extapiv1.Identifier
	switch resp.Type {
	case objlist.Object:
		ids = []extapiv1.Identifier{resp.ObjectVal}
	case objlist.List:
		ids = resp.ListVal
	}

	events := make([]api.Event, 0, len(ids))
	for _, id := range ids {
		event, err := svc.Get(ctx, EventGetRequest{ID: id.Identifier})
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (svc *opencastEvent) List(ctx context.Context, req EventListRequest) ([]api.Event, error) {
	commonReqOpts := []oc.RequestOpts{
		extapiclientv1.WithEventOptions{
			WithACL:                    false,
			WithMetadata:               false,
			WithScheduling:             true,
			WithPublications:           false,
			IncludeInternalPublication: false,
			OnlyWithWriteAccess:        true,
		},
	}
	if req.FilterByText != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventTextFilterFilterKey: req.FilterByText},
		)
	}
	if req.FilterByID != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventIdentifierFilterKey: req.FilterByID},
		)
	}
	if req.FilterByStatus != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{
				extapiclientv1.EventStatusFilterKey: string(conv.EventStatusToOCEventStatus(req.FilterByStatus)),
			},
		)
	}
	if req.FilterByTitle != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventTitleFilterKey: req.FilterByTitle},
		)
	}
	if req.FilterByDescription != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventDescriptionFilterKey: req.FilterByDescription},
		)
	}
	if req.FilterBySeries != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventSeriesFilterKey: req.FilterBySeries},
		)
	}
	if req.FilterBySeriesID != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventIsPartOfFilterKey: req.FilterBySeriesID},
		)
	}
	if req.FilterByCreationDate != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventCreatedFilterKey: req.FilterByCreationDate},
		)
	}
	if req.FilterByStartDate != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventStartFilterKey: req.FilterByStartDate},
		)
	}
	if req.FilterByContributors != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventContributorsFilterKey: req.FilterByContributors},
		)
	}
	if req.FilterByPresenters != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventPresentersFilterKey: req.FilterByPresenters},
		)
	}
	if req.FilterByLocation != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventLocationFilterKey: req.FilterByLocation},
		)
	}
	if req.FilterByLanguage != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventLanguageFilterKey: req.FilterByLanguage},
		)
	}
	if req.FilterByRightsHolder != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventRightsHolderFilterKey: req.FilterByRightsHolder},
		)
	}
	if req.FilterByLicense != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventLicenseFilterKey: req.FilterByLicense},
		)
	}
	if req.FilterBySubjects != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventSubjectFilterKey: req.FilterBySubjects},
		)
	}
	if req.FilterBySource != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventSourceFilterKey: req.FilterBySource},
		)
	}
	if req.FilterByScheduledStartDate != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventTechnicalStartFilterKey: req.FilterByScheduledStartDate},
		)
	}
	if req.FilterByAgentName != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.EventAgentIdFilterKey: req.FilterByAgentName},
		)
	}
	if req.SortBy != "" {
		commonReqOpts = append(commonReqOpts, extapiclientv1.WithSort{{
			By:        extapiclientv1.SortKey(req.SortBy),
			Direction: conv.SortDirectionToOCSortDirection(req.SortDirection),
		}})
	}

	var ocEvents []extapiv1.Event
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			ocReq, err := svc.extAPI.ListEventRequest(
				ctx,
				extapiclientv1.WithPagination{
					Limit:  PageSize,
					Offset: i * PageSize,
				},
			)
			if err != nil {
				return nil, err
			}
			if err = ocReq.ApplyOptions(commonReqOpts...); err != nil {
				return nil, err
			}
			return ocReq, nil
		},
		oc.CollectAllPages(&ocEvents),
	)
	if err != nil {
		return nil, err
	}
	events := conv.Map(ocEvents, conv.OCEventToEvent)
	return events, nil
}

func (svc *opencastEvent) Get(ctx context.Context, req EventGetRequest) (api.Event, error) {
	ocEvent, _, err := svc.extAPI.GetEvent(
		ctx,
		req.ID,
		extapiclientv1.WithEventOptions{
			WithACL:                    false,
			WithMetadata:               false,
			WithScheduling:             true,
			WithPublications:           false,
			IncludeInternalPublication: false,
			OnlyWithWriteAccess:        true,
		},
	)
	if err != nil {
		return api.Event{}, err
	}
	return conv.OCEventToEvent(*ocEvent), nil
}

type talesEvent struct {
	opencastEvent

	infoSvc Info
}

var _ Event = &talesEvent{}

func NewTalesEvent(extAPI extapiclientv1.Client) Event {
	return &talesEvent{
		opencastEvent: opencastEvent{
			extAPI: extAPI,
		},
		infoSvc: NewTalesInfo(extAPI),
	}
}

func (svc *talesEvent) Create(ctx context.Context, req EventCreateRequest) ([]api.Event, error) {
	org, err := svc.infoSvc.Organization(ctx)
	if err != nil {
		return nil, err
	}

	me, err := svc.infoSvc.Me(ctx)
	if err != nil {
		return nil, err
	}

	// ACL

	acl := conv.TalesACLPresetToACL(req.TalesACLPreset, org.ID)
	acl = slices.Grow(acl, 4+ // owner roles
		4+ // user roles
		len(req.TalesACLUsersRead)*2+ // extra read users
		len(req.TalesACLUsersRead)*4, // extra write users
	)

	ownerRolePrefix := "ROLE_OWNER_"
	userRolePrefix := "ROLE_USER_"

	acl = append(acl,
		api.ACE{Role: ownerRolePrefix + me.Username, Action: api.ReadAction, Allow: true},
		api.ACE{Role: ownerRolePrefix + me.Username, Action: api.WriteAction, Allow: true},
		api.ACE{Role: ownerRolePrefix + me.Username, Action: api.AnnotateAction, Allow: true},
		api.ACE{Role: ownerRolePrefix + me.Username, Action: api.AnnotateAdminAction, Allow: true},

		api.ACE{Role: userRolePrefix + me.Username, Action: api.ReadAction, Allow: true},
		api.ACE{Role: userRolePrefix + me.Username, Action: api.WriteAction, Allow: true},
		api.ACE{Role: userRolePrefix + me.Username, Action: api.AnnotateAction, Allow: true},
		api.ACE{Role: userRolePrefix + me.Username, Action: api.AnnotateAdminAction, Allow: true},
	)

	for _, username := range req.TalesACLUsersRead {
		acl = append(acl,
			api.ACE{Role: userRolePrefix + username, Action: api.ReadAction, Allow: true},
			api.ACE{Role: userRolePrefix + username, Action: api.AnnotateAction, Allow: true},
		)
	}

	for _, username := range req.TalesACLUsersWrite {
		acl = append(acl,
			api.ACE{Role: userRolePrefix + username, Action: api.ReadAction, Allow: true},
			api.ACE{Role: userRolePrefix + username, Action: api.WriteAction, Allow: true},
			api.ACE{Role: userRolePrefix + username, Action: api.AnnotateAction, Allow: true},
			api.ACE{Role: userRolePrefix + username, Action: api.AnnotateAdminAction, Allow: true},
		)
	}

	req.ACL = acl

	// Metadata

	if req.MetadataID != "" {
		return nil, errors.New("setting event ID is not allowed in tales.media")
	}
	if req.MetadataSeriesID == "" {
		return nil, errors.New("series ID is required in tales.media")
	}

	return svc.opencastEvent.Create(ctx, req)
}
