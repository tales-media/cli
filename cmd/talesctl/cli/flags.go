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
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

//
// --ace [role:action]
//

func addACEFlag(flags *pflag.FlagSet) {
	flags.StringSlice("ace", nil, "add access control entry form \"role:action\" (can be specified multiple times or a comma separated list)")
}

func getACEFlag(flags *pflag.FlagSet) ([]api.ACE, error) {
	valList := mustGetValue("ace", flags.GetStringSlice)
	val := make([]api.ACE, 0, len(valList))
	for _, c := range valList {
		role, action, ok := strings.Cut(c, ":")
		if !ok {
			return nil, errors.New("invalid config flag: use role:action syntax")
		}
		val = append(val, api.ACE{
			Role:   role,
			Action: api.Action(action),
			Allow:  true,
		})
	}
	return val, nil
}

//
// --acl-preset [public | organization | private]
//

func addACLPresetFlag(flags *pflag.FlagSet) {
	aclPresetValue := &mapValue[api.TalesACLPreset]{
		Default: string(api.PrivateTalesACLPreset),
		Map: map[string]api.TalesACLPreset{
			string(api.PublicTalesACLPreset):       api.PublicTalesACLPreset,
			string(api.OrganizationTalesACLPreset): api.OrganizationTalesACLPreset,
			string(api.PrivateTalesACLPreset):      api.PrivateTalesACLPreset,
		},
	}
	flags.Var(aclPresetValue, "acl-preset", aclPresetValue.Usage("ACL preset"))
}

func getACLPresetFlag(flags *pflag.FlagSet) api.TalesACLPreset {
	flag := mustGetFlag("acl-preset", flags)
	val, ok := flag.Value.(*mapValue[api.TalesACLPreset])
	if !ok {
		panic("BUG: flag 'acl-preset' has incorrect type")
	}
	return val.Value()
}

//
// --acl-users-read [username]
//

func addACLUsersReadFlag(flags *pflag.FlagSet) {
	flags.StringSlice("acl-users-read", nil, "add additional user with read access (can be specified multiple times or a comma separated list)")
}

func getACLUsersReadFlag(flags *pflag.FlagSet) []string {
	return mustGetValue("acl-users-read", flags.GetStringSlice)
}

//
// --acl-users-write [username]
//

func addACLUsersWriteFlag(flags *pflag.FlagSet) {
	flags.StringSlice("acl-users-write", nil, "add additional user with write access (can be specified multiple times or a comma separated list)")
}

func getACLUsersWriteFlag(flags *pflag.FlagSet) []string {
	return mustGetValue("acl-users-write", flags.GetStringSlice)
}

//
// --auth-open-browser
//

func addAuthOpenBrowserFlag(flags *pflag.FlagSet) {
	flags.Bool("auth-open-browser", true, "automatically open auth code URL in browser")
}

func getAuthOpenBrowserFlag(flags *pflag.FlagSet) bool {
	return mustGetValue("auth-open-browser", flags.GetBool)
}

//
// --auth-out-of-band
//

func addAuthOutOfBandFlag(flags *pflag.FlagSet) {
	flags.Bool("auth-out-of-band", false, "authenticate using OIDC out-of-band")
}

func getAuthOutOfBandFlag(flags *pflag.FlagSet) bool {
	return mustGetValue("auth-out-of-band", flags.GetBool)
}

//
// --context [name]
//

func addContextFlag(flags *pflag.FlagSet) {
	flags.String("context", "", "the name of the Opencast context to use")
}

func getContextFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("context", flags)
	return flag.Value.String()
}

//
// --context-use
//

func addContextUseFlag(flags *pflag.FlagSet) {
	flags.Bool("context-use", false, "make new context the current context")
}

func getContextUseFlag(flags *pflag.FlagSet) bool {
	return mustGetValue("context-use", flags.GetBool)
}

//
// --context-service-mapper [static | dynamic]
//

func addContextServiceMapperFlag(flags *pflag.FlagSet) {
	ContextServiceMapperValue := &mapValue[ContextServiceMapper]{
		Default: "static",
		Map: map[string]ContextServiceMapper{
			"static":  StaticContextServiceMapper,
			"dynamic": DynamicContextServiceMapper,
		},
	}
	flags.Var(ContextServiceMapperValue, "context-service-mapper", ContextServiceMapperValue.Usage("set service mapper type"))
}

func getContextServiceMapperFlag(flags *pflag.FlagSet) ContextServiceMapper {
	flag := mustGetFlag("context-service-mapper", flags)
	val, ok := flag.Value.(*mapValue[ContextServiceMapper])
	if !ok {
		panic("BUG: flag 'context-service-mapper' has incorrect type")
	}
	return val.Value()
}

