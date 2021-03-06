package serviceregistry

import (
	"context"
	"fmt"

	"github.com/CzarSimon/httputil/client"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/rtcheap/dto"
)

// Client interface for interacting with the service-registry api.
type Client interface {
	Register(ctx context.Context, svc dto.Service) (dto.Service, error)
	Find(ctx context.Context, id string) (dto.Service, error)
	FindByApplication(ctx context.Context, application string, onlyHealthy bool) ([]dto.Service, error)
	SetStatus(ctx context.Context, id string, status dto.ServiceStatus) error
}

// NewClient creates a new client using the default implementation.
func NewClient(httpClient client.Client) Client {
	if httpClient.Role == "" {
		httpClient.Role = jwt.SystemRole
	}
	if httpClient.UserAgent == "" {
		httpClient.UserAgent = "serviceregistry/restClient"
	}

	return &restClient{
		http: httpClient,
	}
}

type restClient struct {
	http client.Client
}

func (c *restClient) Register(ctx context.Context, svc dto.Service) (dto.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "serviceregistry_rest_client_register")
	defer span.Finish()

	var registeredService dto.Service
	err := c.http.Post(ctx, "/v1/services", svc, &registeredService)
	if err != nil {
		err = fmt.Errorf("failed to register service. %w", err)
		span.LogFields(tracelog.Error(err))
		return dto.Service{}, err
	}

	return registeredService, nil
}

func (c *restClient) Find(ctx context.Context, id string) (dto.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "serviceregistry_rest_client_find")
	defer span.Finish()

	var svc dto.Service
	err := c.http.Get(ctx, "/v1/services/"+id, &svc)
	if err != nil {
		err = fmt.Errorf("failed to find service(id=%s). %w", id, err)
		span.LogFields(tracelog.Error(err))
		return dto.Service{}, err
	}

	return svc, nil
}

func (c *restClient) FindByApplication(ctx context.Context, application string, onlyHealthy bool) ([]dto.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "serviceregistry_rest_client_find_by_application")
	defer span.Finish()

	services := make([]dto.Service, 0)
	path := fmt.Sprintf("/v1/services?application=%s&only-healthy=%t", application, onlyHealthy)
	err := c.http.Get(ctx, path, &services)
	if err != nil {
		err = fmt.Errorf("failed to find service for application = %s. %w", application, err)
		span.LogFields(tracelog.Error(err))
		return nil, err
	}

	return services, nil
}

func (c *restClient) SetStatus(ctx context.Context, id string, status dto.ServiceStatus) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "serviceregistry_rest_client_set_status")
	defer span.Finish()

	path := fmt.Sprintf("/v1/services/%s/status/%s", id, status)
	err := c.http.Put(ctx, path, nil, nil)
	if err != nil {
		err = fmt.Errorf("failed to set status %s for service(id=%s). %w", status, id, err)
		span.LogFields(tracelog.Error(err))
		return err
	}

	return nil
}
