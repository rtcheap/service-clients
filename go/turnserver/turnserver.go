package turnserver

import (
	"context"

	"github.com/CzarSimon/httputil/client"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/opentracing/opentracing-go"
)

// Client client interface for interacting with the turn-server session api.
type Client interface {
	Register(ctx context.Context) error
	GetStatistics(ctx context.Context) error
}

// NewClient creates a new client using the default implementatino.
func NewClient(httpClient client.Client) Client {
	if httpClient.Role == "" {
		httpClient.Role = jwt.SystemRole
	}
	if httpClient.UserAgent == "" {
		httpClient.UserAgent = "turnserver/restClient"
	}

	return &restClient{
		http: httpClient,
	}
}

type restClient struct {
	http client.Client
}

func (c *restClient) Register(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "turnserver.restClient.Register")
	defer span.Finish()

	span.LogFields(tracelog.Bool("success", true))
	return nil
}

func (c *restClient) GetStatistics(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "turnserver.restClient.GetStatistics")
	defer span.Finish()

	span.LogFields(tracelog.Bool("success", true))
	return nil
}
