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
	cmd.PersistentFlags().String(ContextFlag, "", "the name of the Opencast context to use")
	outputValue := OutputValue()
	cmd.PersistentFlags().VarP(outputValue, OutputFlag, OutputFlagShort, outputValue.Usage("the output format"))

	// commands

	cmd.AddGroup(
		ManagementGroup,
		ResourcesGroup,
		AdminGroup,
	)

	cmd.AddCommand(
		agentCommand(cfg),
		workflowDefinitionCommand(cfg),
		infoCommand(cfg),
		versionCommand(cfg),
	)

	return cmd
}
