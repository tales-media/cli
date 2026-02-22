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

	extapiv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11"
	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"
	oc "shio.solutions/tales.media/opencast-client-go/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api/conv"
)

type Series interface {
	List(context.Context, SeriesListRequest) ([]api.Series, error)
	Get(context.Context, SeriesGetRequest) (api.Series, error)
}

type SeriesListRequest struct {
	FilterByText         string
	FilterByID           string
	FilterByTitle        string
	FilterByDescription  string
	FilterByCreationDate string
	FilterByCreator      string
	FilterByContributors string
	FilterByOrganizers   string
	FilterByPublishers   string
	FilterByLanguage     string
	FilterByRightsHolder string
	FilterByLicense      string
	FilterBySubject      string
	FilterByManagedAcl   string
	SortBy               string // TODO: add conversion function
	SortDirection        api.SortDirection
}

type SeriesGetRequest struct {
	ID string
}

type opencastSeries struct {
	extAPI extapiclientv1.Client
}

var _ Series = &opencastSeries{}

func NewOpencastSeries(extAPI extapiclientv1.Client) Series {
	return &opencastSeries{
		extAPI: extAPI,
	}
}

func (svc *opencastSeries) List(ctx context.Context, req SeriesListRequest) ([]api.Series, error) {
	commonReqOpts := []oc.RequestOpts{
		extapiclientv1.WithSeriesOptions{
			WithACL:             false,
			OnlyWithWriteAccess: true,
		},
	}
	if req.FilterByText != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesTextFilterFilterKey: req.FilterByText},
		)
	}
	if req.FilterByID != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesIdentifierFilterKey: req.FilterByID},
		)
	}
	if req.FilterByTitle != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesTitleFilterKey: req.FilterByTitle},
		)
	}
	if req.FilterByDescription != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesDescriptionFilterKey: req.FilterByDescription},
		)
	}
	if req.FilterByCreationDate != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesCreationDateFilterKey: req.FilterByCreationDate},
		)
	}
	if req.FilterByCreator != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesCreatorFilterKey: req.FilterByCreator},
		)
	}
	if req.FilterByContributors != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesContributorsFilterKey: req.FilterByContributors},
		)
	}
	if req.FilterByOrganizers != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesOrganizersFilterKey: req.FilterByOrganizers},
		)
	}
	if req.FilterByPublishers != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesPublishersFilterKey: req.FilterByPublishers},
		)
	}
	if req.FilterByLanguage != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesLanguageFilterKey: req.FilterByLanguage},
		)
	}
	if req.FilterByRightsHolder != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesRightsHolderFilterKey: req.FilterByRightsHolder},
		)
	}
	if req.FilterByLicense != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesLicenseFilterKey: req.FilterByLicense},
		)
	}
	if req.FilterBySubject != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesSubjectFilterKey: req.FilterBySubject},
		)
	}
	if req.FilterByManagedAcl != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.SeriesManagedAclFilterKey: req.FilterByManagedAcl},
		)
	}
	if req.SortBy != "" {
		commonReqOpts = append(commonReqOpts, extapiclientv1.WithSort{{
			By:        extapiclientv1.SortKey(req.SortBy),
			Direction: conv.SortDirectionToOCSortDirection(req.SortDirection),
		}})
	}

	var ocSeries []extapiv1.Series
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			ocReq, err := svc.extAPI.ListSeriesRequest(
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
		oc.CollectAllPages(&ocSeries),
	)
	if err != nil {
		return nil, err
	}
	series := conv.Map(ocSeries, conv.OCSeriesToSeries)
	return series, nil
}

func (svc *opencastSeries) Get(ctx context.Context, req SeriesGetRequest) (api.Series, error) {
	ocSeries, _, err := svc.extAPI.GetSeries(
		ctx,
		req.ID,
		extapiclientv1.WithSeriesOptions{
			WithACL:             false,
			OnlyWithWriteAccess: true,
		},
	)
	if err != nil {
		return api.Series{}, err
	}
	return conv.OCSeriesToSeries(*ocSeries), nil
}

type talesSeries = opencastSeries

var _ Series = &talesSeries{}

var NewTalesSeries = NewOpencastSeries
