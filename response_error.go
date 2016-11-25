package paypalnvp

import (
	"fmt"
)

type (
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

// Error Formatted error string based on properties.
func (r ResponseError) Error() string {
	return fmt.Sprintf(
		"Code '%s', ShortMessage: '%s', LongMessage: '%s', SeverityCode: '%s', ParamID: '%s', ParamValue: '%s'",
		r.Code,
		r.ShortMessage,
		r.LongMessage,
		r.SeverityCode,
		r.ParamID,
		r.ParamValue,
	)

}
