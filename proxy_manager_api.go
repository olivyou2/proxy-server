package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ProxyManagerAPI struct {
	proxyManager    *ProxyManager
	redirectManager *RedirectManager

	echoServer *echo.Echo
}

func createNewProxyManagerAPI(address string, apiHost string, proxyManager *ProxyManager) *ProxyManagerAPI {
	pm := new(ProxyManagerAPI)
	pm.proxyManager = proxyManager
	pm.redirectManager = proxyManager.redirectManager
	proxyManager.redirectManager.addNewRedirect(apiHost, "http://"+address)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", pm.hello)

	e.GET("/proxies", pm.getProxies)
	e.POST("/proxies", pm.addProxies)
	e.DELETE("/proxies", pm.removeProxies)

	pm.echoServer = e

	return pm
}

func (pm *ProxyManagerAPI) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (pm *ProxyManagerAPI) getProxies(c echo.Context) error {
	response := make(map[string]string)

	proxyMap := pm.redirectManager.proxyMap
	proxyMap.Range(func(keyRaw any, valRaw any) bool {
		key := keyRaw.(string)
		val := valRaw.(string)

		response[key] = val
		return true
	})

	return c.JSON(http.StatusOK, response)
}

func (pm *ProxyManagerAPI) addProxies(c echo.Context) error {
	errMsg := ""
	errMsg += checkFormParam(c, "host")
	errMsg += checkFormParam(c, "target")

	if errMsg != "" {
		return c.JSON(http.StatusBadRequest, &APIError{
			Response: &APIResponse{
				Message: "Validation error",
			},
			Details:   errMsg,
			ErrorCode: 1,
		})
	}

	host := c.FormValue("host")
	target := c.FormValue("target")

	_, exists := pm.redirectManager.getRedirect(host)

	if exists {
		fmt.Println("Already exists")
		return c.JSON(http.StatusBadRequest, &APIError{
			Response: &APIResponse{
				Message: "Already used host",
			},
			ErrorCode: 2,
		})
	}

	pm.redirectManager.addNewRedirect(host, target)
	u := &APIResponse{
		Message: "OK",
	}
	return c.JSON(http.StatusOK, u)
}

func (pm *ProxyManagerAPI) removeProxies(c echo.Context) error {
	errMsg := ""
	errMsg += checkQueryParam(c, "host")

	if errMsg != "" {
		return c.JSON(http.StatusBadRequest, &APIError{
			Response: &APIResponse{
				Message: "Validation error",
			},
			Details:   errMsg,
			ErrorCode: 1,
		})
	}

	host := c.QueryParam("host")

	_, exists := pm.redirectManager.getRedirect(host)

	if !exists {
		return c.JSON(http.StatusBadRequest, &APIError{
			Response: &APIResponse{
				Message: "Not exists host",
			},
			ErrorCode: 2,
		})
	}

	pm.redirectManager.removeRedirect(host)

	u := &APIResponse{
		Message: "OK",
	}
	return c.JSON(http.StatusOK, u)
}
