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

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
	extapiv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/cli/pkg/opencast/apis/external-api/v1.11/client"
	oc "shio.solutions/tales.media/cli/pkg/opencast/client"
)

type Event interface {
	List(context.Context, EventListRequest) ([]api.Event, error)
	Get(context.Context, EventGetRequest) (api.Event, error)
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

type talesEvent = opencastEvent

var _ Event = &talesEvent{}

var NewTalesEvent = NewOpencastEvent
