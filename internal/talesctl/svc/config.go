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

package svc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"

	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

const configFile = "config.yaml"

type Config interface {
	Get(context.Context, ConfigGetRequest) (api.Config, error)
	Update(context.Context, ConfigUpdateRequest) error
}

type ConfigGetRequest struct{}

type ConfigUpdateRequest struct {
	Config api.Config
}

type opencastConfig struct {
	ConfigDirName string
}

var _ Config = &opencastConfig{}

func NewOpencastConfig() Config {
	return &opencastConfig{
		ConfigDirName: "opencastctl",
	}
}

func (svc *opencastConfig) Get(ctx context.Context, req ConfigGetRequest) (api.Config, error) {
	cfg := api.Config{}

	f, err := svc.openConfigFile()
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	dec, err := yaml.NewLoader(f)
	if err != nil {
		return cfg, err
	}
	err = dec.Load(&cfg)
	if !errors.Is(err, io.EOF) {
		return cfg, err
	}
	return cfg, nil
}

func (svc *opencastConfig) Update(ctx context.Context, req ConfigUpdateRequest) error {
	f, err := svc.openConfigFile()
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	return enc.Encode(req.Config)
}

func (svc *opencastConfig) openConfigFile() (*os.File, error) {
	// create config root if necessary
	if err := os.MkdirAll(configRoot, 0o755); err != nil {
		return nil, err
	}

	// create config dir if necessary
	configDir := filepath.Join(configRoot, svc.ConfigDirName)
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return nil, err
	}

	// check config dir
	configDirStats, err := os.Stat(configDir)
	if err != nil {
		return nil, err
	}
	if !configDirStats.IsDir() {
		return nil, fmt.Errorf("config: %s is not a directory", configDir)
	}
	if configDirStats.Mode().Perm() != 0o700 {
		return nil, fmt.Errorf("config: %s folder does not have 0700 permissions", configDir)
	}
	// TODO: add owner / group checks

	// open config file
	configFile := filepath.Join(configDir, configFile)
	return os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0o666)
}

type talesConfig = opencastConfig

var _ Config = &talesConfig{}

func NewTalesConfig() Config {
	return &talesConfig{
		ConfigDirName: "talesctl",
	}
}
