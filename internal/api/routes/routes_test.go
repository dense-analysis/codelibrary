package routes_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/database/databasemock"
	"github.com/dense-analysis/codelibrary/internal/api/errorhandler"
	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/dense-analysis/ranges"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type RouteTester struct {
	t   *testing.T
	App *fiber.App
	Ctx *fiber.Ctx
	DB  *databasemock.MockDatabaseAPI
}

func NewRouteTester(t *testing.T) RouteTester {
	app := fiber.New()

	return RouteTester{
		t:   t,
		App: app,
		Ctx: app.AcquireCtx(&fasthttp.RequestCtx{}),
		DB:  databasemock.New(),
	}
}

func (r RouteTester) Release() {
	r.App.ReleaseCtx(r.Ctx)
}

const maxParams = 30

// forceFieldAccess forces Go reflection to allow access to unexported fields.
func forceFieldAccess(field reflect.Value) reflect.Value {
	pointer := reflect.NewAt(
		field.Type(),
		unsafe.Pointer(field.UnsafeAddr()),
	)

	return pointer.Elem()
}

func (r *RouteTester) AssertStatus(
	wrappedHandler func(db database.DatabaseAPI) fiber.Handler,
	expectedStatusCode int,
) {
	err := wrappedHandler(r.DB)(r.Ctx)

	if err != nil {
		// Pass the error through the error handler.
		err = errorhandler.ErrorHandler(r.Ctx, err)
	}

	r.t.Helper()
	statusCode := r.Ctx.Response().StatusCode()

	if err != nil || statusCode != expectedStatusCode {
		r.t.Fatalf("Request status %d != %d, error: %v", statusCode, expectedStatusCode, err)
	}
}

func (r *RouteTester) SetParams(params ...ranges.Pair[string, string]) {
	// Get the route or the default route.
	route := r.Ctx.Route()

	// Set the route pointer
	routeField := reflect.ValueOf(r.Ctx).Elem().FieldByName("route")
	routeField = forceFieldAccess(routeField)
	routeField.Set(reflect.ValueOf(route))

	// Get the values out of the
	valuesField := reflect.ValueOf(r.Ctx).Elem().FieldByName("values")
	valuesField = forceFieldAccess(valuesField)
	values := valuesField.Interface().([maxParams]string)

	route.Params = []string{}

	// Clear out previous param values.
	for i := range values {
		values[i] = ""
	}

	// Set new param names and values.
	for i, pair := range params {
		route.Params = append(route.Params, pair.A)
		values[i] = pair.B
	}

	valuesField.Set(reflect.ValueOf(values))
}

func (r *RouteTester) SetQueryArgs(value any) {
	queryString := ""
	structValue := reflect.ValueOf(value)
	structType := reflect.TypeOf(value)

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		if fieldValue.CanInterface() {
			name := fieldType.Tag.Get("query")

			if len(name) == 0 {
				name = fieldType.Name
			}

			value := fieldValue.Interface()

			if len(queryString) > 0 {
				queryString += "&"
			}

			queryString += name
			queryString += "="

			switch v := value.(type) {
			case int:
				queryString += strconv.Itoa(v)
			case int64:
				queryString += strconv.FormatInt(v, 10)
			case uint64:
				queryString += strconv.FormatUint(v, 10)
			case string:
				queryString += url.QueryEscape(v)
			case []string:
				queryString += ranges.JoinStrings(
					ranges.MapF[string](
						ranges.SliceRange(v),
						url.QueryEscape,
					),
					",",
				)
			default:
				panic("Unhandled type: " + fieldType.Name)
			}
		}
	}

	uri := r.Ctx.Request().URI()
	uri.SetQueryString(queryString)
}

func (r *RouteTester) SetRequestBody(value any) {
	r.t.Helper()
	data, err := json.Marshal(value)

	// Ensure JSON deserialization worked.
	assert.Nil(r.t, err)

	req := r.Ctx.Request()
	req.Header.SetContentType("application/json")
	req.SetBody(data)
}

func (r *RouteTester) GetResponse(v any) {
	r.t.Helper()
	// Load the response data JSON
	data := r.Ctx.Response().Body()
	err := json.Unmarshal(data, v)

	// Ensure JSON serialization worked.
	assert.Nil(r.t, err)
}

func (r *RouteTester) AssertResponseError(details ...models.ErrorLocation) {
	expectedError := models.NewError(details...)
	var actualError models.Error
	r.GetResponse(&actualError)
	assert.Equal(r.t, expectedError, actualError)
}
