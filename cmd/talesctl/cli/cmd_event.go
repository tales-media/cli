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
	"io"

	"github.com/spf13/cobra"

	extapiclientv1 "shio.solutions/tales.media/opencast-client-go/apis/external-api/v1.11/client"

	"shio.solutions/tales.media/cli/internal/talesctl/svc"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

// TODO: delete
// TODO: update

func eventCommand(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "event",
		Short:                 "Manage Events",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.GroupID = ResourcesGroup.ID
	cmd.AddGroup(
		ResourcesGroup,
		ManagementGroup,
	)
	cmd.AddCommand(
		// resources
		eventACLCommand(cfg),
		eventTrackCommand(cfg),
		eventCatalogCommand(cfg),
		eventPublicationCommand(cfg),

		// management
		eventCreateCommand(cfg),
		eventListCommand(cfg),
		eventGetCommand(cfg),
	)
	return cmd
}

func eventCreateCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"create",
		"Create an Event",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Event
				req svc.EventCreateRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastEvent(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesEvent(extAPI) },
			})()

			// common configuration

			// Metadata

			req.MetadataTitle = getMetadataTitleFlag(cmd.Flags())
			req.MetadataDescription = getMetadataDescriptionFlag(cmd.Flags())
			req.MetadataSeriesID = getMetadataSeriesIDFlag(cmd.Flags())
			req.MetadataStartDate = getMetadataStartDateFlag(cmd.Flags())
			req.MetadataDuration = getMetadataDurationFlag(cmd.Flags())
			req.MetadataLocation = getMetadataLocationFlag(cmd.Flags())
			req.MetadataRightsHolder = getMetadataRightsHolderFlag(cmd.Flags())
			req.MetadataSubjects = getMetadataSubjectFlag(cmd.Flags())
			req.MetadataSource = getMetadataSourceFlag(cmd.Flags())

			// TODO: Scheduling

			// specific configuration

			err := mustSelect(cfg.AliasType, map[AliasType]func() error{
				OpencastAlias: func() error {
					// ACL

					acl, err := getACEFlag(cmd.Flags())
					if err != nil {
						return err
					}
					req.ACL = acl

					// Metadata

					req.MetadataID = getMetadataIDFlag(cmd.Flags())

					// Workflow

					req.WorkflowDefinition = getWorkflowDefinitionFlag(cmd.Flags())
					req.WorkflowProperties, err = getWorkflowPropertiesFlag(cmd.Flags())
					if err != nil {
						return err
					}

					// Upload

					var (
						fn string
						r  io.ReadCloser
					)

					fn, r, err = getTrackXFlag("presenter", cmd)
					if err != nil {
						return err
					}
					if fn != "" {
						req.PresenterStreamFilename = new(fn)
						req.PresenterStream = r
					}

					fn, r, err = getTrackXFlag("presentation", cmd)
					if err != nil {
						return err
					}
					if fn != "" {
						req.PresentationStreamFilename = new(fn)
						req.PresentationStream = r
					}

					fn, r, err = getTrackXFlag("audio", cmd)
					if err != nil {
						return err
					}
					if fn != "" {
						req.AudioStreamFilename = new(fn)
						req.AudioStream = r
					}

					return nil
				},
				TalesAlias: func() error {
					// ACL

					req.TalesACLPreset = getACLPresetFlag(cmd.Flags())
					req.TalesACLUsersRead = getACLUsersReadFlag(cmd.Flags())
					req.TalesACLUsersWrite = getACLUsersWriteFlag(cmd.Flags())

					// Workflow

					req.WorkflowDefinition = "tales-media-main-ingest-generic"

					// Upload

					var (
						fn  string
						r   io.ReadCloser
						err error
					)

					fn, r, err = getTrackXFlag("main", cmd)
					if err != nil {
						return err
					}
					if fn != "" {
						req.PresenterStreamFilename = new(fn)
						req.PresenterStream = r
					}

					fn, r, err = getTrackXFlag("secondary", cmd)
					if err != nil {
						return err
					}
					if fn != "" {
						req.PresentationStreamFilename = new(fn)
						req.PresentationStream = r
					}

					return nil
				},
			})()
			if err != nil {
				return nil, err
			}

			return s.Create(cmd.Context(), req)
		},
	)
	cmd.GroupID = ManagementGroup.ID

	// common flags

	addMetadataTitleFlag(cmd.Flags())
	addMetadataDescriptionFlag(cmd.Flags())
	addMetadataSeriesIDFlag(cmd.Flags())
	addMetadataStartDateFlag(cmd.Flags())
	addMetadataDurationFlag(cmd.Flags())
	addMetadataLocationFlag(cmd.Flags())
	addMetadataRightsHolderFlag(cmd.Flags())
	addMetadataSubjectFlag(cmd.Flags())
	addMetadataSourceFlag(cmd.Flags())

	// specific flags

	mustSelect(cfg.AliasType, map[AliasType]func(){
		OpencastAlias: func() {
			addACEFlag(cmd.Flags())
			addMetadataIDFlag(cmd.Flags())
			addWorkflowDefinitionFlag(cmd.Flags())
			addWorkflowPropertiesFlag(cmd.Flags())
			addTrackXFlag("presenter", cmd.Flags())
			addTrackXFlag("presentation", cmd.Flags())
			addTrackXFlag("audio", cmd.Flags())
		},
		TalesAlias: func() {
			addACLPresetFlag(cmd.Flags())
			addACLUsersReadFlag(cmd.Flags())
			addACLUsersWriteFlag(cmd.Flags())
			addTrackXFlag("main", cmd.Flags())
			addTrackXFlag("secondary", cmd.Flags())
		},
	})()

	return cmd
}

func eventListCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"list",
		"List Events",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Event
				req svc.EventListRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastEvent(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesEvent(extAPI) },
			})()

			req.FilterByText = getFilterByXStringFlag("text", cmd.Flags())
			req.FilterByID = getFilterByXStringFlag("id", cmd.Flags())
			filterByStatus := getFilterByXMapValueFlag[api.EventStatus]("status", cmd.Flags())
			if filterByStatus != "all" {
				req.FilterByStatus = filterByStatus
			}
			req.FilterByTitle = getFilterByXStringFlag("title", cmd.Flags())
			req.FilterByDescription = getFilterByXStringFlag("description", cmd.Flags())
			req.FilterBySeries = getFilterByXStringFlag("series", cmd.Flags())
			req.FilterBySeriesID = getFilterByXStringFlag("series-id", cmd.Flags())
			req.FilterByCreationDate = getFilterByXStringFlag("creation-date", cmd.Flags())
			req.FilterByStartDate = getFilterByXStringFlag("start-date", cmd.Flags())
			req.FilterByContributors = getFilterByXStringFlag("contributors", cmd.Flags())
			req.FilterByPresenters = getFilterByXStringFlag("presenters", cmd.Flags())
			req.FilterByLocation = getFilterByXStringFlag("location", cmd.Flags())
			req.FilterByLanguage = getFilterByXStringFlag("language", cmd.Flags())
			req.FilterByRightsHolder = getFilterByXStringFlag("rights-holder", cmd.Flags())
			req.FilterByLicense = getFilterByXStringFlag("license", cmd.Flags())
			req.FilterBySubjects = getFilterByXStringFlag("subjects", cmd.Flags())
			req.FilterBySource = getFilterByXStringFlag("source", cmd.Flags())
			req.FilterByScheduledStartDate = getFilterByXStringFlag("scheduled-start-date", cmd.Flags())
			req.FilterByAgentName = getFilterByXStringFlag("agent-name", cmd.Flags())
			req.SortBy = getSortByFlag(cmd.Flags())
			req.SortDirection = getSortDirectionFlag(cmd.Flags())

			return s.List(cmd.Context(), req)
		},
	)
	addFilterByXStringFlag("text", cmd.Flags())
	addFilterByXStringFlag("id", cmd.Flags())
	addFilterByXMapValueFlag("status", listValue(
		string(api.ProcessingEventStatus),
		[]api.EventStatus{
			api.EventStatus("all"),
			api.IngestingEventStatus,
			api.PausedEventStatus,
			api.PendingEventStatus,
			api.ProcessedEventStatus,
			api.ProcessingEventStatus,
			api.ProcessingCancelledEventStatus,
			api.ProcessingFailureEventStatus,
			api.RecordingEventStatus,
			api.RecordingFailureEventStatus,
			api.ScheduledEventStatus,
		},
	), cmd.Flags())
	addFilterByXStringFlag("title", cmd.Flags())
	addFilterByXStringFlag("description", cmd.Flags())
	addFilterByXStringFlag("series", cmd.Flags())
	addFilterByXStringFlag("series-id", cmd.Flags())
	addFilterByXStringFlag("creation-date", cmd.Flags())
	addFilterByXStringFlag("start-date", cmd.Flags())
	addFilterByXStringFlag("contributors", cmd.Flags())
	addFilterByXStringFlag("presenters", cmd.Flags())
	addFilterByXStringFlag("location", cmd.Flags())
	addFilterByXStringFlag("language", cmd.Flags())
	addFilterByXStringFlag("rights-holder", cmd.Flags())
	addFilterByXStringFlag("license", cmd.Flags())
	addFilterByXStringFlag("subjects", cmd.Flags())
	addFilterByXStringFlag("source", cmd.Flags())
	addFilterByXStringFlag("scheduled-start-date", cmd.Flags())
	addFilterByXStringFlag("agent-name", cmd.Flags())
	addSortByFlag(&mapValue[string]{
		Default: "start_date",
		Map: map[string]string{
			"title":                "title",
			"presenters":           "presenter",
			"start_date":           "start_date",
			"end_date":             "end_date",
			"review_status":        "review_status",
			"workflow_status":      "workflow_state",
			"scheduling_status":    "scheduling_status",
			"series":               "series_name",
			"location":             "location",
			"scheduled_start_date": "technical_start",
			"scheduled_end_date":   "technical_end",
		},
	}, cmd.Flags())
	addSortDirectionFlag(cmd.Flags())
	cmd.GroupID = ManagementGroup.ID
	return cmd
}

func eventGetCommand(cfg *Config) *cobra.Command {
	cmd := extAPICommand(
		"get [id]",
		"Get an Event",
		cfg,
		func(cmd *cobra.Command, args []string, extAPI extapiclientv1.Client) (any, error) {
			var (
				s   svc.Event
				req svc.EventGetRequest
			)

			mustSelect(cfg.AliasType, map[AliasType]func(){
				OpencastAlias: func() { s = svc.NewOpencastEvent(extAPI) },
				TalesAlias:    func() { s = svc.NewTalesEvent(extAPI) },
			})()

			req.ID = args[0]

			return s.Get(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.GroupID = ManagementGroup.ID
	return cmd
}
