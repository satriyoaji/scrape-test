package scraper

import (
	"math/rand"
	"strings"
)

type ProxyManager struct {
	proxies []string
}

func NewProxyManager(proxyList string) *ProxyManager {
	return &ProxyManager{
		proxies: strings.Split(proxyList, ","),
	}
}

func (pm *ProxyManager) GetRandomProxy() string {
	return pm.proxies[rand.Intn(len(pm.proxies))]
}
