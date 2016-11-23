package paypalnvp

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/vidsy/go-paypalnvp/payload"
)

const (
	baseAPIEndpoint                  = "https://%s.paypal.com/nvp"
	sandboxAPISignatureRequestPrefix = "api.sandbox"
	apiSignatureRequestPrefix        = "api"

	//APIVersion version of the API to use.
	APIVersion = "2.3"

	// Sandsbox environment
	Sandbox = "sandbox"

	// Live environment
	Live = "live"
)

type (
	// Client struct used to interact with the NVP API.
	Client struct {
		client      TransportClient
		environment string
		User        string
		Password    string
		Signature   string
	}

	// TransportClient interface for client providing HTTP transport
	// functionality.
	TransportClient interface {
		Do(req *http.Request) (resp *http.Response, err error)
	}
)

// NewClient Creates a new client.
func NewClient(client TransportClient, environment string, user string, password string, signature string) *Client {
	if client == nil {
		client = &http.Client{}
	}

	return &Client{client, environment, user, password, signature}
}

// Execute performs the NVP request and returns the results.
func (c Client) Execute(item payload.Serializer) (*Response, error) {
	item.SetCredentials(
		c.User,
		c.Password,
		c.Signature,
		APIVersion,
	)

	data, err := item.Serialize()
	if err != nil {
		return nil, err
	}

	response, err := c.perform(data)
	if err != nil {
		return nil, err
	}

	return &Response{Response: response}, nil
}

func (c Client) perform(serializedData string) (*http.Response, error) {
	request, _ := http.NewRequest(
		"POST",
		c.generateEndpoint(),
		bytes.NewBuffer([]byte(serializedData)),
	)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c Client) generateEndpoint() string {
	endpointPrefix := sandboxAPISignatureRequestPrefix
	if c.environment == Live {
		endpointPrefix = apiSignatureRequestPrefix
	}

	return fmt.Sprintf(baseAPIEndpoint, endpointPrefix)
}