type ContextServiceMapper int

const (
	StaticContextServiceMapper ContextServiceMapper = iota
	DynamicContextServiceMapper
)

//
// --context-static-service-mapping [host=url]
//

func addContextStaticServiceMappingFlag(flags *pflag.FlagSet) {
	flags.StringSlice("context-static-service-mapping", nil, "set static service mapping form \"host=url\" (can be specified multiple times or a comma separated list)")
}

func getContextStaticServiceMappingFlag(flags *pflag.FlagSet) (map[string]string, error) {
	valList := mustGetValue("context-static-service-mapping", flags.GetStringSlice)
	val := make(map[string]string, len(valList))
	for _, c := range valList {
		k, v, ok := strings.Cut(c, "=")
		if !ok {
			return nil, errors.New("invalid config flag: use host=url syntax")
		}
		val[k] = v
	}
	return val, nil
}

//
// --context-dynamic-service-mapper-ttl [time.Duration]
//

func addContextDynamicServiceMapperTTLFlag(flags *pflag.FlagSet) {
	// TODO: use default from api package
	flags.Duration("context-dynamic-service-mapper-ttl", 10*time.Minute, "set TTL for dynamic service mapper")
}

func getContextDynamicServiceMapperTTLFlag(flags *pflag.FlagSet) time.Duration {
	return mustGetValue("context-dynamic-service-mapper-ttl", flags.GetDuration)
}

//
// --context-basic-auth [username:password]
//

func addContextBasicAuthFlag(flags *pflag.FlagSet) {
	flags.String("context-basic-auth", "", "username and passwort for HTTP Basic Auth in the form \"username:password\"")
}

func getContextBasicAuthFlag(flags *pflag.FlagSet) (username, password string, err error) {
	flag := mustGetFlag("context-basic-auth", flags)
	usernamePassword := flag.Value.String()
	if usernamePassword == "" {
		return
	}
	var ok bool
	if username, password, ok = strings.Cut(usernamePassword, ":"); !ok {
		err = errors.New("incorrect format")
	}
	return
}

//
// --context-jwt-auth [token]
//

func addContextJWTAuthFlag(flags *pflag.FlagSet) {
	flags.String("context-jwt-auth", "", "JWT token for JWT Auth")
}

func getContextJWTAuthFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("context-jwt-auth", flags)
	return flag.Value.String()
}

//
// --context-jwt-auth-header [header]
//

func addContextJWTAuthHeaderFlag(flags *pflag.FlagSet) {
	flags.String("context-jwt-auth-header", "Authorization", "the header name for JWT Auth")
}

func getContextJWTAuthHeaderFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("context-jwt-auth-header", flags)
	return flag.Value.String()
}

//
// --context-jwt-auth-prefix [prefix]
//

func addContextJWTAuthPrefixFlag(flags *pflag.FlagSet) {
	flags.String("context-jwt-auth-prefix", "Bearer ", "the header value prefix for JWT Auth")
}

func getContextJWTAuthPrefixFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("context-jwt-auth-prefix", flags)
	return flag.Value.String()
}

//
// --filter-by-{key} [value]
//

func addFilterByXStringFlag(key string, flags *pflag.FlagSet) {
	flags.String("filter-by-"+key, "", "filter resource list by "+key)
}

func getFilterByXStringFlag(key string, flags *pflag.FlagSet) string {
	flagName := "filter-by-" + key
	return mustGetValue(flagName, flags.GetString)
}

func addFilterByXMapValueFlag[T any](key string, filterValue *mapValue[T], flags *pflag.FlagSet) {
	flags.Var(filterValue, "filter-by-"+key, filterValue.Usage("filter resource list by "+key))
}

func getFilterByXMapValueFlag[T any](key string, flags *pflag.FlagSet) T {
	flagName := "filter-by-" + key
	flag := mustGetFlag(flagName, flags)
	val, ok := flag.Value.(*mapValue[T])
	if !ok {
		panic("BUG: flag '" + flagName + "' has incorrect type")
	}
	return val.Value()
}

//
// TODO: --metadata-contributor [contributor] ***
//

//
// TODO: --metadata-created [time.Date]
//

//
// TODO: --metadata-creator [creator] ***
//

//
// --metadata-description [description]
//

func addMetadataDescriptionFlag(flags *pflag.FlagSet) {
	flags.String("metadata-description", "", "the description in the standard metadata catalog")
}

func getMetadataDescriptionFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-description", flags)
	return flag.Value.String()
}

//
// --metadata-duration [time.Duration]
//

func addMetadataDurationFlag(flags *pflag.FlagSet) {
	flags.Duration("metadata-duration", 0, "the duration in the standard metadata catalog")
}

