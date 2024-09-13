package main

import (
	"sync"
)

type RedirectManager struct {
	proxyMap *sync.Map
}

func createNewRedirectManager() *RedirectManager {
	manager := new(RedirectManager)

	manager.proxyMap = new(sync.Map)

	return manager
}

func (rm *RedirectManager) addNewRedirect(host string, target string) (ok bool) {
	_, exists := rm.proxyMap.Load(host)

	if exists {
		return false
	}

	rm.proxyMap.Store(host, target)
	return true
}

func (rm *RedirectManager) getRedirect(host string) (value string, ok bool) {
	target, exists := rm.proxyMap.Load(host)

	if !exists {
		return "", false
	}

	return target.(string), exists
}

func (rm *RedirectManager) removeRedirect(host string) (ok bool) {
	rm.proxyMap.Delete(host)

	return true
}
