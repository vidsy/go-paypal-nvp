package paypalnvp

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type (
	// Response struct for response from NVP request.
	Response struct {
		*http.Response
		ParsedQueryParams *url.Values
		Acknowledgement   string    `nvp_field:"ACK"`
		CorrelationID     string    `nvp_field:"CORRELATIONID"`
		TimeStamp         time.Time `nvp_field:"TIMESTAMP"`
		Version           string    `nvp_field:"VERSION"`
		Build             string    `nvp_field:"BUILD"`
		Errors            []ResponseError
	}

	// ResponseError struct for any errors returned from an NVP request.
	ResponseError struct {
		Code         string `nvp_field:"L_ERRORCODE%d"`
		ShortMessage string `nvp_field:"L_SHORTMESSAGE%d"`
		LongMessage  string `nvp_field:"L_LONGMESSAGE%d"`
		SeverityCode string `nvp_field:"L_SEVERITYCODE%d"`
		ParamID      string `nvp_field:"L_ERRORPARAMID%d"`
		ParamValue   string `nvp_field:"L_ERRORPARAMVALUE%d"`
	}
)

// NewResponse Creates new response from net/http response.
func NewResponse(httpResponse *http.Response) (*Response, error) {
	response := &Response{Response: httpResponse}
	data, err := response.parseBody()
	if err != nil {
		return nil, err
	}
	response.ParsedQueryParams = data
	response.mapFields()

	return response, nil
}

// Successful indicates if the request was valid based on status code and
// NVP response fields.
func (r Response) Successful() bool {
	if r.StatusCode != 200 {
		return false
	}

	if r.ErrorCount() > 0 {
		return false
	}

	return true
}

// ErrorCount count of errors returned in response.
func (r *Response) ErrorCount() int {
	errorCount := 0
	for key, _ := range *r.ParsedQueryParams {
		if strings.Contains(key, "L_ERRORCODE") {
			errorCount += 1
		}
	}
	return errorCount
}

func (r *Response) parseBody() (*url.Values, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}

	if err = r.Body.Close(); err != nil {
		return nil, err
	}

	data, err := url.ParseQuery(string(body[:]))
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Response) mapFields() {
	mapping := r.tagToFieldMapping()

	for key, values := range *r.ParsedQueryParams {
		if fieldName, exists := mapping[key]; exists {
			responseValue := reflect.ValueOf(r)
			field := responseValue.Elem().FieldByName(fieldName)
			switch field.Kind() {
			case reflect.String:
				field.SetString(values[0])
			}
		}
	}

	for i := 0; i < r.ErrorCount(); i++ {
		errorItem := ResponseError{}
		errorItemType := reflect.TypeOf(errorItem)
		errorItemValue := reflect.ValueOf(&errorItem)

		for ii := 0; ii < errorItemType.NumField(); ii++ {
			field := errorItemType.Field(ii)
			if fieldTag, ok := field.Tag.Lookup("nvp_field"); ok {
				fieldValue := errorItemValue.Elem().FieldByName(field.Name)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(r.ParsedQueryParams.Get(fmt.Sprintf(fieldTag, i)))
				}
			}
		}

		r.Errors = append(r.Errors, errorItem)
	}
}

func (r *Response) tagToFieldMapping() map[string]string {
	mapping := make(map[string]string)
	itemType := reflect.TypeOf(*r)

	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)
		if fieldTag, ok := field.Tag.Lookup("nvp_field"); ok {
			mapping[fieldTag] = field.Name
		}
	}

	return mapping
}
