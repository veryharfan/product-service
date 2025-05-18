package pkg

import (
	"context"
	"net/http"
)

type AuthInternalHeader string

const (
	AuthInternalHeaderKey AuthInternalHeader = "X-Internal-Auth"
)

func AddRequestHeader(ctx context.Context, internalAuthHeader string, httpRequest *http.Request) {
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Accept", "application/json")

	if reqID := ctx.Value("request_id"); reqID != nil {
		httpRequest.Header.Add("X-Request-ID", reqID.(string))
	}

	httpRequest.Header.Add(string(AuthInternalHeaderKey), internalAuthHeader)
}
