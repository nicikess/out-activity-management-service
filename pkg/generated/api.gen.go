// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for RunStatus.
const (
	Active    RunStatus = "active"
	Completed RunStatus = "completed"
	Paused    RunStatus = "paused"
)

// Coordinate defines model for Coordinate.
type Coordinate struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Run defines model for Run.
type Run struct {
	EndTime *time.Time         `json:"endTime,omitempty"`
	Id      openapi_types.UUID `json:"id"`
	Route   *struct {
		Coordinates *[]struct {
			Latitude  *float32   `json:"latitude,omitempty"`
			Longitude *float32   `json:"longitude,omitempty"`
			Timestamp *time.Time `json:"timestamp,omitempty"`
		} `json:"coordinates,omitempty"`
	} `json:"route,omitempty"`
	StartTime time.Time `json:"startTime"`
	Stats     *struct {
		// AveragePace Average pace in meters per second
		AveragePace *float32 `json:"averagePace,omitempty"`

		// Distance Total distance in meters
		Distance *float32 `json:"distance,omitempty"`

		// Duration Total duration in seconds
		Duration *int `json:"duration,omitempty"`
	} `json:"stats,omitempty"`
	Status RunStatus          `json:"status"`
	UserId openapi_types.UUID `json:"userId"`
}

// RunStatus defines model for Run.Status.
type RunStatus string

// StartRunRequest defines model for StartRunRequest.
type StartRunRequest struct {
	InitialLocation Coordinate `json:"initialLocation"`
}

// StartRunJSONRequestBody defines body for StartRun for application/json ContentType.
type StartRunJSONRequestBody = StartRunRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Start a new run
	// (POST /runs)
	StartRun(ctx echo.Context) error
	// Get user's active run
	// (GET /runs/active)
	GetActiveRun(ctx echo.Context) error
	// Get run by ID
	// (GET /runs/{runId})
	GetRun(ctx echo.Context, runId openapi_types.UUID) error
	// End an active or paused run
	// (PUT /runs/{runId}/end)
	EndRun(ctx echo.Context, runId openapi_types.UUID) error
	// Pause an active run
	// (PUT /runs/{runId}/pause)
	PauseRun(ctx echo.Context, runId openapi_types.UUID) error
	// Resume a paused run
	// (PUT /runs/{runId}/resume)
	ResumeRun(ctx echo.Context, runId openapi_types.UUID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// StartRun converts echo context to params.
func (w *ServerInterfaceWrapper) StartRun(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.StartRun(ctx)
	return err
}

// GetActiveRun converts echo context to params.
func (w *ServerInterfaceWrapper) GetActiveRun(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetActiveRun(ctx)
	return err
}

// GetRun converts echo context to params.
func (w *ServerInterfaceWrapper) GetRun(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "runId", runtime.ParamLocationPath, ctx.Param("runId"), &runId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter runId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetRun(ctx, runId)
	return err
}

// EndRun converts echo context to params.
func (w *ServerInterfaceWrapper) EndRun(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "runId", runtime.ParamLocationPath, ctx.Param("runId"), &runId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter runId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.EndRun(ctx, runId)
	return err
}

// PauseRun converts echo context to params.
func (w *ServerInterfaceWrapper) PauseRun(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "runId", runtime.ParamLocationPath, ctx.Param("runId"), &runId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter runId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PauseRun(ctx, runId)
	return err
}

// ResumeRun converts echo context to params.
func (w *ServerInterfaceWrapper) ResumeRun(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "runId", runtime.ParamLocationPath, ctx.Param("runId"), &runId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter runId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ResumeRun(ctx, runId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/runs", wrapper.StartRun)
	router.GET(baseURL+"/runs/active", wrapper.GetActiveRun)
	router.GET(baseURL+"/runs/:runId", wrapper.GetRun)
	router.PUT(baseURL+"/runs/:runId/end", wrapper.EndRun)
	router.PUT(baseURL+"/runs/:runId/pause", wrapper.PauseRun)
	router.PUT(baseURL+"/runs/:runId/resume", wrapper.ResumeRun)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RWTW/jNhD9KwRboBd15bR7qW7ZNl2kaIuFE6CHIIeJOHa4kEjtcOjWCPzfiyHlD9lq",
	"Ey829mFvljh6M/Pem6GfdO3bzjt0HHT1pEP9iC2knz97T8Y6YJSnjnyHxBbTWQNsOZp0MvPUAutKzxoP",
	"rAvNyw51pV1sH5D0qtCNd/MXh68KTfgpWkKjq7ttpl2Y+81X/uEj1ixJrog8HVZa+5y2jw9M1s0lvsUQ",
	"YD52tldBQtjGj+WeRneYGZ25te2wZwOM37O8LQ4rsmYQG6M1Y2Hk45gk9Uau9GgZ23BS4QotjQWGtntp",
	"z6sRLvsXQATL8YjAQHwctYGBR9iABRLM8QPUCctgqMl2bL3Tlb7Mh6qDGpV1qkVGCqpDUgFr70Sc5zkx",
	"NjC4Mfxbz9Co9fk2w8tgI0HG+Q/Y/lxgc7VhC2Md47yftTF2OfYGjq0MANRsF0JrBzGgtC07o0FGszMN",
	"W7JjQLp+iZf3Bi2F9B/vqrypaWz0biRsGt0UP0UMfKixdZYtNL/7ekPYt4QzXelvyu3uK/vFV+5svYP6",
	"9pAOyxH6sI5keXkjeLmCdwiEdBn5UZ4e0tOva2Z+++tWGkzRuupPt1Q9Mnd6JcDWzfyh2jdIC1ujmnlS",
	"LTiYWzdXM8sOQ1BJOMtLRdEl+S03AjqNTv0hwdiiY9Vj6EIvkELGvXgzeTMRfn2HDjqrK/1jeiUu4MfU",
	"V5lghXCfiRfaEzWi/kYZnVnEwO+8WeZV5Rhd+gS6rrGZz/JjyPJkKZ4Tal/41VAupojpRei8C1mIHyYX",
	"Xyy9NJZSDvUQapNz0agQ6xpDmMWmWQqVbyeTL5Y+33YjBVy7BTTWKFrTInl/ev28l2lJiNMUNIRglgr/",
	"sYHDYCp0dTech7v71X2hQ2xboOXaNAqUw78FK32cfFb2W6h60nMccdt75FzC2nED3SevrftO+zMfncm8",
	"v3193v/0ecwHuY9h/D2ykqX7XdgB2uH9iaK7Nqv/I/48lE9PzrVkdJ4/l2UR6GGprn+Rgjsg6O97+dYK",
	"vmxWXWgH6SpIvOv9nVbsNPHc9Xq/L2KJLt3LJ0le6C6O+OXKmfP5BZ0ZXcynN8+J1rJktWGzkrd/3I5y",
	"75UzCtx6P3hS+Z/g+KYo0+G5bfZBijif0XqCvj6nSWK7sYr8ccfjvJaE23HbuMcIQ2zPbrJpquJ8Lsss",
	"fMU2W4/Z8TbL0ikYrrKEQIu1mSI1utIldLZcXOjV/erfAAAA///AAe5gLxMAAA==",
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
