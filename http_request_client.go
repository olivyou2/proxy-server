package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
)

type HttpRequestClient interface {
	write()
	recv()
}

type Request struct {
	conn net.Conn

	host string
	uri  string
}

func createNewRequestClient(urlValue string) (*Request, error) {
	urlParsed, err := url.Parse(urlValue)

	if err != nil {
		return nil, err
	}

	req := new(Request)
	req.host = urlParsed.Hostname()
	req.uri = urlParsed.RequestURI()

	isHttp := false
	isHttps := false

	requestHost := urlParsed.Hostname()
	requestPort := "80"

	if urlParsed.Scheme == "http" {
		isHttp = true
	} else if urlParsed.Scheme == "https" {
		isHttps = true
		requestPort = "443"
	} else if urlParsed.Scheme == "" {
		// Consider as http
		isHttp = true
	} else {
		return nil, errors.New("UNKNOWN SCHEME " + urlParsed.Scheme)
	}

	if urlParsed.Port() != "" {
		requestPort = urlParsed.Port()
	}

	if isHttp {
		httpConn, err := net.Dial("tcp", requestHost+":"+requestPort)

		if nil != err {
			panic(err)
		}

		req.conn = httpConn
		// conn = httpConn
	} else if isHttps {
		httpsConn, err := tls.Dial("tcp", requestHost+":"+requestPort, &tls.Config{})

		if nil != err {
			panic(err)
		}

		req.conn = httpsConn
		// conn = httpsConn
	}
	return req, nil
}
