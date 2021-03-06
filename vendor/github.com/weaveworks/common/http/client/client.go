package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/instrument"
	oldcontext "golang.org/x/net/context"
)

// Requester executes an HTTP request.
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// TimeRequestHistogram performs an HTTP client request and records the duration in a histogram
func TimeRequestHistogram(ctx context.Context, operation string, metric *prometheus.HistogramVec, client Requester, request *http.Request) (*http.Response, error) {
	var response *http.Response
	doRequest := func(_ oldcontext.Context) error {
		var err error
		response, err = client.Do(request)
		return err
	}
	toStatusCode := func(err error) string {
		if err == nil {
			return strconv.Itoa(response.StatusCode)
		}
		return "error"
	}
	err := instrument.TimeRequestHistogramStatus(ctx, fmt.Sprintf("%s %s", request.Method, operation), metric, toStatusCode, doRequest)
	return response, err
}
