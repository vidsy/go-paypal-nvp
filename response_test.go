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
		})

		t.Run("ErrorsMappedCorrectly", func(t *testing.T) {
			data := `L_ERRORCODE0=15005&L_SHORTMESSAGE0=Processor%20Decline&L_LONGMESSAGE0=This%20transaction%20cannot%20be%20processed%2e&L_SEVERITYCODE0=Error&L_ERRORPARAMID0=ProcessorResponse&L_ERRORPARAMVALUE0=0051`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)
			responseError := response.Errors[0]

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

	t.Run(".Successful", func(t *testing.T) {
		t.Run("WithNoErrorsAndValidStatusCode", func(t *testing.T) {
			data := `TIMESTAMP=2011%2d11%2d15T20%3a27%3a02Z&CORRELATIONID=5be53331d9700&ACK=Success&VERSION=78&BUILD=000000`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if !response.Successful() {
				t.Fatalf("Expected Successful() to be true, got: %t", response.Successful())
			}

		})

		t.Run("NotWithInvalidStatusCode", func(t *testing.T) {
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
				StatusCode: 500,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if response.Successful() {
				t.Fatalf("Expected Successful() to be false, got: %t", response.Successful())
			}
		})

		t.Run("NotWithErrors", func(t *testing.T) {
			data := `L_ERRORCODE0=15005&L_SHORTMESSAGE0=Processor%20Decline&L_LONGMESSAGE0=This%20transaction%20cannot%20be%20processed%2e&L_SEVERITYCODE0=Error&L_ERRORPARAMID0=ProcessorResponse&L_ERRORPARAMVALUE0=0051`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if response.Successful() {
				t.Fatalf("Expected Successful() to be false, got: %t", response.Successful())
			}
		})
	})

	t.Run(".ErrorCount", func(t *testing.T) {
		t.Run("WithErrors", func(t *testing.T) {
			data := `L_ERRORCODE0=15005&L_SHORTMESSAGE0=a&L_LONGMESSAGE0=b&L_SEVERITYCODE0=Error&L_ERRORPARAMID0=ProcessorResponse&L_ERRORPARAMVALUE0=0051&L_ERRORCODE1=15006&L_SHORTMESSAGE1=c&L_LONGMESSAGE0=d&L_SEVERITYCODE1=Error&L_ERRORPARAMID1=ProcessorResponse&L_ERRORPARAMVALUE1=0052`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if response.ErrorCount() != 2 {
				t.Fatalf("Expected ErrorCount() to be 2, got: %d", response.ErrorCount())
			}
		})

		t.Run("WithNoErrors", func(t *testing.T) {
			data := `TIMESTAMP=2011%2d11%2d15T20%3a27%3a02Z&CORRELATIONID=5be53331d9700&ACK=Success&VERSION=78&BUILD=000000`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)

			if response.ErrorCount() != 0 {
				t.Fatalf("Expected ErrorCount() to be 0, got: %d", response.ErrorCount())
			}
		})

	})
}
