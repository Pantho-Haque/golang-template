package providers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"magic.pathao.com/parcel/prism/internal/monitoring/beatbox"
)

type HttpProvider interface {
	makeRequest(method string, fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string, enforceOK bool) (*http.Response, error)
	makeGetRequest(fullUrl string, headers map[string]any, maxRetry int, timeout int, caller string) (*http.Response, error)
	makePostRequest(fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string) (*http.Response, error)
	makePatchRequest(fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string) (*http.Response, error)
	makeDeleteRequest(fullUrl string, headers map[string]any, maxRetry int, timeout int, caller string) (*http.Response, error)
}
type httpProvider struct {
	log     *zap.Logger
	beatbox *beatbox.BeatBox
}

func NewHttpProvider(logger *zap.Logger, beatbox *beatbox.BeatBox) HttpProvider {
	return &httpProvider{
		log: logger,
	}
}

// makeRequest centralizes common HTTP call behavior (headers, timeout, retries, metrics)
// enforceOK controls whether a non-200 response should be treated as an error (true for GET)
func (s *httpProvider) makeRequest(method string, fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string, enforceOK bool) (*http.Response, error) {
	var resp *http.Response

	path := fullUrl
	parsedURL, parseErr := url.Parse(fullUrl)
	if parseErr != nil {
		s.log.Error("failed to parse url", zap.Error(parseErr))
	} else {
		path = parsedURL.Path
	}

	for attempt := 1; attempt <= maxRetry; attempt++ {
		then := time.Now()
		client := &http.Client{}
		client.Timeout = time.Duration(timeout) * time.Millisecond

		req, reqErr := http.NewRequest(method, fullUrl, nil)
		if method != "GET" && method != "DELETE" {
			var bodyReader *bytes.Reader
			if body != nil {
				bodyReader = bytes.NewReader(body)
			}
			req, reqErr = http.NewRequest(method, fullUrl, bodyReader)
		}

		if reqErr != nil {
			s.log.Error("failed to create request", zap.Error(reqErr))
			return nil, reqErr
		}

		for key, value := range headers {
			req.Header.Add(key, fmt.Sprintf("%v", value))
		}

		resp, parseErr = client.Do(req)
		s.beatbox.TimeTakenInServiceCall(caller, path, parseErr == nil, then)
		if resp != nil {
			s.beatbox.StatusCodeCountInServiceCall(caller, path, resp.StatusCode)
		}

		if enforceOK && resp != nil && resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("status code: %d", resp.StatusCode)
			s.log.Error(fmt.Sprintf("failed to get response from %s service", caller),
				zap.Error(err),
				zap.String("url", fullUrl),
				zap.Int("attempt", attempt),
			)
			return nil, err
		}

		if parseErr == nil {
			return resp, nil
		}

		s.log.Error(fmt.Sprintf("failed to get response from %s service", caller),
			zap.Error(parseErr),
			zap.String("url", fullUrl),
			zap.Int("attempt", attempt),
		)
	}

	return nil, fmt.Errorf("failed to get response from %s service after %d attempts", caller, maxRetry)
}

func (s *httpProvider) makeGetRequest(fullUrl string, headers map[string]any, maxRetry int, timeout int, caller string) (*http.Response, error) {
	return s.makeRequest("GET", fullUrl, headers, nil, maxRetry, timeout, caller, true)
}

func (s *httpProvider) makePostRequest(fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string) (*http.Response, error) {
	return s.makeRequest("POST", fullUrl, headers, body, maxRetry, timeout, caller, false)
}

func (s *httpProvider) makePatchRequest(fullUrl string, headers map[string]any, body []byte, maxRetry int, timeout int, caller string) (*http.Response, error) {
	return s.makeRequest("PATCH", fullUrl, headers, body, maxRetry, timeout, caller, false)
}

func (s *httpProvider) makeDeleteRequest(fullUrl string, headers map[string]any, maxRetry int, timeout int, caller string) (*http.Response, error) {
	return s.makeRequest("DELETE", fullUrl, headers, nil, maxRetry, timeout, caller, true)
}
