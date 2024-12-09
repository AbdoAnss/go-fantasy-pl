package api

import "net/http"

type Client interface {
	Get(endpoint string) (*http.Response, error)
}