func getMetadataDurationFlag(flags *pflag.FlagSet) time.Duration {
	return mustGetValue("metadata-duration", flags.GetDuration)
}

//
// --metadata-id [id]
//

func addMetadataIDFlag(flags *pflag.FlagSet) {
	flags.String("metadata-id", "", "the identifier in the standard metadata catalog")
}

func getMetadataIDFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-id", flags)
	return flag.Value.String()
}

//
// TODO: --metadata-language [language]
//

//
// TODO: --metadata-license [license]
//

//
// --metadata-location [location]
//

func addMetadataLocationFlag(flags *pflag.FlagSet) {
	flags.String("metadata-location", "", "the location in the standard metadata catalog")
}

func getMetadataLocationFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-location", flags)
	return flag.Value.String()
}

//
// TODO: --metadata-presenter [presenter] ***
//

//
// TODO: --metadata-publisher [publisher] ????
//

//
// --metadata-rights-holder [rights-holder]
//

func addMetadataRightsHolderFlag(flags *pflag.FlagSet) {
	flags.String("metadata-rights-holder", "", "the rights-holder in the standard metadata catalog")
}

func getMetadataRightsHolderFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-rights-holder", flags)
	return flag.Value.String()
}

//
// --metadata-series-id [series-id]
//

func addMetadataSeriesIDFlag(flags *pflag.FlagSet) {
	flags.String("metadata-series-id", "", "the series-id in the standard metadata catalog")
}

func getMetadataSeriesIDFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-series-id", flags)
	return flag.Value.String()
}

//
// --metadata-source [source]
//

func addMetadataSourceFlag(flags *pflag.FlagSet) {
	flags.String("metadata-source", "", "the source in the standard metadata catalog")
}

func getMetadataSourceFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-source", flags)
	return flag.Value.String()
}

//
// --metadata-start-date [time.Time]
//

func addMetadataStartDateFlag(flags *pflag.FlagSet) {
	flags.Time("metadata-start-date", time.Time{}, []string{time.RFC3339}, "the start-date in the standard metadata catalog in RFC3339 format (e.g. '"+time.RFC3339+"')")
}

func getMetadataStartDateFlag(flags *pflag.FlagSet) time.Time {
	return mustGetValue("metadata-start-date", flags.GetTime)
}

//
// --metadata-subject [subject]
//

func addMetadataSubjectFlag(flags *pflag.FlagSet) {
	flags.StringSlice("metadata-subject", nil, "the subject in the standard metadata catalog (can be specified multiple times or a comma separated list)")
}

func getMetadataSubjectFlag(flags *pflag.FlagSet) []string {
	return mustGetValue("metadata-subject", flags.GetStringSlice)
}

//
// --metadata-title [title]
//

func addMetadataTitleFlag(flags *pflag.FlagSet) {
	flags.String("metadata-title", "", "the title in the standard metadata catalog")
}

func getMetadataTitleFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("metadata-title", flags)
	return flag.Value.String()
}

//
// -o, --output [Output]
//

func addOutputFlag(flags *pflag.FlagSet) {
	outputValue := &mapValue[Output]{
		Default: "human",
		Map: map[string]Output{
			"human": HumanOutput,
			"wide":  WideOutput,
			"json":  JSONOutput,
			"yaml":  YAMLOutput,
			"none":  NoneOutput,
		},
	}
	flags.VarP(outputValue, "output", "o", outputValue.Usage("the output format"))
}

func getOutputFlag(flags *pflag.FlagSet) Output {
	flag := mustGetFlag("output", flags)
	val, ok := flag.Value.(*mapValue[Output])
	if !ok {
		panic("BUG: flag 'output' has incorrect type")
	}
	return val.Value()
}

type Output int

const (
	HumanOutput Output = iota
	WideOutput
	JSONOutput
	YAMLOutput
	NoneOutput
)

//
// --sort-by [key]
//

func addSortByFlag(sortValue *mapValue[string], flags *pflag.FlagSet) {
	flags.Var(sortValue, "sort-by", sortValue.Usage("sort resource list by key"))
}

func getSortByFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("sort-by", flags)
	val, ok := flag.Value.(*mapValue[string])
	if !ok {
		panic("BUG: flag 'sort-by' has incorrect type")
	}
	return val.Value()
}

//
// --sort-direction [ascending | asc | descending | desc]
//

func addSortDirectionFlag(flags *pflag.FlagSet) {
	sortDirectionValue := &mapValue[api.SortDirection]{
		Default: "ascending",
		Map: map[string]api.SortDirection{
			"ascending":  api.Ascending,
			"asc":        api.Ascending,
			"descending": api.Descending,
			"desc":       api.Descending,
		},
	}
	flags.Var(sortDirectionValue, "sort-direction", sortDirectionValue.Usage("sort resource list in this direction"))
}

