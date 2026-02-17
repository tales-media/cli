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

package api

import "encoding/json"

func DeepCopy[T any](in *T) *T {
	copyJSON, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	copy := new(T)
	err = json.Unmarshal(copyJSON, copy)
	if err != nil {
		panic(err)
	}
	return copy
}
