// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version 2.4.1 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// Defines values for DownloadFontParamsSubset.
const (
	Cyrillic    DownloadFontParamsSubset = "cyrillic"
	CyrillicExt DownloadFontParamsSubset = "cyrillic-ext"
	Greek       DownloadFontParamsSubset = "greek"
	GreekExt    DownloadFontParamsSubset = "greek-ext"
	Hebrew      DownloadFontParamsSubset = "hebrew"
	Latin       DownloadFontParamsSubset = "latin"
	LatinExt    DownloadFontParamsSubset = "latin-ext"
	Vietnamese  DownloadFontParamsSubset = "vietnamese"
)

// Defines values for DownloadFontParamsStyle.
const (
	Italic DownloadFontParamsStyle = "italic"
	Normal DownloadFontParamsStyle = "normal"
)

// DownloadFontParamsSubset defines parameters for DownloadFont.
type DownloadFontParamsSubset string

// DownloadFontParamsStyle defines parameters for DownloadFont.
type DownloadFontParamsStyle string

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetFonts request
	GetFonts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DownloadFont request
	DownloadFont(ctx context.Context, id string, subset DownloadFontParamsSubset, weight string, style DownloadFontParamsStyle, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DownloadLicense request
	DownloadLicense(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSubsets request
	GetSubsets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetFonts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetFontsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DownloadFont(ctx context.Context, id string, subset DownloadFontParamsSubset, weight string, style DownloadFontParamsStyle, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDownloadFontRequest(c.Server, id, subset, weight, style)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DownloadLicense(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDownloadLicenseRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetSubsets(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSubsetsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetFontsRequest generates requests for GetFonts
func NewGetFontsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/fonts.json")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDownloadFontRequest generates requests for DownloadFont
func NewDownloadFontRequest(server string, id string, subset DownloadFontParamsSubset, weight string, style DownloadFontParamsStyle) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "subset", runtime.ParamLocationPath, subset)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "weight", runtime.ParamLocationPath, weight)
	if err != nil {
		return nil, err
	}

	var pathParam3 string

	pathParam3, err = runtime.StyleParamWithLocation("simple", false, "style", runtime.ParamLocationPath, style)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/fonts/%s_%s_%s_%s.woff2", pathParam0, pathParam1, pathParam2, pathParam3)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDownloadLicenseRequest generates requests for DownloadLicense
func NewDownloadLicenseRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/licenses/%s-LICENSE.txt", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetSubsetsRequest generates requests for GetSubsets
func NewGetSubsetsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/subsets.json")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetFontsWithResponse request
	GetFontsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetFontsResponse, error)

	// DownloadFontWithResponse request
	DownloadFontWithResponse(ctx context.Context, id string, subset DownloadFontParamsSubset, weight string, style DownloadFontParamsStyle, reqEditors ...RequestEditorFn) (*DownloadFontResponse, error)

	// DownloadLicenseWithResponse request
	DownloadLicenseWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*DownloadLicenseResponse, error)

	// GetSubsetsWithResponse request
	GetSubsetsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSubsetsResponse, error)
}

type GetFontsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]struct {
		// Designer Name(s) of the designer(s)
		Designer string `json:"designer"`

		// Id Unique identifier for the font family
		Id string `json:"id"`

		// Name Name of the font family
		Name string `json:"name"`

		// Styles Available styles for the font family
		Styles []GetFonts200Styles `json:"styles"`

		// Subsets Available subsets for the font family
		Subsets []GetFonts200Subsets `json:"subsets"`

		// Weights Available font weights for the font family
		Weights []string `json:"weights"`
	}
}
type GetFonts200Styles string
type GetFonts200Subsets string

// Status returns HTTPResponse.Status
func (r GetFontsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetFontsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DownloadFontResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DownloadFontResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DownloadFontResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DownloadLicenseResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DownloadLicenseResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DownloadLicenseResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSubsetsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]struct {
		// Ranges The Unicode ranges covered by the subset, formatted as a comma-separated list of hexadecimal ranges
		Ranges string `json:"ranges"`

		// Subset The name of the subset
		Subset GetSubsets200Subset `json:"subset"`
	}
}
type GetSubsets200Subset string

// Status returns HTTPResponse.Status
func (r GetSubsetsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSubsetsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetFontsWithResponse request returning *GetFontsResponse
func (c *ClientWithResponses) GetFontsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetFontsResponse, error) {
	rsp, err := c.GetFonts(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetFontsResponse(rsp)
}

// DownloadFontWithResponse request returning *DownloadFontResponse
func (c *ClientWithResponses) DownloadFontWithResponse(ctx context.Context, id string, subset DownloadFontParamsSubset, weight string, style DownloadFontParamsStyle, reqEditors ...RequestEditorFn) (*DownloadFontResponse, error) {
	rsp, err := c.DownloadFont(ctx, id, subset, weight, style, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDownloadFontResponse(rsp)
}

// DownloadLicenseWithResponse request returning *DownloadLicenseResponse
func (c *ClientWithResponses) DownloadLicenseWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*DownloadLicenseResponse, error) {
	rsp, err := c.DownloadLicense(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDownloadLicenseResponse(rsp)
}

// GetSubsetsWithResponse request returning *GetSubsetsResponse
func (c *ClientWithResponses) GetSubsetsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSubsetsResponse, error) {
	rsp, err := c.GetSubsets(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSubsetsResponse(rsp)
}

// ParseGetFontsResponse parses an HTTP response from a GetFontsWithResponse call
func ParseGetFontsResponse(rsp *http.Response) (*GetFontsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetFontsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []struct {
			// Designer Name(s) of the designer(s)
			Designer string `json:"designer"`

			// Id Unique identifier for the font family
			Id string `json:"id"`

			// Name Name of the font family
			Name string `json:"name"`

			// Styles Available styles for the font family
			Styles []GetFonts200Styles `json:"styles"`

			// Subsets Available subsets for the font family
			Subsets []GetFonts200Subsets `json:"subsets"`

			// Weights Available font weights for the font family
			Weights []string `json:"weights"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDownloadFontResponse parses an HTTP response from a DownloadFontWithResponse call
func ParseDownloadFontResponse(rsp *http.Response) (*DownloadFontResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DownloadFontResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseDownloadLicenseResponse parses an HTTP response from a DownloadLicenseWithResponse call
func ParseDownloadLicenseResponse(rsp *http.Response) (*DownloadLicenseResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DownloadLicenseResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetSubsetsResponse parses an HTTP response from a GetSubsetsWithResponse call
func ParseGetSubsetsResponse(rsp *http.Response) (*GetSubsetsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetSubsetsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []struct {
			// Ranges The Unicode ranges covered by the subset, formatted as a comma-separated list of hexadecimal ranges
			Ranges string `json:"ranges"`

			// Subset The name of the subset
			Subset GetSubsets200Subset `json:"subset"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
