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

package main

import (
	"context"
	"os"
	"os/signal"

	"shio.solutions/tales.media/cli/cmd/talesctl/cli"
	"shio.solutions/tales.media/cli/internal/pkg/version"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// handle termination signals
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, shutdownSignals...)
	go func() {
		<-sig
		cancel()
		<-sig
		os.Exit(1)
	}()

	// configure
	cfg := cli.Configure(os.Args)
	version.CLI.Name = cfg.Alias

	// execute
	cmd := cli.New(cfg)
	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
