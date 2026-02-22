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
	"reflect"
	"slices"

	"github.com/spf13/cobra"

	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"
	oc "shio.solutions/tales.media/opencast-client-go/client"

	"shio.solutions/tales.media/cli/internal/pkg/formatter"
	"shio.solutions/tales.media/cli/internal/talesctl/svc"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/client"
)

func baseCommand(use, short string, valueFunc func(*cobra.Command, []string) (any, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   use,
		Short:                 short,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// at this point (after args parsing) we handle error output with formatter
			cmd.SilenceErrors = true
			formatter := detectFormatter(cmd)

			val, err := valueFunc(cmd, args)
			if err != nil {
				fmtErr := formatter.Error(cmd.ErrOrStderr(), err)
				if fmtErr != nil {
					panic(fmtErr)
				}
				return err
			}

			if val == nil {
				return nil
			}

			valType := reflect.TypeOf(val)
			switch valType.Kind() {
			case reflect.Slice:
				return formatter.List(cmd.OutOrStdout(), val)
			case reflect.Struct:
				return formatter.Object(cmd.OutOrStdout(), val)
			default:
				panic(fmt.Sprintf("BUG: cannot format values of type %s", valType.Name()))
			}
		},
	}
	return cmd
}

func cfgCommand(use, short string, cfg *Config, configValueFunc func(*cobra.Command, []string) (any, error)) *cobra.Command {
	return baseCommand(use, short, func(cmd *cobra.Command, args []string) (any, error) {
		var (
			s   svc.Config
			req svc.ConfigGetRequest
		)

		// TODO: allow setting configuration file path
		mustSelect(cfg.AliasType, map[AliasType]func(){
			OpencastAlias: func() { s = svc.NewOpencastConfig() },
			TalesAlias:    func() { s = svc.NewTalesConfig() },
		})()

		config, err := s.Get(cmd.Context(), req)
		if err != nil {
			return nil, err
		}
		cfg.Config = config

		// determine context
		cfg.ContextName = getContextFlag(cmd.Flags())
		if cfg.ContextName == "" {
			cfg.ContextName = cfg.CurrentContext
		}

		return configValueFunc(cmd, args)
	})
}

func occCommand(use, short string, cfg *Config, occValueFunc func(*cobra.Command, []string, oc.Client) (any, error)) *cobra.Command {
	return cfgCommand(use, short, cfg, func(cmd *cobra.Command, args []string) (any, error) {
		if cfg.ContextName == "" {
			return nil, errors.New("cli: no Opencast context selected")
		}
		i := slices.IndexFunc(cfg.Contexts, func(c api.Context) bool { return c.Name == cfg.ContextName })
		if i < 0 {
			return nil, fmt.Errorf("cli: current Opencast context '%s' not found", cfg.ContextName)
		}
		cfg.Context = &cfg.Contexts[i]

		c, err := client.New(*cfg.Context)
		if err != nil {
			return nil, err
		}

		return occValueFunc(cmd, args, c)
	})
}

func extAPICommand(use, short string, cfg *Config, extAPIClientValueFunc func(*cobra.Command, []string, extapiclientv1.Client) (any, error)) *cobra.Command {
	return occCommand(use, short, cfg, func(cmd *cobra.Command, args []string, occ oc.Client) (any, error) {
		extAPI := extapiclientv1.New(occ)
		return extAPIClientValueFunc(cmd, args, extAPI)
	})
}

func detectFormatter(cmd *cobra.Command) formatter.Formatter {
	switch getOutputFlag(cmd.Flags()) {
	default:
		fallthrough
	case HumanOutput:
		return &formatter.Human{}
	case WideOutput:
		return &formatter.Human{Wide: true}
	case JSONOutput:
		return &formatter.JSON{}
	case YAMLOutput:
		return &formatter.YAML{}
	case NoneOutput:
		return &formatter.None{}
	}
}
