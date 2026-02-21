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
	"errors"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	oc "shio.solutions/tales.media/cli/pkg/opencast/client"
)

func New(ctx api.Context) (oc.Client, error) {
	// authentication
	authRequestOpts := []oc.RequestOpts{}
	switch {
	case ctx.Authentication.Basic != nil:
		authRequestOpts = append(authRequestOpts, oc.WithBasicAuth(
			ctx.Authentication.Basic.Username,
			ctx.Authentication.Basic.Password,
		))

	case ctx.Authentication.JWT != nil:
		authRequestOpts = append(authRequestOpts, oc.WithJWTHeader(
			ctx.Authentication.JWT.Header,
			ctx.Authentication.JWT.Prefix,
			ctx.Authentication.JWT.Token,
		))

	default:
		return nil, errors.New("client: no authentication configured")
	}

	// service mapper
	var sm oc.ServiceMapper
	switch {
	case ctx.ServiceMapper.Static != nil:
		sm = &oc.StaticServiceMapper{
			Default:     ctx.ServiceMapper.Static.Default,
			ServiceHost: ctx.ServiceMapper.Static.ServiceHosts,
		}

	case ctx.ServiceMapper.Dynamic != nil:
		smStatic := &oc.StaticServiceMapper{
			Default: ctx.ServiceMapper.Dynamic.ServiceRegistry,
		}

		staticClient, err := oc.New(smStatic, oc.WithRequestOptions(authRequestOpts...))
		if err != nil {
			return nil, err
		}

		sm = oc.NewDynamicServiceMapper(staticClient, ctx.ServiceMapper.Dynamic.TTL)

	default:
		return nil, errors.New("client: no service mapper configured")
	}

	return oc.New(sm, oc.WithRequestOptions(authRequestOpts...))
}
