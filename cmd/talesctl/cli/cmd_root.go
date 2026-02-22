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

import "github.com/spf13/cobra"

var (
	ManagementGroup = &cobra.Group{ID: "management", Title: "Management Commands:"}
	ResourcesGroup  = &cobra.Group{ID: "resources", Title: "Resources Commands:"}
	AdminGroup      = &cobra.Group{ID: "admin", Title: "Admin Commands:"}
)

func New(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: cfg.Alias,
		Short: mustSelect(cfg.AliasType, map[AliasType]string{
			OpencastAlias: "Opencast CLI",
			TalesAlias:    "tales.media CLI",
		}),
		Long: mustSelect(cfg.AliasType, map[AliasType]string{
			OpencastAlias: "A command-line interface for Opencast",
			TalesAlias:    "A command-line interface for tales.media",
		}),
		TraverseChildren:      true,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		CompletionOptions: cobra.CompletionOptions{
			// TODO: add support for completions
			DisableDefaultCmd: true,
		},
	}
	cmd.SetArgs(cfg.Args)

	// global flags
	addContextFlag(cmd.PersistentFlags())
	addOutputFlag(cmd.PersistentFlags())

	// commands

	cmd.AddGroup(
		ResourcesGroup,
		AdminGroup,
	)

	cmd.AddCommand(
		// resources
		agentCommand(cfg),
		eventCommand(cfg),
		groupCommand(cfg),
		playlistCommand(cfg),
		seriesCommand(cfg),
		workflowCommand(cfg),
		workflowDefinitionCommand(cfg),

		// admin
		configCommand(cfg),

		// additional
		infoCommand(cfg),
		versionCommand(cfg),
	)

	// tales.media specific commands
	if cfg.AliasType == TalesAlias {
		cmd.AddCommand(
			// admin
			authCommand(cfg),
		)
	}

	return cmd
}
