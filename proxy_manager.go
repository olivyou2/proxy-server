package main

import "strings"

type ProxyManager struct {
	redirectManager *RedirectManager
}

func createNewProxyManager() *ProxyManager {
	manager := new(ProxyManager)
	manager.redirectManager = createNewRedirectManager()

	return manager
}

func (pm *ProxyManager) connectHandler(client *HttpClient, req *HttpRequest) {
	host := req.host

	target, ok := pm.redirectManager.getRedirect(host)

	if !ok {
		resp := createHttpResponse("400")
		data := resp.createContent("No host " + host)

		(*client.connection).Write([]byte(data))
		(*client.connection).Close()
		return
	}

	targetHost := strings.Join(strings.Split(target, "://")[1:], "://")

	request, err := createNewRequestClient(target)

	if err != nil {
		resp := createHttpResponse("400")
		data := resp.createContent(err.Error())

		(*client.connection).Write([]byte(data))
		(*client.connection).Close()
		return
	}

	go func() {
		buffer := make([]byte, 1024)

		for {
			n, err := request.conn.Read(buffer)

			if nil != err || n == 0 {
				(*client.connection).Close()
				break
			}

			(*client.connection).Write(buffer[:n])
		}
	}()

	req.setHeader("Host", targetHost)
	header := req.toHttpText()

	req.connection = &request.conn

	mode := true

	if mode {
		request.conn.Write([]byte(req.headerOriginal))
		request.conn.Write([]byte("\r\n\r\n"))
	} else {
		request.conn.Write([]byte(header))
	}

	request.conn.Write([]byte(client.data))
}

func (pm *ProxyManager) dataHandler(client *HttpClient, req *HttpRequest, data []byte) {

	(*req.connection).Write(data)
}

/**

0x3133A


**/
