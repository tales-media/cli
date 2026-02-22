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

type Group interface {
	List(context.Context, GroupListRequest) ([]api.Group, error)
	Get(context.Context, GroupGetRequest) (api.Group, error)
}

type GroupListRequest struct {
	FilterByName  string
	SortBy        string // TODO: add conversion function
	SortDirection api.SortDirection
}

type GroupGetRequest struct {
	ID string
}

type opencastGroup struct {
	extAPI extapiclientv1.Client
}

var _ Group = &opencastGroup{}

func NewOpencastGroup(extAPI extapiclientv1.Client) Group {
	return &opencastGroup{
		extAPI: extAPI,
	}
}

func (svc *opencastGroup) List(ctx context.Context, req GroupListRequest) ([]api.Group, error) {
	commonReqOpts := []oc.RequestOpts{}
	if req.FilterByName != "" {
		commonReqOpts = append(commonReqOpts,
			extapiclientv1.WithFilter{extapiclientv1.GroupNameFilterKey: req.FilterByName},
		)
	}
	if req.SortBy != "" {
		commonReqOpts = append(commonReqOpts, extapiclientv1.WithSort{{
			By:        extapiclientv1.SortKey(req.SortBy),
			Direction: conv.SortDirectionToOCSortDirection(req.SortDirection),
		}})
	}

	var ocGroups []extapiv1.Group
	err := oc.Paginate(
		svc.extAPI,
		func(i int) (*oc.Request, error) {
			ocReq, err := svc.extAPI.ListGroupRequest(
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
		oc.CollectAllPages(&ocGroups),
	)
	if err != nil {
		return nil, err
	}

	groups := conv.Map(ocGroups, conv.OCGroupToGroup)
	return groups, nil
}

func (svc *opencastGroup) Get(ctx context.Context, req GroupGetRequest) (api.Group, error) {
	ocGroup, _, err := svc.extAPI.GetGroup(ctx, req.ID)
	if err != nil {
		return api.Group{}, err
	}
	// BUG: this endpoint returns broken data (roles with [] and members as Java object toString())
	return conv.OCGroupToGroup(*ocGroup), nil
}

type talesGroup = opencastGroup

var _ Group = &talesGroup{}

var NewTalesGroup = NewOpencastGroup
