package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	errs "pantho/golang/pkg"
)


func ServeData(c *gin.Context, data any, codes ...int) {
	status := http.StatusOK
	if len(codes) > 0 && codes[0] >= 100 && codes[0] < 600 {
		status = codes[0]
	}
	c.JSON(status, data)
}

func internal(c *gin.Context) {
	err := errs.NewInternalServer("unhandled internal error", "something went wrong!")

	formatter := err.(errs.HTTPFormatter)
	statusCode := err.(errs.StatusCoder)

	c.JSON(statusCode.StatusCode(), formatter.HTTPFormat())
}

// ServeErr inspects an error and sends an appropriate HTTP error response.
// It handles custom errors from the errs package and provides a generic
// internal server error for all other cases.
func ServeErr(c *gin.Context, e error) {
	log.Printf("serving http error: %v", e)

	// Check if the error (or any error it wraps) is a custom error
	// that can be formatted for an HTTP response.
	var customErr interface {
		errs.StatusCoder
		errs.HTTPFormatter
	}
	if errors.As(e, &customErr) {
		c.JSON(customErr.StatusCode(), customErr.HTTPFormat())
		return
	}

	// For any other error type, respond with a generic internal server error.
	internal(c)
}