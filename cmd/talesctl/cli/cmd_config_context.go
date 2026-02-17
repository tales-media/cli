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

package cli

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"shio.solutions/tales.media/cli/internal/talesctl/svc"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

// TODO: update

func configContextCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context",
		Short:                 "Manage Contexts",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ManagementGroup,
	)
	cmd.AddCommand(
		configContextCreateCommand(cfg),
		configContextListCommand(cfg),
		configContextGetCommand(cfg),
		configContextUseCommand(cfg),
		configContextDeleteCommand(cfg),
	)
	return cmd
}

func configContextCreateCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"create [context name] [url]",
		"Create a Context",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			var (
				s   svc.Config
				req svc.ConfigUpdateRequest
			)

			// TODO: allow setting configuration file path
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastConfig() },
				TalesAlias:    func() { s = svc.NewTalesConfig() },
			})()

			ctx := api.Context{
				Name: args[0],
			}

			// check if context already exists
			found := -1 != slices.IndexFunc(cfg.Contexts, func(c api.Context) bool { return c.Name == ctx.Name })
			if found {
				return nil, fmt.Errorf("cli: context '%s' already exists", ctx.Name)
			}

			// service mapper
			switch getContextServiceMapperFlag(cmd.Flags()) {
			case StaticContextServiceMapper:
				serviceHosts, err := getContextStaticServiceMappingFlag(cmd.Flags())
				if err != nil {
					return nil, err
				}
				ctx.ServiceMapper.Static = &api.StaticServiceMapper{
					Default:      args[1],
					ServiceHosts: serviceHosts,
				}

			case DynamicContextServiceMapper:
				ctx.ServiceMapper.Dynamic = &api.DynamicServiceMapper{
					ServiceRegistry: args[1],
					TTL:             getContextDynamicServiceMapperTTLFlag(cmd.Flags()),
				}
			}

			// authentication
			username, password, err := getContextBasicAuthFlag(cmd.Flags())
			if err != nil {
				return nil, err
			}
			if username != "" || password != "" {
				ctx.Authentication.Basic = &api.BasicAuthentication{
					Username: username,
					Password: password,
				}
			}

			token := getContextJWTAuthFlag(cmd.Flags())
			if token != "" {
				ctx.Authentication.JWT = &api.JWTAuthentication{
					Token: token,
				}
			}

			req.Config = cfg.Config
			req.Config.Contexts = append(req.Config.Contexts, ctx)
			if getContextUseFlag(cmd.Flags()) {
				req.Config.CurrentContext = ctx.Name
			}

			err = s.Update(cmd.Context(), req)
			return nil, err
		},
	)
	cmd.Args = cobra.ExactArgs(2)
	addContextUseFlag(cmd.Flags())
	addContextServiceMapperFlag(cmd.Flags())
	addContextStaticServiceMappingFlag(cmd.Flags())
	addContextDynamicServiceMapperTTLFlag(cmd.Flags())
	addContextBasicAuthFlag(cmd.Flags())
	addContextJWTAuthFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func configContextListCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"list",
		"List Contexts",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			config := api.DeepCopy(&cfg.Config)
			config.RedactCredentials() // TODO: add flag to no redact credentials?

			for i := range config.Contexts {
				ctx := &config.Contexts[i]
				if ctx.Name == cfg.ContextName {
					ctx.Name = ctx.Name + " (*)"
				}
			}

			// TODO: use custom column for marking current context

			return config.Contexts, nil
		},
	)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func configContextGetCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"get [context name]",
		"Get a Context",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			ctxName := cfg.ContextName
			if len(args) > 0 {
				ctxName = args[0]
			}

			var ctx *api.Context
			for _, c := range cfg.Contexts {
				if c.Name == ctxName {
					ctx = api.DeepCopy(&c)
				}
			}
			if ctx == nil {
				return nil, errors.New("Not Found")
			}

			ctx.RedactCredentials() // TODO: add flag to no redact credentials?

			return *ctx, nil
		},
	)
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func configContextUseCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"use [context name]",
		"Set Current Context",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			var (
				s   svc.Config
				req svc.ConfigUpdateRequest
			)

			// TODO: allow setting configuration file path
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastConfig() },
				TalesAlias:    func() { s = svc.NewTalesConfig() },
			})()

			req.Config = cfg.Config
			req.Config.CurrentContext = args[0]

			// check if new current context exists
			found := -1 != slices.IndexFunc(cfg.Contexts, func(c api.Context) bool { return c.Name == req.Config.CurrentContext })
			if !found {
				return nil, fmt.Errorf("cli: context '%s' not found", req.Config.CurrentContext)
			}

			return nil, s.Update(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func configContextDeleteCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"delete [context name]",
		"Delete a Context",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			var (
				s   svc.Config
				req svc.ConfigUpdateRequest
			)

			// TODO: allow setting configuration file path
			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastConfig() },
				TalesAlias:    func() { s = svc.NewTalesConfig() },
			})()

			ctxDelName := args[0]

			if cfg.CurrentContext == ctxDelName {
				return nil, fmt.Errorf("cli: cannot delete current context '%s'", ctxDelName)
			}

			req.Config = cfg.Config
			req.Config.Contexts = make([]api.Context, 0, len(cfg.Contexts))

			found := 0
			for _, ctx := range cfg.Contexts {
				if ctx.Name == ctxDelName {
					found += 1
					continue
				}
				req.Config.Contexts = append(req.Config.Contexts, ctx)
			}
			if found == 0 {
				return nil, errors.New("Not Found")
			}
			// TODO: handle found > 1?

			return nil, s.Update(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
