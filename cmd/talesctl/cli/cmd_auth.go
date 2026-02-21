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
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"sync/atomic"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"shio.solutions/tales.media/cli/internal/talesctl/svc"
	"shio.solutions/tales.media/cli/internal/talesctl/svc/api"
)

func authCommand(cfg *Config) *cobra.Command {
	cmd := cfgCommand(
		"auth [org name]",
		"Authenticate with tales.media",
		cfg,
		func(cmd *cobra.Command, args []string) (any, error) {
			s := svc.NewTalesConfig()
			req := svc.ConfigUpdateRequest{}

			// authenticate
			jwt, err := authTalesMedia(cmd)
			if err != nil {
				return nil, err
			}

			// construct new context
			org := args[0]
			ctx := api.Context{
				Name: org + ".tales.media",
				ServiceMapper: api.ServiceMapper{
					Static: &api.StaticServiceMapper{
						Default: fmt.Sprintf("https://api.%s.tales.services", org),
					},
				},
				Authentication: api.Authentication{
					JWT: &api.JWTAuthentication{
						Token:  jwt,
						Header: "X-Forwarded-ID-Token",
						Prefix: "",
					},
				},
			}

			req.Config = cfg.Config
			req.Config.CurrentContext = ctx.Name

			// check if config already contains context entry
			i := slices.IndexFunc(cfg.Contexts, func(c api.Context) bool { return c.Name == ctx.Name })
			if i < 0 {
				req.Config.Contexts = append(req.Config.Contexts, ctx)
			} else {
				req.Config.Contexts[i] = ctx
			}

			return nil, s.Update(cmd.Context(), req)
		},
	)
	cmd.Args = cobra.ExactArgs(1)
	addAuthOpenBrowserFlag(cmd.Flags())
	addAuthOutOfBandFlag(cmd.Flags())
	cmd.GroupID = AdminGroup.ID
	return cmd
}

func authTalesMedia(cmd *cobra.Command) (string, error) {
	ctx := cmd.Context()

	// for user interactions
	scanner := bufio.NewScanner(cmd.InOrStdin())

	// setup OIDC
	var (
		providerURL    = "https://auth.shio.solutions"
		clientID       = "talesctl-zgjm67i8xuxumuoq29mcr0rbrss0e27g"
		rawExtraClaims = `{"id_token":{"name":{"essential":true},"email":{"essential":true},"federated_claims":{"essential":true},"groups":{"essential":true}}}`
		scopes         = []string{
			oidc.ScopeOpenID,
			"email",
			"profile",
			"groups",
			"federated:id",
		}
	)

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return "", err
	}

	oauth2Config := oauth2.Config{
		ClientID: clientID,
		Endpoint: provider.Endpoint(),
		Scopes:   scopes,
	}

	state := rand.Text()
	pkceVerifier := oauth2.GenerateVerifier()

	// determine return channel used to pass the auth code:
	//   OOB  -> user has to copy-paste auth code
	//   HTTP -> start HTTP server and configure callback URL
	var getAuthCode func() (string, error)

	if getAuthOutOfBandFlag(cmd.Flags()) {
		oauth2Config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
		getAuthCode = func() (string, error) {
			cmd.Println()
			cmd.Print("Paste auth code from website: ")
			if !scanner.Scan() {
				return "", errors.New("cli: could not read auth code")
			}
			return scanner.Text(), nil
		}
	} else {
		// create listener on random local port
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return "", err
		}
		defer l.Close()
		oauth2Config.RedirectURL = fmt.Sprintf("http://%s/callback", l.Addr().String())

		// create HTTP handler
		var (
			authCode      string
			callbackCount = &atomic.Int32{}
			callbackDone  = make(chan struct{})
		)
		http.HandleFunc("GET /callback", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/plain")

			// check if /callback was accessed more than once
			if !callbackCount.CompareAndSwap(0, 1) {
				c := callbackCount.Add(1)
				fmt.Fprintf(w, "ERROR: callback can only be called once; called %d times", c)
				r.Response.StatusCode = http.StatusBadRequest
				return
			}

			// closing this channel will teardown the HTTP server
			defer close(callbackDone)
			q := r.URL.Query()

			// check for error
			if err := q.Get("error"); err != "" {
				fmt.Fprint(w, "ERROR: ")
				fmt.Fprint(w, q.Get("error_description"))
				return
			}

			// retrieve auth code
			authCode = q.Get("code")
			fmt.Fprint(w, "DONE: Check terminal output")
		})

		getAuthCode = func() (string, error) {
			var (
				httpErr error
				httpSrv = &http.Server{}
			)

			// start HTTP server
			go func() { httpErr = httpSrv.Serve(l) }()

			// wait for request to /callback
			select {
			case <-callbackDone:
			case <-ctx.Done():
				cmd.PrintErrln("aborting")
			case <-time.After(10 * time.Minute):
				cmd.PrintErrln("timeout after 10m")
			}

			// terminate HTTP server
			_ = httpSrv.Shutdown(ctx)
			if !errors.Is(httpErr, http.ErrServerClosed) {
				return "", httpErr
			}

			return authCode, nil
		}
	}

	// direct user to auth code URL
	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(pkceVerifier),
		oauth2.SetAuthURLParam("claims", rawExtraClaims),
	)

	if getAuthOpenBrowserFlag(cmd.Flags()) {
		cmd.Println("Press enter to open the following auth code URL in our browser:")
		cmd.Println()
		cmd.Print(authCodeURL)
		scanner.Scan()

		if err = browser.OpenURL(authCodeURL); err != nil {
			cmd.PrintErrf("failed to open URL in browser: %s\n", err.Error())
			cmd.PrintErrln("open URL manually:")
			cmd.PrintErrln(authCodeURL)
		}
	} else {
		cmd.Println("Open the following auth code URL in our browser:")
		cmd.Println(authCodeURL)
	}

	// retrieve auth code
	authCode, err := getAuthCode()
	if err != nil {
		return "", err
	}
	if authCode == "" {
		return "", errors.New("cli: received empty auth code")
	}

	// retrieve token
	token, err := oauth2Config.Exchange(
		ctx,
		authCode,
		oauth2.VerifierOption(pkceVerifier),
	)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("cli: received no ID token")
	}

	// validate token
	tokenVerifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})
	_, err = tokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	cmd.Println()
	cmd.Println("Authenticated successfully")
	return rawIDToken, nil
}
