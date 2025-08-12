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

package formatter

import (
	"io"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

type YAML struct{}

var _ Formatter = &YAML{}

func (f *YAML) List(w io.Writer, list any) error {
	gl, err := NewGenericList(list)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	return enc.Encode(gl)
}

func (f *YAML) Object(w io.Writer, obj any) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(obj)
}

func (f *YAML) Error(w io.Writer, err error) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(&errorObj{Error: err.Error()})
}
