package paypalnvp_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/vidsy/go-paypalnvp"
)

type (
	SerializedDataMock struct{}
)

func (sdm SerializedDataMock) Serialize() (string, error) {
	return "SOME=Data", nil
}

type (
	// MockClient used for mocking interactions with paypal nvp.
	MockClient struct {
		MockDo func(*http.Request) (*http.Response, error)
	}

	// MockReadCloser mock used as Body return for response.
	MockReadCloser struct {
		io.Reader
	}
)

// Close Dummy method for implementing MockReadCloser interface.
func (mrc MockReadCloser) Close() error {
	return nil
}

// Do optionally calls a MockDo then returns a http.Response.
func (mc MockClient) Do(req *http.Request) (*http.Response, error) {
	if mc.MockDo != nil {
		return mc.MockDo(req)
	}

	return NewMockResponse(nil)
}

func NewMockResponse(response []byte) (*http.Response, error) {
	if response == nil {
		response = []byte(`[]`)
	}

	return &http.Response{
		Body:       MockReadCloser{io.MultiReader(bytes.NewReader(response))},
		StatusCode: 200,
	}, nil
}

func TestClient(t *testing.T) {
	t.Run("NewClient", func(t *testing.T) {
		t.Run("CreatesClientWithDefaultHTTPClient", func(t *testing.T) {
			client := paypalnvp.NewClient(nil, "test")

			if client == nil {
				t.Fatalf("Expected new Client, got: %v", client)
			}
		})
	})

	t.Run(".Execute", func(t *testing.T) {
		t.Run("PerformsRequestWithSerializedData", func(t *testing.T) {
			httpClient := MockClient{}
			client := paypalnvp.NewClient(httpClient, "test")
			payload := SerializedDataMock{}
			client.Execute(payload)

			// if response == nil {
			// 	t.Fatalf("Expected response to be nil, got: %v", response)
			// }
		})
	})
}
