package main

import (
	"github.com/labstack/echo/v4"
)

func checkFormParam(c echo.Context, key string) string {
	ret := c.FormValue(key)

	if ret == "" {
		return "FORM DOESN'T HAVE KEY " + key + "\n"
	}

	return ""
}

func checkQueryParam(c echo.Context, key string) string {
	ret := c.QueryParam(key)

	if ret == "" {
		return "QUERY DOESN'T HAVE KEY " + key + "\n"
	}

	return ""
}

type APIError struct {
	Response  *APIResponse `json:"response" xml:"response"`
	Details   any          `json:"details" xml:"details"`
	ErrorCode int          `json:"errorCode" xml:"errorCode"`
}

type APIResponse struct {
	Message string `json:"message" xml:"message"`
}
