package playground

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ilbw97/validator/rtnmsg"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// custom validator, binder
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

type CustomBinder struct {
	Validator *CustomValidator
}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	return cb.BindAndValidate(i, c)
}

const (
	TagRequired = "required"
	TagBase64   = "base64"
	TagOneof    = "oneof"
	TagEmail    = "email"
	TagMin      = "min"
	TagMax      = "max"
	TagLen      = "len"
	TagURL      = "url"
	TagNumeric  = "numeric"
	TagAlpha    = "alpha"
	TagAlphanum = "alphanumeric"
	TagGte      = "gte"
	TagLte      = "lte"
	TagGt       = "gt"
	TagLt       = "lt"
	TagDatetime = "datetime"
	TagIP       = "ip"
	TagIPv4     = "ipv4"
	TagIPv6     = "ipv6"
	TagUnique   = "unique"
)

// formatValidationError formats a validation error into a user-friendly message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case TagRequired:
		return fmt.Sprintf("%s: '%s' is required", rtnmsg.MsgInvalidParam, err.Field())
	case TagBase64:
		return fmt.Sprintf("%s: '%s' should be base64 encoded", rtnmsg.MsgInvalidParam, err.Field())
	case TagOneof:
		possibleValues := strings.ReplaceAll(err.Param(), " ", ", ")
		return fmt.Sprintf("%s: '%s', possible values: '%s'", rtnmsg.MsgInvalidParam, err.Field(), possibleValues)
	case TagEmail:
		return fmt.Sprintf("%s: '%s' should be a valid email address", rtnmsg.MsgInvalidParam, err.Field())
	case TagMin:
		return fmt.Sprintf("%s: '%s' should be at least %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagMax:
		return fmt.Sprintf("%s: '%s' should be at most %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagLen:
		return fmt.Sprintf("%s: '%s' should be exactly %s characters long", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagURL:
		return fmt.Sprintf("%s: '%s' should be a valid URL", rtnmsg.MsgInvalidParam, err.Field())
	case TagNumeric:
		return fmt.Sprintf("%s: '%s' should be numeric", rtnmsg.MsgInvalidParam, err.Field())
	case TagAlpha:
		return fmt.Sprintf("%s: '%s' should contain only alphabetic characters", rtnmsg.MsgInvalidParam, err.Field())
	case TagAlphanum:
		return fmt.Sprintf("%s: '%s' should contain only alphanumeric characters", rtnmsg.MsgInvalidParam, err.Field())
	case TagGte:
		return fmt.Sprintf("%s: '%s' should be greater than or equal to %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagLte:
		return fmt.Sprintf("%s: '%s' should be less than or equal to %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagGt:
		return fmt.Sprintf("%s: '%s' should be greater than %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagLt:
		return fmt.Sprintf("%s: '%s' should be less than %s", rtnmsg.MsgInvalidParam, err.Field(), err.Param())
	case TagDatetime:
		return fmt.Sprintf("%s: '%s' should be a valid datetime", rtnmsg.MsgInvalidParam, err.Field())
	case TagIP:
		return fmt.Sprintf("%s: '%s' should be a valid IP address", rtnmsg.MsgInvalidParam, err.Field())
	case TagIPv4:
		return fmt.Sprintf("%s: '%s' should be a valid IPv4 address", rtnmsg.MsgInvalidParam, err.Field())
	case TagIPv6:
		return fmt.Sprintf("%s: '%s' should be a valid IPv6 address", rtnmsg.MsgInvalidParam, err.Field())
	case TagUnique:
		return fmt.Sprintf("%s: '%s' should contain unique values", rtnmsg.MsgInvalidParam, err.Field())
	default:
		log.Debugf("%s: field[%s], tag[%s], type[%s], value[%s], kind[%s], param[%s]",
			rtnmsg.MsgInvalidParam, err.Field(), err.Tag(), err.Type(), err.Value(),
			err.Kind().String(), err.Param())
		return fmt.Sprintf("%s: '%s'", rtnmsg.MsgInvalidParam, err.Field())
	}
}

// handleBindingError processes binding errors and returns appropriate HTTP errors
func handleBindingError(err error, c echo.Context) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		if httpErr.Code == http.StatusUnsupportedMediaType {
			return handleContentTypeError(c, httpErr)
		}
		return handleInternalBindingError(httpErr)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("unknown binding error: %v", err))
}

// handleContentTypeError processes content type related errors
func handleContentTypeError(c echo.Context, httpErr *echo.HTTPError) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "" {
		log.Debugf("Content-Type header missing")
		return echo.NewHTTPError(http.StatusUnsupportedMediaType,
			fmt.Sprintf("%s: '%s'", rtnmsg.MsgUnsupportedContentType, rtnmsg.MsgContentTypeMissing))
	}
	if contentType != echo.MIMEApplicationJSON {
		log.Debugf("Unsupported Content-Type: %s", contentType)
		return echo.NewHTTPError(http.StatusUnsupportedMediaType,
			fmt.Sprintf("%s: '%s'", rtnmsg.MsgUnsupportedContentType, contentType))
	}
	return echo.NewHTTPError(http.StatusUnsupportedMediaType,
		fmt.Sprintf("%s: '%s'", rtnmsg.MsgUnsupportedContentType, httpErr.Message))
}

// handleInternalBindingError processes internal binding errors
func handleInternalBindingError(httpErr *echo.HTTPError) error {
	switch internalErr := httpErr.Internal.(type) {
	case *json.SyntaxError:
		log.Debugf("json syntax error: %s", internalErr.Error())
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("invalid json syntax: %s", internalErr.Error()))
	case *json.UnmarshalTypeError:
		log.Debugf("type error: field[%s], type[%s], struct[%s], value[%s]",
			internalErr.Field, internalErr.Type, internalErr.Struct, internalErr.Value)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("%s: '%s'", rtnmsg.MsgInvalidParam, internalErr.Field))
	case *json.InvalidUnmarshalError:
		log.Debugf("invalid unmarshal error: %s", internalErr)
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("%s: '%s'", rtnmsg.MsgInvalidParam, internalErr.Type))
	default:
		log.Errorf("other binding error: %v", httpErr)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("%s: %v", rtnmsg.MsgInternalError, httpErr))
	}
}

// validateStruct performs validation on the struct and returns the first validation error
func (cb *CustomBinder) validateStruct(i interface{}) error {
	if cb.Validator == nil {
		log.Error("Validator is not initialized")
		return errors.New(rtnmsg.MsgInternalError)
	}

	if err := cb.Validator.Validate(i); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok && len(validationErrs) > 0 {
			// Return only the first validation error
			return echo.NewHTTPError(http.StatusBadRequest, formatValidationError(validationErrs[0]))
		}
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("unknown validation error: %v", err))
	}
	return nil
}

// BindAndValidate binds and validates the request
func (cb *CustomBinder) BindAndValidate(i interface{}, c echo.Context) error {
	db := new(echo.DefaultBinder)

	if err := db.Bind(i, c); err != nil {
		return handleBindingError(err, c)
	}

	return cb.validateStruct(i)
}
