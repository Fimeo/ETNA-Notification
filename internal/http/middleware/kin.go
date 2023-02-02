package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
)

type KinValidator struct {
	R routers.Router
}

// NewKinValidator service performs input request validation from OpenApi V3 documentation.
// The documentation could be found in api/swagger.json directory.
//
// Kin service loads and validate the openapi file and boostrap internal gorillamux router.
func NewKinValidator() KinValidator {
	doc, err := openapi3.NewLoader().LoadFromFile("api/swagger.json")
	if err != nil {
		panic(fmt.Sprintf("cannot read swagger json : %s", err.Error()))
	}

	ctx := context.Background()
	err = doc.Validate(ctx)
	if err != nil {
		panic(fmt.Sprintf("cannot validate swagger json : %s", err.Error()))
	}

	router, _ := gorillamux.NewRouter(doc)

	return KinValidator{R: router}
}

// Kin middleware performs input request validation from service.Kin configuration.
//
// This middleware checks in first time in the request path and method is defined in openapi
// documentation (by matching mux router). If the request does not match any route, 404 error page
// is sent back to client.
//
// Then, is request go through first step, the input body, query parameters, headers, content type and so on
// are validated against the defined expected types, regular expressions, enums...
// If the input body does is incorrect, http.StatusUnprocessableEntity is sent back to client.
func Kin(kin KinValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		route, params, err := kin.R.FindRoute(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "route not found"})
			return
		}
		if err := openapi3filter.ValidateRequest(c, &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: params,
			Route:      route,
		}); err != nil {
			var schemaError *openapi3.SchemaError
			if errors.As(err, &schemaError) {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": schemaError.Reason})
			} else {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			}
			return
		}

		c.Next()
	}
}
