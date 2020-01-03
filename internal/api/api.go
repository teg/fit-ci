// Package api implements the endpoints accessed by GitHub
package api

import (
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"goji.io/pat"
)

// An API object wraps the http server
type API struct {
	server *baseapp.Server
}

// New sets up the http server and configures it to serve the github webhook handers we suport
func New(serverConfig baseapp.HTTPConfig, githubConfig githubapp.Config, logger zerolog.Logger) *API {
	server, err := baseapp.NewServer(
		serverConfig,
		baseapp.DefaultParams(logger, "fit-ci.")...,
	)
	if err != nil {
		panic(err)
	}

	cc, err := githubapp.NewDefaultCachingClientCreator(
		githubConfig,
		githubapp.WithClientUserAgent("fit-ci/1.0.0"),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(server.Registry()),
		),
	)
	if err != nil {
		panic(err)
	}

	checkSuiteHandler := &checkSuiteHandler{
		ClientCreator: cc,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(githubConfig, checkSuiteHandler)
	server.Mux().Handle(pat.Post(githubapp.DefaultWebhookRoute), webhookHandler)

	return &API{
		server: server,
	}
}

// Start the server and block, serving the webhook API
func (api *API) Start() error {
	return api.server.Start()
}
