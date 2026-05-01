package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
)

func JSONResponse(status int, title, detail string, err error) fuego.HTTPError {
	return fuego.HTTPError{
		Status: status,
		Title:  title,
		Detail: detail,
		Err:    err,
	}
}

func ValidationResponse(err error) fuego.HTTPError {
	return JSONResponse(
		http.StatusBadRequest,
		"invalid request body",
		formatValidationError(err),
		err,
	)
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}

	var parts []string
	for _, fe := range ve {
		switch fe.Tag() {
		case "required":
			parts = append(parts, fmt.Sprintf("%s is required", fe.Field()))
		case "oneof":
			parts = append(
				parts,
				fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param()),
			)
		case "email":
			parts = append(
				parts,
				fmt.Sprintf("%s must be a valid email", fe.Field()),
			)
		case "datetime":
			parts = append(
				parts,
				fmt.Sprintf("%s must match format %s", fe.Field(), fe.Param()),
			)
		default:
			parts = append(parts, fmt.Sprintf("%s is invalid", fe.Field()))
		}
	}

	return strings.Join(parts, "; ")
}
