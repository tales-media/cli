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
	"github.com/spf13/cobra"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

func configCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Manage Configuration",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = AdminGroup.ID
	cmd.AddGroup(
		ResourcesGroup,
		ManagementGroup,
	)
	cmd.AddCommand(
		// resources
		configContextCommand(cfg),

		// management
		configGetCommand(cfg),
	)
	return cmd
}

func configGetCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"get",
		"Get Configuration",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			config := api.DeepCopy(&cfg.Config)
			config.RedactCredentials() // TODO: add flag to no redact credentials?
			return *config, nil
		},
	)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
