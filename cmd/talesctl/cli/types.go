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
	"fmt"
	"path"
	"strings"
)

const TalesctlPrefix = "talesctl"

type Config struct {
	Alias     string
	AliasType AliasType
	Args      []string
}

func Configure(args []string) *Config {
	alias, args := path.Base(args[0]), args[1:]
	return &Config{
		Alias:     alias,
		AliasType: DetectAliasType(alias),
		Args:      args,
	}
}

type AliasType int

const (
	OpencastAlias AliasType = iota
	TalesAlias
)

func DetectAliasType(cmd string) AliasType {
	if strings.HasPrefix(cmd, TalesctlPrefix) {
		return TalesAlias
	}
	return OpencastAlias
}

func mustSelect[K comparable, V any](a K, m map[K]V) V {
	if v, ok := m[a]; ok {
		return v
	}
	panic(fmt.Sprintf("missing value for %v", a))
}

type Output int

const (
	HumanOutput Output = iota
	WideOutput
	JSONOutput
	YAMLOutput
	NoneOutput
)

func OutputValue() *mapValue[Output] {
	return &mapValue[Output]{
		Default: HumanOutput,
		Map: map[string]Output{
			"human": HumanOutput,
			"wide":  WideOutput,
			"json":  JSONOutput,
			"yaml":  YAMLOutput,
			"none":  NoneOutput,
		},
	}
}
