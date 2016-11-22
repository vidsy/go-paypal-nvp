package paypalnvp

import (
	"net/http"
)

const (

	// Base API endpoint for NVP requests
	baseAPIEndpoint = "https://%s.paypal.com/nvp"

	// Sandbox prefix for api signature requests
	sandboxAPISignatureRequestPrefix = "api.sandbox"

	// Live prefix for api signature requests
	apiSignatureRequestPrefix = "api"
)

type (
	// Client struct used to interact with the NVP API.
	Client struct {
		client TransportClient
	}

	// TransportClient interface for client providing HTTP transport
	// functionality.
	TransportClient interface {
		Do(req *http.Request) (resp *http.Response, err error)
	}
)

func NewClient(client TransportClient) *Client {
	if client == nil {
		client = &http.Client{}
	}

	return &Client{client}
}
