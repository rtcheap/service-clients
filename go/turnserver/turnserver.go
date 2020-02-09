package turnserver

import (
	"context"
	"fmt"

	"github.com/CzarSimon/httputil/client"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/rtcheap/dto"
)

// Client interface for interacting with the turn-server session api.
type Client interface {
	Register(ctx context.Context, baseURL string, session dto.Session) error
	Unregister(ctx context.Context, baseURL, userID string) error
	GetStatistics(ctx context.Context, baseURL string) (dto.SessionStatistics, error)
}

// NewClient creates a new turnserver client using the default implementation.
func NewClient(httpClient client.Client) Client {
	if httpClient.Role == "" {
		httpClient.Role = jwt.SystemRole
	}
	if httpClient.UserAgent == "" {
		httpClient.UserAgent = "turnserver/restClient"
	}
	httpClient.BaseURL = ""

	return &restClient{
		http: httpClient,
	}
}

type restClient struct {
	http client.Client
}

func (c *restClient) Register(ctx context.Context, baseURL string, session dto.Session) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "turnserver.restClient.Register")
	defer span.Finish()

	url := fmt.Sprintf("%s/v1/sessions", baseURL)
	err := c.http.Post(ctx, url, session, nil)
	if err != nil {
		err = fmt.Errorf("failed to register %s. %w", session, err)
		span.LogFields(tracelog.Bool("success", false), tracelog.Error(err))
		return err
	}

	span.LogFields(tracelog.Bool("success", true))
	return nil
}

func (c *restClient) Unregister(ctx context.Context, baseURL, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "turnserver.restClient.Unregister")
	defer span.Finish()

	url := fmt.Sprintf("%s/v1/sessions/%s", baseURL, userID)
	err := c.http.Delete(ctx, url, nil)
	if err != nil {
		err = fmt.Errorf("failed to unregister user(id=%s). %w", userID, err)
		span.LogFields(tracelog.Bool("success", false), tracelog.Error(err))
		return err
	}

	span.LogFields(tracelog.Bool("success", true))
	return nil
}

func (c *restClient) GetStatistics(ctx context.Context, baseURL string) (dto.SessionStatistics, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "turnserver.restClient.GetStatistics")
	defer span.Finish()

	var stats dto.SessionStatistics
	url := fmt.Sprintf("%s/v1/sessions/statistics", baseURL)
	err := c.http.Get(ctx, url, &stats)
	if err != nil {
		err = fmt.Errorf("failed to retrieve turn-server session statistics. %w", err)
		span.LogFields(tracelog.Bool("success", false), tracelog.Error(err))
		return dto.SessionStatistics{}, err
	}

	span.LogFields(tracelog.Bool("success", true))
	return stats, nil
}
