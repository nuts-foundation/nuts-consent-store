// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ConsentCheckRequest defines model for ConsentCheckRequest.
type ConsentCheckRequest struct {
	Actor        Identifier `json:"actor"`
	Custodian    Identifier `json:"custodian"`
	ResourceType string     `json:"resourceType"`
	Subject      Identifier `json:"subject"`
	ValidAt      *string    `json:"validAt,omitempty"`
}

// ConsentCheckResponse defines model for ConsentCheckResponse.
type ConsentCheckResponse struct {
	ConsentGiven *string `json:"consentGiven,omitempty"`
	Limitations  *string `json:"limitations,omitempty"`
}

// ConsentQueryRequest defines model for ConsentQueryRequest.
type ConsentQueryRequest struct {
	Actor     *Identifier     `json:"actor,omitempty"`
	Custodian *Identifier     `json:"custodian,omitempty"`
	Page      *PageDefinition `json:"page,omitempty"`
	Query     interface{}     `json:"query"`
	ValidAt   *time.Time      `json:"validAt,omitempty"`
}

// ConsentQueryResponse defines model for ConsentQueryResponse.
type ConsentQueryResponse struct {
	Page         PageDefinition      `json:"page"`
	Results      []SimplifiedConsent `json:"results"`
	TotalResults int                 `json:"totalResults"`
}

// Identifier defines model for Identifier.
type Identifier string

// PageDefinition defines model for PageDefinition.
type PageDefinition struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// SimplifiedConsent defines model for SimplifiedConsent.
type SimplifiedConsent struct {
	Actor      Identifier `json:"actor"`
	Custodian  Identifier `json:"custodian"`
	Id         string     `json:"id"`
	RecordHash *string    `json:"recordHash,omitempty"`
	Resources  []string   `json:"resources"`
	Subject    Identifier `json:"subject"`
	ValidFrom  ValidFrom  `json:"validFrom"`
	ValidTo    ValidTo    `json:"validTo"`
}

// SubjectQuery defines model for SubjectQuery.
type SubjectQuery string

// ValidFrom defines model for ValidFrom.
type ValidFrom string

// ValidTo defines model for ValidTo.
type ValidTo string

// createConsentJSONBody defines parameters for CreateConsent.
type createConsentJSONBody SimplifiedConsent

// checkConsentJSONBody defines parameters for CheckConsent.
type checkConsentJSONBody ConsentCheckRequest

// queryConsentJSONBody defines parameters for QueryConsent.
type queryConsentJSONBody ConsentQueryRequest

// CreateConsentRequestBody defines body for CreateConsent for application/json ContentType.
type CreateConsentJSONRequestBody createConsentJSONBody

// CheckConsentRequestBody defines body for CheckConsent for application/json ContentType.
type CheckConsentJSONRequestBody checkConsentJSONBody

// QueryConsentRequestBody defines body for QueryConsent for application/json ContentType.
type QueryConsentJSONRequestBody queryConsentJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(req *http.Request, ctx context.Context) error

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// HTTP client with any customized settings, such as certificate chains.
	Client http.Client

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// The interface specification for the client above.
type ClientInterface interface {
	// CreateConsent request  with any body
	CreateConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	CreateConsent(ctx context.Context, body CreateConsentJSONRequestBody) (*http.Response, error)

	// CheckConsent request  with any body
	CheckConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	CheckConsent(ctx context.Context, body CheckConsentJSONRequestBody) (*http.Response, error)

	// QueryConsent request  with any body
	QueryConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	QueryConsent(ctx context.Context, body QueryConsentJSONRequestBody) (*http.Response, error)

	// DeleteConsent request
	DeleteConsent(ctx context.Context, proofHash string) (*http.Response, error)
}

