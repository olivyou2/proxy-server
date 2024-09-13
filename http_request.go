package main

import (
	"net"
)

type HeaderPair struct {
	key   string
	value string
}

type HttpRequest struct {
	method         string
	uri            string
	version        string
	host           string
	headers        []HeaderPair
	headerOriginal string

	connection *net.Conn
}

func createRequestContext() *HttpRequest {
	req := new(HttpRequest)

	req.headers = make([]HeaderPair, 0)

	return req
}

func (req HttpRequest) toHttpText() string {
	text := ""

	text += req.method
	text += " "

	text += req.uri
	text += " "

	text += req.version
	text += "\r\n"

	for _, value := range req.headers {
		text += value.key
		text += ": "
		text += value.value
		text += "\r\n"
	}

	text += "\r\n"

	return text
}

func (req *HttpRequest) setHeader(key string, value string) {
	for index, header := range req.headers {
		if header.key == key {
			req.headers[index].value = value
			return
		}
	}

	req.headers = append(req.headers, HeaderPair{key: key, value: value})
}
