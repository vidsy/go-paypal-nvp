package paypalnvp_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/vidsy/go-paypalnvp"
)

type (
	SerializedDataMock struct {
		mockSerialize func() (string, error)
	}
)

func (sdm SerializedDataMock) Serialize() (string, error) {
	if sdm.mockSerialize != nil {
		return sdm.mockSerialize()
	}

	return "SOME=Data", nil
}

type (
	MockClient struct {
		MockDo func(*http.Request) (*http.Response, error)
	}

	MockReadCloser struct {
		io.Reader
	}
)

func (mrc MockReadCloser) Close() error {
	return nil
}

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
			response, _ := client.Execute(payload)

			if response.StatusCode != 200 {
				t.Fatalf("Expected StatusCode to be 200, got: %d", response.StatusCode)
			}
		})

		t.Run("ReturnsErrorOnSerializeError", func(t *testing.T) {
			httpClient := MockClient{}
			client := paypalnvp.NewClient(httpClient, "test")
			payload := SerializedDataMock{
				mockSerialize: func() (string, error) {
					return "", errors.New("Serializer error")
				},
			}
			_, err := client.Execute(payload)

			if err == nil {
				t.Fatalf("Expected an error, got: %v", err)
			}
		})

		t.Run("ReturnsErrorOnClientRequestError", func(t *testing.T) {
			httpClient := MockClient{
				MockDo: func(request *http.Request) (*http.Response, error) {
					return nil, errors.New("Client error")
				},
			}
			client := paypalnvp.NewClient(httpClient, "test")
			payload := SerializedDataMock{}
			_, err := client.Execute(payload)

			if err == nil {
				t.Fatalf("Expected an error, got: %v", err)
			}

		})
	})
}
