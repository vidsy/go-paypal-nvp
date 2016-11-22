package paypalnvp_test

import (
	"testing"

	"github.com/vidsy/go-paypalnvp"
)

func TestClient(t *testing.T) {
	t.Run("NewClient", func(t *testing.T) {
		t.Run("CreatesClientWithDefaultHTTPClient", func(t *testing.T) {
			client := paypalnvp.NewClient(nil)

			if client == nil {
				t.Fatalf("Expected new Client, got: %v", client)
			}
		})
	})
}
