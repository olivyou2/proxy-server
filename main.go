package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	proxyHost := "0.0.0.0"

	listenHttp := true
	listenHttps := true

	httpPort := "80"
	httpsPort := "443"

	certFile := "fullchain.pem"
	keyFile := "privkey.pem"

	useApi := false
	apiHost := "proxy.default.site"
	apiPort := "1010"

	(&cli.App{
		Name: "proxy-server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       proxyHost,
				Usage:       "To set the server's host",
				Destination: &proxyHost,
			},
			&cli.BoolFlag{
				Name:        "use-api",
				Value:       useApi,
				Usage:       "To use API",
				Destination: &useApi,
			},
			&cli.StringFlag{
				Name:        "api-host",
				Value:       apiHost,
				Usage:       "To use API",
				Destination: &apiHost,
			},
			&cli.StringFlag{
				Name:        "api-port",
				Value:       apiPort,
				Usage:       "To set the API port",
				Destination: &apiPort,
			},
			&cli.BoolFlag{
				Name:        "listen-http",
				Value:       listenHttp,
				Usage:       "To turn on http listener",
				Destination: &listenHttp,
			},
			&cli.BoolFlag{
				Name:        "listen-https",
				Value:       listenHttps,
				Usage:       "To turn on https listener",
				Destination: &listenHttps,
			},
			&cli.StringFlag{
				Name:        "cert-file",
				Value:       certFile,
				Usage:       "To set the TLS certification file",
				Destination: &certFile,
			},
			&cli.StringFlag{
				Name:        "privkey-file",
				Value:       keyFile,
				Usage:       "To set the TLS private key file",
				Destination: &keyFile,
			},
		},
		Usage: "a http(s) proxy fits perfectly with you",
		Action: func(context *cli.Context) error {
			pm := createNewProxyManager()

			if listenHttp {
				tlsServer := createHttpServer()
				tlsServer.connectedHandler = pm.connectHandler
				tlsServer.dataHandler = pm.dataHandler
				go tlsServer.listenTls(proxyHost+":"+httpsPort, certFile, keyFile)

				fmt.Println("HTTPS proxy server is listening on " + proxyHost + ":" + httpsPort)
			}

			if listenHttps {
				server := createHttpServer()
				server.connectedHandler = pm.connectHandler
				server.dataHandler = pm.dataHandler
				go server.listen(proxyHost + ":" + httpPort)

				fmt.Println("HTTP proxy server is listening on " + proxyHost + ":" + httpPort)
			}

			if useApi {
				fmt.Println("API server is listening on " + proxyHost + ":" + httpPort + " (" + apiHost + ")")

				pma := createNewProxyManagerAPI(proxyHost+":"+apiPort, apiHost, pm)
				go pma.echoServer.Logger.Fatal(pma.echoServer.Start(proxyHost + ":" + apiPort))
			}

			select {}
		},
	}).Run(os.Args)

}
