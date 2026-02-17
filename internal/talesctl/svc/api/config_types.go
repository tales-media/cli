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

import "time"

// TODO: add validation and defaulting

const RedactedValue = "REDACTED"

type Meta struct {
	APIVersion string `human:"API Version,wideonly" json:"apiVersion" yaml:"apiVersion"`
	Kind       string `human:"Kind,wideonly" json:"kind" yaml:"kind"`
}

type Config struct {
	Meta `human:"Meta,wideonly" json:",inline" yaml:",inline"`

	CurrentContext string    `human:"Current Context" json:"currentContext" yaml:"currentContext"`
	Contexts       []Context `human:"Contexts,wideonly" json:"contexts" yaml:"contexts"`
}

func (cfg *Config) RedactCredentials() {
	for i := range cfg.Contexts {
		cfg.Contexts[i].RedactCredentials()
	}
}

type Context struct {
	Name           string         `human:"Name" json:"name" yaml:"name"`
	ServiceMapper  ServiceMapper  `human:"Service Mapper,wideonly" json:"serviceMapper" yaml:"serviceMapper"`
	Authentication Authentication `human:"Authentication,wideonly" json:"authentication" yaml:"authentication"`
}

func (ctx *Context) RedactCredentials() {
	if ctx.Authentication.Basic != nil {
		ctx.Authentication.Basic.Password = RedactedValue
	}
	if ctx.Authentication.JWT != nil {
		ctx.Authentication.JWT.Token = RedactedValue
	}
}

type ServiceMapper struct {
	Static  *StaticServiceMapper  `human:"Static" json:"static,omitempty" yaml:"static,omitempty"`
	Dynamic *DynamicServiceMapper `human:"Dynamic" json:"dynamic,omitempty" yaml:"dynamic,omitempty"`
}

type StaticServiceMapper struct {
	Default      string            `human:"Default" json:"default" yaml:"default"`
	ServiceHosts map[string]string `human:"Service Hosts,wideonly" json:"serviceHosts" yaml:"serviceHosts"`
}

type DynamicServiceMapper struct {
	ServiceRegistry string        `human:"Service Registry" json:"serviceRegistry" yaml:"serviceRegistry"`
	TTL             time.Duration `human:"TTL" json:"ttl" yaml:"ttl"`
}

type Authentication struct {
	Basic *BasicAuthentication `human:"Basic" json:"basic,omitempty" yaml:"basic,omitempty"`
	JWT   *JWTAuthentication   `human:"JWT" json:"jwt,omitempty" yaml:"jwt,omitempty"`
}

type BasicAuthentication struct {
	Username string `human:"Username,wideonly" json:"username" yaml:"username"`
	Password string `human:"Password,wideonly" json:"password" yaml:"password"`
}

type JWTAuthentication struct {
	Token string `human:"Token,wideonly" json:"token" yaml:"token"`
}
