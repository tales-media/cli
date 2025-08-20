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
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/spf13/pflag"
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

func listValue[T any](defaultValue string, list []T) *mapValue[T] {
	v := &mapValue[T]{
		Default: defaultValue,
		Map:     make(map[string]T),
	}
	for _, e := range list {
		v.Map[fmt.Sprintf("%v", e)] = e
	}
	return v
}

type mapValue[T any] struct {
	key     string
	Default string
	Map     map[string]T
}

var _ pflag.Value = &mapValue[any]{}

func (v *mapValue[T]) String() string {
	if v.key == "" {
		return v.Default
	}
	return v.key
}

func (v *mapValue[T]) Value() T {
	if val, ok := v.Map[v.key]; ok {
		return val
	}
	return v.Map[v.Default]
}

func (v *mapValue[T]) Set(s string) error {
	if _, ok := v.Map[strings.ToLower(s)]; ok {
		v.key = s
		return nil
	}
	return fmt.Errorf("unknown option")
}

func (v *mapValue[T]) Type() string {
	return "string"
}

func (v *mapValue[T]) Usage(s string) string {
	keys := slices.Sorted(maps.Keys(v.Map))
	return s + " ( " + strings.Join(keys, " | ") + " )"
}
