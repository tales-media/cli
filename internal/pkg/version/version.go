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

package version

import (
	"fmt"
	"runtime"
)

var (
	version      = "0.0.0"   // set during build
	gitCommit    = "unknown" // set during build
	gitTreeState = "unknown" // set during build
	buildDate    = "unknown" // set during build

	goVersion = runtime.Version()
	compiler  = runtime.Compiler
	platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	CLI = Info{
		Name:         "unknown",
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    goVersion,
		Compiler:     compiler,
		Platform:     platform,
	}
)

// Info holds version and build information. The fields are largely the same as in the `k8s.io/kubernetes/pkg/version`
// package of the Kubernetes project.
type Info struct {
	// Name of the versioned object.
	Name string `human:"Name,wideonly" json:"name" yaml:"name"`

	// Version number.
	Version string `human:"Version" json:"version" yaml:"version"`

	// GitCommit SHA.
	GitCommit string `human:"Git Commit,wideonly" json:"gitCommit" yaml:"gitCommit"`

	// GitTreeState is either "clean" or "dirty".
	GitTreeState string `human:"Git Tree State,wideonly" json:"gitTreeState" yaml:"gitTreeState"`

	// BuildDate of the binary.
	BuildDate string `human:"Build Date,wideonly" json:"buildDate" yaml:"buildDate"`

	// GoVersion of the binary.
	GoVersion string `human:"Go Version,wideonly" json:"goVersion" yaml:"goVersion"`

	// Compiler used for the binary.
	Compiler string `human:"Compiler,wideonly" json:"compiler" yaml:"compiler"`

	// Platform the binary is compiled for.
	Platform string `human:"Platform,wideonly" json:"platform" yaml:"platform"`
}
