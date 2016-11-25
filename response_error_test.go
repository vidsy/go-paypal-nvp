package paypalnvp_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/vidsy/go-paypalnvp"
)

func TestResponseError(t *testing.T) {
	t.Run(".Error()", func(t *testing.T) {
		t.Run("CorrectErrorFormat", func(t *testing.T) {
			data := `L_ERRORCODE0=15005&L_SHORTMESSAGE0=Processor%20Decline&L_LONGMESSAGE0=This%20transaction%20cannot%20be%20processed%2e&L_SEVERITYCODE0=Error&L_ERRORPARAMID0=ProcessorResponse&L_ERRORPARAMVALUE0=0051`
			httpResponse := &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
				StatusCode: 200,
			}

			response, _ := paypalnvp.NewResponse(httpResponse)
			responseError := response.Errors[0]

			expectedMessage := `Code '15005', ShortMessage: 'Processor Decline', LongMessage: 'This transaction cannot be processed.', SeverityCode: 'Error', ParamID: 'ProcessorResponse', ParamValue: '0051'`

			if expectedMessage != responseError.Error() {
				t.Fatalf("Expected Error() to be '%s', got: '%s'", expectedMessage, responseError.Error())
			}
		})
	})
}
