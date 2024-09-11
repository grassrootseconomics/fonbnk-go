package fonbnk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FonbnkClient struct {
	clientId     string
	clientSecret string
	sourceParam  string
	endpoint     string
	httpClient   *http.Client
}

const (
	userAgent   = "fonbnk-go"
	contentType = "application/json"

	baseLiveEndpoint    = "https://aten.fonbnk-services.com"
	baseSandboxEndpoint = "https://sandbox-api.fonbnk.com"
)

// New returns an instance of a Fonbnk client reusbale across different products
func New(clientId string, clientSecret string, sourceParam string, sandbox bool) *FonbnkClient {
	fonbnkClient := &FonbnkClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		sourceParam:  sourceParam,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	if sandbox {
		fonbnkClient.endpoint = baseSandboxEndpoint
	} else {
		fonbnkClient.endpoint = baseLiveEndpoint
	}

	return fonbnkClient
}

// SetHTTPClient can be used to override the default client with a custom set one
func (fc *FonbnkClient) SetHTTPClient(httpClient *http.Client) *FonbnkClient {
	fc.httpClient = httpClient

	return fc
}

// setHeaders sets the headers required by the Fonbnk API
func (fc *FonbnkClient) setHeaders(req *http.Request) (*http.Request, error) {
	if err := fc.setAuthHeaders(req); err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// requestWithCtx builds the HTTP request
func (fc *FonbnkClient) requestWithCtx(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return fc.do(req)
}

// do executes the built http request, setting appropriate headers
func (fc *FonbnkClient) do(req *http.Request) (*http.Response, error) {
	builtRequest, err := fc.setHeaders(req)
	if err != nil {
		return nil, err
	}

	return fc.httpClient.Do(builtRequest)
}

// parseResponse is a general utility to decode JSON responses correctly
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Fonbnk server error: code=%s: response_body=%s", resp.Status, string(b))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