func (c *Client) CreateConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewCreateConsentRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) CreateConsent(ctx context.Context, body CreateConsentJSONRequestBody) (*http.Response, error) {
	req, err := NewCreateConsentRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) CheckConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewCheckConsentRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) CheckConsent(ctx context.Context, body CheckConsentJSONRequestBody) (*http.Response, error) {
	req, err := NewCheckConsentRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) QueryConsentWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewQueryConsentRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) QueryConsent(ctx context.Context, body QueryConsentJSONRequestBody) (*http.Response, error) {
	req, err := NewQueryConsentRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteConsent(ctx context.Context, proofHash string) (*http.Response, error) {
	req, err := NewDeleteConsentRequest(c.Server, proofHash)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewCreateConsentRequest calls the generic CreateConsent builder with application/json body
func NewCreateConsentRequest(server string, body CreateConsentJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateConsentRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateConsentRequestWithBody generates requests for CreateConsent with any type of body
func NewCreateConsentRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl := fmt.Sprintf("%s/consent", server)

	req, err := http.NewRequest("POST", queryUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewCheckConsentRequest calls the generic CheckConsent builder with application/json body
func NewCheckConsentRequest(server string, body CheckConsentJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCheckConsentRequestWithBody(server, "application/json", bodyReader)
}

// NewCheckConsentRequestWithBody generates requests for CheckConsent with any type of body
func NewCheckConsentRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl := fmt.Sprintf("%s/consent/check", server)

	req, err := http.NewRequest("POST", queryUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewQueryConsentRequest calls the generic QueryConsent builder with application/json body
func NewQueryConsentRequest(server string, body QueryConsentJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryConsentRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryConsentRequestWithBody generates requests for QueryConsent with any type of body
func NewQueryConsentRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl := fmt.Sprintf("%s/consent/query", server)

	req, err := http.NewRequest("POST", queryUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewDeleteConsentRequest generates requests for DeleteConsent
func NewDeleteConsentRequest(server string, proofHash string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "proofHash", proofHash)
	if err != nil {
		return nil, err
	}

	queryUrl := fmt.Sprintf("%s/consent/%s", server, pathParam0)

	req, err := http.NewRequest("DELETE", queryUrl, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses returns a ClientWithResponses with a default Client:
func NewClientWithResponses(server string) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client: http.Client{},
			Server: server,
		},
	}
}

// NewClientWithResponsesAndRequestEditorFunc takes in a RequestEditorFn callback function and returns a ClientWithResponses with a default Client:
func NewClientWithResponsesAndRequestEditorFunc(server string, reqEditorFn RequestEditorFn) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client:        http.Client{},
			Server:        server,
			RequestEditor: reqEditorFn,
		},
	}
}

type createConsentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r createConsentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r createConsentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type checkConsentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ConsentCheckResponse
}

// Status returns HTTPResponse.Status
func (r checkConsentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r checkConsentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type queryConsentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ConsentQueryResponse
}

// Status returns HTTPResponse.Status
func (r queryConsentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r queryConsentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type deleteConsentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r deleteConsentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r deleteConsentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateConsentWithBodyWithResponse request with arbitrary body returning *CreateConsentResponse
func (c *ClientWithResponses) CreateConsentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*createConsentResponse, error) {
	rsp, err := c.CreateConsentWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsecreateConsentResponse(rsp)
}

func (c *ClientWithResponses) CreateConsentWithResponse(ctx context.Context, body CreateConsentJSONRequestBody) (*createConsentResponse, error) {
	rsp, err := c.CreateConsent(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsecreateConsentResponse(rsp)
}

// CheckConsentWithBodyWithResponse request with arbitrary body returning *CheckConsentResponse
func (c *ClientWithResponses) CheckConsentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*checkConsentResponse, error) {
	rsp, err := c.CheckConsentWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsecheckConsentResponse(rsp)
}

func (c *ClientWithResponses) CheckConsentWithResponse(ctx context.Context, body CheckConsentJSONRequestBody) (*checkConsentResponse, error) {
	rsp, err := c.CheckConsent(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsecheckConsentResponse(rsp)
}

// QueryConsentWithBodyWithResponse request with arbitrary body returning *QueryConsentResponse
func (c *ClientWithResponses) QueryConsentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*queryConsentResponse, error) {
	rsp, err := c.QueryConsentWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsequeryConsentResponse(rsp)
}

func (c *ClientWithResponses) QueryConsentWithResponse(ctx context.Context, body QueryConsentJSONRequestBody) (*queryConsentResponse, error) {
	rsp, err := c.QueryConsent(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsequeryConsentResponse(rsp)
}

// DeleteConsentWithResponse request returning *DeleteConsentResponse
func (c *ClientWithResponses) DeleteConsentWithResponse(ctx context.Context, proofHash string) (*deleteConsentResponse, error) {
	rsp, err := c.DeleteConsent(ctx, proofHash)
	if err != nil {
		return nil, err
	}
	return ParsedeleteConsentResponse(rsp)
}

// ParsecreateConsentResponse parses an HTTP response from a CreateConsentWithResponse call
func ParsecreateConsentResponse(rsp *http.Response) (*createConsentResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &createConsentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParsecheckConsentResponse parses an HTTP response from a CheckConsentWithResponse call
func ParsecheckConsentResponse(rsp *http.Response) (*checkConsentResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &checkConsentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		response.JSON200 = &ConsentCheckResponse{}
		if err := json.Unmarshal(bodyBytes, response.JSON200); err != nil {
			return nil, err
		}

	}

	return response, nil
}

// ParsequeryConsentResponse parses an HTTP response from a QueryConsentWithResponse call
func ParsequeryConsentResponse(rsp *http.Response) (*queryConsentResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &queryConsentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		response.JSON200 = &ConsentQueryResponse{}
		if err := json.Unmarshal(bodyBytes, response.JSON200); err != nil {
			return nil, err
		}

	}

	return response, nil
}

// ParsedeleteConsentResponse parses an HTTP response from a DeleteConsentWithResponse call
func ParsedeleteConsentResponse(rsp *http.Response) (*deleteConsentResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &deleteConsentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new consent record for a C-S-A combination.// (POST /consent)
	CreateConsent(ctx echo.Context) error
	// Send a request for checking if the given combination exists// (POST /consent/check)
	CheckConsent(ctx echo.Context) error
	// Do a query for available consent// (POST /consent/query)
	QueryConsent(ctx echo.Context) error
	// Remove a consent record for a C-S-A combination.// (DELETE /consent/{proofHash})
	DeleteConsent(ctx echo.Context, proofHash string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateConsent converts echo context to params.
func (w *ServerInterfaceWrapper) CreateConsent(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateConsent(ctx)
	return err
}

// CheckConsent converts echo context to params.
func (w *ServerInterfaceWrapper) CheckConsent(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CheckConsent(ctx)
	return err
}

// QueryConsent converts echo context to params.
func (w *ServerInterfaceWrapper) QueryConsent(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.QueryConsent(ctx)
	return err
}

// DeleteConsent converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteConsent(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "proofHash" -------------
	var proofHash string

	err = runtime.BindStyledParameter("simple", false, "proofHash", ctx.Param("proofHash"), &proofHash)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter proofHash: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteConsent(ctx, proofHash)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/consent", wrapper.CreateConsent)
	router.POST("/consent/check", wrapper.CheckConsent)
	router.POST("/consent/query", wrapper.QueryConsent)
	router.DELETE("/consent/:proofHash", wrapper.DeleteConsent)

}

