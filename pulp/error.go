package pulp

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"net/url"
)

// Pulp Api docs:
// http://pulp.readthedocs.org/en/latest/dev-guide/conventions/exceptions.html#exception-handling
type ErrorResponse struct {
	Response     *http.Response // HTTP response that caused this error
	ResourceID   string         `json:"resource_id"`
	Message      string         `json:"error_message"` // error message
	ErrorDetails *Error         `json:"error"`         // more detail on individual errors

}

func (r *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(r.Response.Request.URL.Opaque)
	ru := fmt.Sprintf("%s://%s%s", r.Response.Request.URL.Scheme, r.Response.Request.URL.Host, path)
	return fmt.Sprintf("%v %s: %d %v", r.Response.Request.Method, ru, r.Response.StatusCode, r.Message)
}

// Pulp Api docs:
// http://pulp.readthedocs.org/en/latest/dev-guide/conventions/exceptions.html#error-details
type Error struct {
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Data        json.RawMessage `json:"data"`
	Sub_errors  json.RawMessage `json:"sub_errors"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error: %v",
		e.Code, e.Description)
}
