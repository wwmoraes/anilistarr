// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version 2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// CustomList defines model for CustomList.
type CustomList = []struct {
	TvdbID *float32 `json:"TvdbID,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /user/{name}/id)
	GetUserID(w http.ResponseWriter, r *http.Request, name string)

	// (GET /user/{name}/media)
	GetUserMedia(w http.ResponseWriter, r *http.Request, name string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /user/{name}/id)
func (_ Unimplemented) GetUserID(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /user/{name}/media)
func (_ Unimplemented) GetUserMedia(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetUserID operation middleware
func (siw *ServerInterfaceWrapper) GetUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "name" -------------
	var name string

	name = chi.URLParam(r, "name")

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserID(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetUserMedia operation middleware
func (siw *ServerInterfaceWrapper) GetUserMedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "name" -------------
	var name string

	name = chi.URLParam(r, "name")

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserMedia(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{name}/id", wrapper.GetUserID)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/user/{name}/media", wrapper.GetUserMedia)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{
	"H4sIAAAAAAAC/8xUy27bSgz9lQHvXV0IGie53WjVIE4LA0k3TVdNFrREW0w1j3IoO0Ggfy9mHNc1gj4W",
	"LdCVqWM+z+HwCdrgYvDkNUHzBKntyWExL8akwV1x0vzFSq7AUUIkUabydbPplot5tvQxEjTgR7ckgWmq",
	"9khY3lOrcABQBB9hyi70oCQeh3loS7qOUisclYOHBt6w70wY1bggZHCZTe3JRAklZQWjDNBArxoba9es",
	"/bis2+DsduuCICWLngdOiiK5AfarkMu0wSu2ZS5yyDnHwfE1ipKEuqNNjjlu6SL4DYkmg96c70LMmEjM",
	"FrXt2a9NgTQYNG0hcAesgjhUs+257c1/KHLrMcZk0hhjEK2hgoFb8olyUx5dJup6cXM0Y2qsFdzWu0Fz",
	"2TwIef3ezNZhUhJ7tbi4fPf+smjAOuTc5wdiKtiQpN18s/qknmW/EMljZGjgrEAVRNS+iGRzZfuUm5ws",
	"dxlaUyEzbwZmphYdNPCW9EMiWcxLsKAjJUnQfNwJQL7EKD2ojQOyL3I8oIulwf04ZZXyn6UBqPbklJ8K",
	"hD6PLNRBozLSdJeRFINPu/08nc32gv+43snp2f/TC7mLtKsw+i5z8upXk8EKeaAur8GanhdkMW9MXde3",
	"Hl6WIdaexKARVDIDO1YTxISCckojmS1rX5Z/jEmF0BkVbD+RmB5jJE9dflFTdayOo47xZwJdF6e/SSOM",
	"ceC2NGrvU/CHy5Stf4VW0MA/9nC67PPdst8creklz4WOrw+y0LnmDfmi0O9S+I/qO1WQSDZ7jY5vw+Hh",
	"16vhcXfA7qYvAQAA//8z3Gun4wUAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
