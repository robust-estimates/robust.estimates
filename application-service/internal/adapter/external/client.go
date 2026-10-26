package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"robust/application-service/internal/domain"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func (c *Client) Registry(ctx context.Context, inn string) (domain.RegistryData, error) {
	var v domain.RegistryData
	return v, c.get(ctx, "/fns/v1/companies/"+url.PathEscape(inn), &v)
}

func (c *Client) CreditReport(ctx context.Context, inn string) (domain.CreditReport, error) {
	var v domain.CreditReport
	return v, c.get(ctx, "/bki/v1/reports/"+url.PathEscape(inn), &v)
}

func (c *Client) Account(ctx context.Context, inn string) (domain.AccountData, error) {
	var v domain.AccountData
	return v, c.get(ctx, "/bank/v1/accounts/"+url.PathEscape(inn), &v)
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	request, e := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if e != nil {
		return e
	}

	response, e := c.HTTP.Do(request)
	if e != nil {
		var urlError *url.Error
		if errors.As(e, &urlError) {
			return urlError.Err
		}

		return e
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("external source responded with status %d", response.StatusCode)
	}

	return json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(v)
}
