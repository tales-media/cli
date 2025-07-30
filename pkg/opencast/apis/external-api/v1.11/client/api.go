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

package client

import (
	"context"
	"net/http"

	extapiv1 "github.com/tales-media/cli/pkg/opencast/apis/external-api/v1.11"
	oc "github.com/tales-media/cli/pkg/opencast/client"
)

const (
	AcceptJSONHeader = "application/" + extapiv1.Version + "+json"
)

type Client interface {
	Do(*oc.Request) (*oc.Response, error)

	// API

	GetAPI(context.Context, ...oc.RequestOpts) (*extapiv1.API, *oc.Response, error)
	GetAPIRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetAPIVersion(context.Context, ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error)
	GetAPIVersionRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetAPIVersionDefault(context.Context, ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error)
	GetAPIVersionDefaultRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	// Info

	GetInfoOrganization(context.Context, ...oc.RequestOpts) (*extapiv1.Organization, *oc.Response, error)
	GetInfoOrganizationRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetInfoOrganizationProperties(context.Context, ...oc.RequestOpts) (*extapiv1.OrganizationProperties, *oc.Response, error)
	GetInfoOrganizationPropertiesRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetInfoOrganizationPropertiesEngageUIURL(context.Context, ...oc.RequestOpts) (*extapiv1.OrganizationProperties, *oc.Response, error)
	GetInfoOrganizationPropertiesEngageUIURLRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetInfoMe(context.Context, ...oc.RequestOpts) (*extapiv1.Me, *oc.Response, error)
	GetInfoMeRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)

	GetInfoMeRoles(context.Context, ...oc.RequestOpts) (*extapiv1.StringList, *oc.Response, error)
	GetInfoMeRolesRequest(context.Context, ...oc.RequestOpts) (*oc.Request, error)
}

type client struct {
	occ oc.Client
}

var _ Client = &client{}

func New(opencastClient oc.Client) *client {
	return &client{
		occ: opencastClient,
	}
}

func (c *client) Do(req *oc.Request) (*oc.Response, error) {
	if err := req.ApplyOptions(
		oc.WithHeader("Accept", AcceptJSONHeader),
	); err != nil {
		return nil, err
	}
	return c.occ.Do(req)
}

func (c *client) GetAPI(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.API, *oc.Response, error) {
	return oc.GenericAutoDecodedCall[extapiv1.API](
		func() (*oc.Request, error) { return c.GetAPIRequest(ctx, opts...) },
		c.Do,
	)
}

func (c *client) GetAPIRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		extapiv1.ServiceType,
		"/api/",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetAPIVersion(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error) {
	return oc.GenericAutoDecodedCall[extapiv1.APIVersion](
		func() (*oc.Request, error) { return c.GetAPIVersionRequest(ctx, opts...) },
		c.Do,
	)
}

func (c *client) GetAPIVersionRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		extapiv1.ServiceType,
		"/api/version",
		oc.NoBody,
		opts...,
	)
}

func (c *client) GetAPIVersionDefault(ctx context.Context, opts ...oc.RequestOpts) (*extapiv1.APIVersion, *oc.Response, error) {
	return oc.GenericAutoDecodedCall[extapiv1.APIVersion](
		func() (*oc.Request, error) { return c.GetAPIVersionDefaultRequest(ctx, opts...) },
		c.Do,
	)
}

func (c *client) GetAPIVersionDefaultRequest(ctx context.Context, opts ...oc.RequestOpts) (*oc.Request, error) {
	return oc.NewRequest(
		ctx,
		http.MethodGet,
		extapiv1.ServiceType,
		"/api/default",
		oc.NoBody,
		opts...,
	)
}
