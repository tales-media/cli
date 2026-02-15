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
	"strings"

	"github.com/spf13/pflag"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

//
// --context [name]
//

func addContextFlag(flags *pflag.FlagSet) {
	flags.String("context", "", "the name of the Opencast context to use")
}

//
// --filter-by-{key} [value]
//

func addFilterByXStringFlag(key string, flags *pflag.FlagSet) {
	flags.String("filter-by-"+key, "", "filter resource list by "+key)
}

func getFilterByXStringFlag(key string, flags *pflag.FlagSet) string {
	flagName := "filter-by-" + key
	val, err := flags.GetString(flagName)
	if err != nil {
		panicBugUndefinedFlag(flagName)
	}
	return val
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
// --workflow-property [key=value]
//

func addWorkflowPropertiesFlag(flags *pflag.FlagSet) {
	flags.StringSliceP("workflow-property", "p", nil, "set workflow configuration property in the form \"key=value\" (can be specified multiple times or a comma separated list)")
}

func getWorkflowPropertiesFlag(flags *pflag.FlagSet) (map[string]string, error) {
	valList, err := flags.GetStringSlice("workflow-property")
	if err != nil {
		panicBugUndefinedFlag("workflow-property")
	}
	val := make(map[string]string)
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
	workflowStatusValue := WorkflowStatusValue()
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
	NoneWorkflowState WorkflowStatus = iota
	RunningWorkflowStatus
	PausedWorkflowStatus
	StoppedWorkflowStatus
)

func WorkflowStatusValue() *mapValue[WorkflowStatus] {
	return &mapValue[WorkflowStatus]{
		Default: "none",
		Map: map[string]WorkflowStatus{
			"none":    NoneWorkflowState,
			"running": RunningWorkflowStatus,
			"paused":  PausedWorkflowStatus,
			"stopped": StoppedWorkflowStatus,
		},
	}
}

func mustGetFlag(name string, flags *pflag.FlagSet) *pflag.Flag {
	flag := flags.Lookup(name)
	if flag == nil {
		panicBugUndefinedFlag(name)
	}
	return flag
}

func panicBugUndefinedFlag(name string) {
	panic(fmt.Sprintf("BUG: flag '%s' undefined for this command", name))
}
