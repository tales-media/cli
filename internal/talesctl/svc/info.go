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

	"github.com/tales-media/cli/internal/talesctl/svc/api"
	"github.com/tales-media/cli/internal/talesctl/svc/api/conv"
	extapiclientv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11/client"
)

type Info interface {
	API(context.Context) (api.APIInfo, error)
	Me(context.Context) (api.Me, error)
	Organization(context.Context) (api.Organization, error)
}

type opencastInfo struct {
	extAPI extapiclientv1.Client
}

var _ Info = &opencastInfo{}

func NewOpencastInfo(extAPI extapiclientv1.Client) Info {
	return &opencastInfo{
		extAPI: extAPI,
	}
}

func (svc *opencastInfo) API(ctx context.Context) (api.APIInfo, error) {
	ocAPI, _, err := svc.extAPI.GetAPI(ctx)
	if err != nil {
		return api.APIInfo{}, err
	}
	ocVersion, _, err := svc.extAPI.GetAPIVersion(ctx)
	if err != nil {
		return api.APIInfo{}, err
	}
	return conv.OCAPIToAPIInfo(*ocAPI, *ocVersion), nil
}

func (svc *opencastInfo) Me(ctx context.Context) (api.Me, error) {
	ocMe, _, err := svc.extAPI.GetInfoMe(ctx)
	if err != nil {
		return api.Me{}, err
	}
	roles, _, err := svc.extAPI.GetInfoMeRoles(ctx)
	if err != nil {
		return api.Me{}, err
	}
	return conv.OCMeToMe(*ocMe, roles), nil
}

func (svc *opencastInfo) Organization(ctx context.Context) (api.Organization, error) {
	ocOrganization, _, err := svc.extAPI.GetInfoOrganization(ctx)
	if err != nil {
		return api.Organization{}, err
	}
	orgProperties, _, err := svc.extAPI.GetInfoOrganizationProperties(ctx)
	if err != nil {
		return api.Organization{}, err
	}
	return conv.OCOrganizationToOrganization(*ocOrganization, orgProperties), nil
}

type talesInfo = opencastInfo

var _ Info = &talesInfo{}

var NewTalesInfo = NewOpencastInfo
