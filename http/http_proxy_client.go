package http_client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"
)

func NewHttpProxyClient(proxyAddr string) *http.Client {
	proxyUrl := &url.URL{
		Scheme: "http",
		User:   url.UserPassword(_ProxyHttpUsername, _ProxyHttpPassword),
		Host:   proxyAddr,
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		DialContext: func(ctx context.Context, nwk, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(nwk, addr, time.Second*time.Duration(100))
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}
