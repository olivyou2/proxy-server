package main

import (
	"net/http"
	"time"
)

type HttpResponse struct {
	code     string
	datetime string
}

func createHttpResponse(code string) *HttpResponse {
	resp := new(HttpResponse)
	resp.code = code
	resp.datetime = time.Now().UTC().Format(http.TimeFormat)

	return resp
}

func (resp *HttpResponse) createContent(payload string) string {
	result := ""
	result += "HTTP/1.1 " + resp.code + "OK\r\n"
	result += "Date " + resp.datetime + "\r\n"
	result += "Content-Type: text/html" + "\r\n"
	result += "\r\n"
	result += payload

	return result
}
