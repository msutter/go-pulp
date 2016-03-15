//
// Copyright 2016, Marc Sutter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pulp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	// comment the std libs as they do not support InsecureSkipVerify
	// https://github.com/golang/go/issues/5742
	// replace them with a fixed version
	// "net/http"
  // "crypto/tls"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"github.com/Azure/azure-sdk-for-go/core/tls"
)

const (
	libraryVersion = "0.1"
	apiVersion     = "v2"
	userAgent      = "go-pulp/" + libraryVersion
)

type Client struct {
	client     		*http.Client
	DisableSsl       bool
	InsecureSkipVerify bool
	baseURL   *url.URL
	UserAgent string
	apiUser   string
	apiPasswd string

	// Services used for talking to different parts of the Pulp API.
	Repositories *RepositoriesService
	Tasks        *TasksService
}

type ListOptions struct {
	Page    int `url:"page,omitempty" json:"page,omitempty"`
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty"`
}

func NewClient(host string, User string, Passwd string, DisableSsl bool, InsecureSkipVerify bool, httpClient *http.Client) (client *Client, err error) {

	ssl := &tls.Config{}
	if InsecureSkipVerify {
		ssl.InsecureSkipVerify =  true
	}

	transport := &http.Transport{
		TLSClientConfig: 			 ssl,
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: transport,
		}
	}

	client = &Client{
		client:    					httpClient,
		UserAgent: 					userAgent,
		apiUser:   					User,
		apiPasswd: 					Passwd,
		DisableSsl:   		 	DisableSsl,
		InsecureSkipVerify: InsecureSkipVerify,
	}

	if err := client.SetHost(host); err != nil {
		return nil, err
	}

	client.Repositories = &RepositoriesService{client: client}
	client.Tasks = &TasksService{client: client}

	return
}

func (c *Client) SetHost(hostStr string) error {
	var err error

	p := "https"
	if c.DisableSsl {
		p = "http"
	}

	err = c.SetBaseURL(p + "://" + hostStr + "/pulp/api/" + apiVersion + "/")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

func (c *Client) SetBaseURL(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseURL, err = url.Parse(urlStr)
	if err != nil {
		return err
	}

	return nil
}

type Response struct {
	*http.Response
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*http.Request, error) {
	u := *c.baseURL
	// Set the encoded opaque data
	u.Opaque = c.baseURL.Path + path

	q, err := query.Values(opt)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}

	if opt != nil && (method == "POST" || method == "PUT") {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.apiUser, c.apiPasswd)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return response, err
}

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

// Pulp Api docs:
// http://pulp.readthedocs.org/en/latest/dev-guide/conventions/sync-v-async.html#call-report
type CallReport struct {
	Result       string `json:"result"`
	Error        *Error `json:"error"`
	SpawnedTasks []struct {
		Href   string `json:"_href"`
		TaskId string `json:"task_id"`
	} `json:"spawned_tasks"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error: %v",
		e.Code, e.Description)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
