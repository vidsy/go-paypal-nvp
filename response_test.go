package paypalnvp_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/vidsy/go-paypalnvp"
)

func TestResponse(t *testing.T) {
	t.Run("NewResponse", func(t *testing.T) {
		t.Run("QueryDataMappedToFields", func(t *testing.T) {
			data := `TIMESTAMP=2011%2d11%2d15T20%3a27%3a02Z&CORRELATIONID=5be53331d9700&ACK=Success&VERSION=78&BUILD=000000`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if response.Acknowledgement != "Success" {
				t.Fatalf("Expected Acknowledgement == 'Success', got: '%s'", response.Acknowledgement)
			}

			if response.CorrelationID != "5be53331d9700" {
				t.Fatalf("Expected CorrelationID == '5be53331d9700', got: '%s'", response.CorrelationID)
			}

			if response.Version != "78" {
				t.Fatalf("Expected Version == '78', got: '%s'", response.Version)
			}

			if response.Build != "000000" {
				t.Fatalf("Expected Build == '78', got: '%s'", response.Build)
			}

			if !response.Successful() {
				t.Fatalf("Expected response to be successful, got: %t", response.Successful())
			}
		})

		t.Run("ErrorsMappedCorrectly", func(t *testing.T) {
			data := `L_ERRORCODE0=15005&L_SHORTMESSAGE0=Processor%20Decline&L_LONGMESSAGE0=This%20transaction%20cannot%20be%20processed%2e&L_SEVERITYCODE0=Error&L_ERRORPARAMID0=ProcessorResponse&L_ERRORPARAMVALUE0=0051`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)
			responseError := response.Errors[0]

			if response.Successful() {
				t.Fatalf("Expected response to not be successful, got: %t", response.Successful())
			}

			if response.ErrorCount() != 1 {
				t.Fatalf("Expected error count to be 1, got: %d", response.ErrorCount())
			}

			if responseError.Code != "15005" {
				t.Fatalf("Expected Code to be '15005', got '%s'", responseError.Code)
			}

			if responseError.ShortMessage != "Processor Decline" {
				t.Fatalf("Expected ShortMessage to be 'Processor Decline', got '%s'", responseError.ShortMessage)
			}

			if responseError.LongMessage != "This transaction cannot be processed." {
				t.Fatalf("Expected LongMessage to be 'This transaction cannot be processed.', got '%s'", responseError.LongMessage)
			}

			if responseError.SeverityCode != "Error" {
				t.Fatalf("Expected SeverityCode to be 'Error', got '%s'", responseError.SeverityCode)
			}

			if responseError.ParamID != "ProcessorResponse" {
				t.Fatalf("Expected ParamID to be 'ProcessorResponse', got '%s'", responseError.ParamID)
			}

			if responseError.ParamValue != "0051" {
				t.Fatalf("Expected ParamID to be '0051', got '%s'", responseError.ParamValue)
			}

		})
	})
}
