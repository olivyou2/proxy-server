package main

import (
	"crypto/tls"
	"net"
	"strings"
)

type HttpServer struct {
	server *net.Listener

	connectedHandler func(client *HttpClient, req *HttpRequest)
	dataHandler      func(client *HttpClient, req *HttpRequest, data []byte)
}

type HttpClient struct {
	connection   *net.Conn
	data         string
	headerParsed bool
}

func createHttpClient(connection *net.Conn) *HttpClient {
	client := new(HttpClient)
	client.connection = connection

	return client
}

func (s *HttpServer) listen(address string) {

	server, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}

	s.server = &server
	s.accept()
}

func (s *HttpServer) listenTls(address string, cert string, privkey string) {
	cer, err := tls.LoadX509KeyPair(cert, privkey)

	if err != nil {
		panic(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	server, err := tls.Listen("tcp", address, config)
	// server, err := net.Listen("tcp", address)

	if err != nil {
		panic(err)
	}

	s.server = &server
	s.accept()
}

func (s *HttpServer) accept() {
	for {
		conn, err := (*s.server).Accept()

		if err != nil {
			panic(err)
		}

		client := createHttpClient(&conn)

		go s.recv(client)
	}
}

func (s *HttpServer) recv(client *HttpClient) {
	conn := client.connection
	buffer := make([]byte, 1024)

	var context *HttpRequest = nil

	for {
		n, err := (*conn).Read(buffer)

		if nil != err || n == 0 {
			break
		}

		buffer_string := buffer[:n]

		if client.headerParsed {

			if s.dataHandler != nil && context != nil {
				s.dataHandler(client, context, buffer_string)
			}
		} else {
			client.data += string(buffer_string)

			if strings.Contains(client.data, "\r\n\r\n") {
				// All headers arived
				client.headerParsed = true

				sep_index := strings.Index(client.data, "\r\n\r\n")

				headers := client.data[:sep_index]
				client.data = client.data[sep_index+4:]

				context = parseHeader(headers)

				if s.connectedHandler != nil {
					s.connectedHandler(client, context)
				}
			}
		}

	}
}

func parseHeader(headerText string) *HttpRequest {

	lines := strings.Split(headerText, "\r\n")
	request := createRequestContext()

	firstLine := lines[0]
	info := strings.Split(firstLine, " ")

	method := info[0]
	uri := info[1]
	version := info[2]

	request.method = method
	request.uri = uri
	request.version = version

	lines = lines[1:]

	for _, line := range lines {
		kv := strings.Split(line, ":")

		key := strings.TrimSpace(kv[0])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(strings.Join(kv[1:], ":"))

		hkv := HeaderPair{key: key, value: value}

		request.headers = append(request.headers, hkv)

		lowerKey := strings.ToLower(key)

		if lowerKey == "host" {
			request.host = value
		}
	}
	request.headerOriginal = headerText

	return request
}

func createHttpServer() *HttpServer {
	server := new(HttpServer)

	return server
}