func getSortDirectionFlag(flags *pflag.FlagSet) api.SortDirection {
	flag := mustGetFlag("sort-direction", flags)
	val, ok := flag.Value.(*mapValue[api.SortDirection])
	if !ok {
		panic("BUG: flag 'sort-direction' has incorrect type")
	}
	return val.Value()
}

//
// --track{-name} [file | filename:-]
//

func addTrackXFlag(name string, flags *pflag.FlagSet) {
	flagName := "track"
	flagDescription := "upload file as track"
	if name != "" {
		flagName += "-" + name
		flagDescription += " (" + name + ")"
	}
	flags.String(flagName, "", flagDescription)
}

func getTrackXFlag(name string, cmd *cobra.Command) (filename string, r io.ReadCloser, err error) {
	flagName := "track"
	if name != "" {
		flagName += "-" + name
	}
	val := mustGetValue(flagName, cmd.Flags().GetString)

	if val == "" {
		return
	}

	file, isStdin := strings.CutSuffix(val, ":-")
	if isStdin {
		filename = file
		r = io.NopCloser(cmd.InOrStdin())
		return
	}

	filename = filepath.Base(file)
	r, err = os.Open(file)
	return
}

//
// --workflow-definition [name]
//

func addWorkflowDefinitionFlag(flags *pflag.FlagSet) {
	flags.String("workflow-definition", "schedule-and-upload", "the name of the workflow definition")
}

func getWorkflowDefinitionFlag(flags *pflag.FlagSet) string {
	flag := mustGetFlag("workflow-definition", flags)
	return flag.Value.String()
}

//
// --workflow-property [key=value]
//

func addWorkflowPropertiesFlag(flags *pflag.FlagSet) {
	flags.StringSliceP("workflow-property", "p", nil, "set workflow configuration property in the form \"key=value\" (can be specified multiple times or a comma separated list)")
}

func getWorkflowPropertiesFlag(flags *pflag.FlagSet) (map[string]string, error) {
	valList := mustGetValue("workflow-property", flags.GetStringSlice)
	val := make(map[string]string, len(valList))
	for _, c := range valList {
		k, v, ok := strings.Cut(c, "=")
		if !ok {
			return nil, errors.New("invalid config flag: use key=value syntax")
		}
		val[k] = v
	}
	return val, nil
}

//
// --workflow-status [WorkflowStatus]
//

func addWorkflowStatusFlag(flags *pflag.FlagSet) {
	workflowStatusValue := &mapValue[WorkflowStatus]{
		Default: "none",
		Map: map[string]WorkflowStatus{
			"none":    NoneWorkflowStatus,
			"running": RunningWorkflowStatus,
			"paused":  PausedWorkflowStatus,
			"stopped": StoppedWorkflowStatus,
		},
	}
	flags.Var(workflowStatusValue, "workflow-status", workflowStatusValue.Usage("set a new workflow status"))
}

func getWorkflowStatusFlag(flags *pflag.FlagSet) WorkflowStatus {
	flag := mustGetFlag("workflow-status", flags)
	val, ok := flag.Value.(*mapValue[WorkflowStatus])
	if !ok {
		panic("BUG: flag 'workflow-status' has incorrect type")
	}
	return val.Value()
}

type WorkflowStatus int

const (
	NoneWorkflowStatus WorkflowStatus = iota
	RunningWorkflowStatus
	PausedWorkflowStatus
	StoppedWorkflowStatus
)

// TODO: workflow-normalization-skip-remux
// TODO: workflow-normalization-skip-audio
// TODO: workflow-normalization-skip-audio-sample-rate
// TODO: workflow-normalization-skip-audio-channels
// TODO: workflow-normalization-skip-audio-loudnorm
// TODO: workflow-normalization-skip-video
// TODO: workflow-normalization-skip-video-resolution
// TODO: workflow-normalization-skip-video-frame-rate
// TODO: workflow-normalization-skip-video-colorspace
// TODO: workflow-normalization-skip-video-deinterlace
// TODO: workflow-normalization-skip-video-rotation
// TODO: workflow-edit-required

func mustGetFlag(name string, flags *pflag.FlagSet) *pflag.Flag {
	flag := flags.Lookup(name)
	if flag == nil {
		panicBugUndefinedFlag(name)
	}
	return flag
}

func mustGetValue[T any](name string, f func(string) (T, error)) T {
	val, err := f(name)
	if err != nil {
		panicBugUndefinedFlag(name)
	}
	return val
}

func panicBugUndefinedFlag(name string) {
	panic(fmt.Sprintf("BUG: flag '%s' undefined for this command", name))
}
