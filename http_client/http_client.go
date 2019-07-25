package http_client

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"
)

func NewHttpClient(proxyAddr string) *http.Client {
	//direct
	//proxy2, err := url.Parse(proxyAddr)
	//if err != nil {
	//	return nil
	//}
	//if proxy2.User != nil {
	//	auth := &proxy.Auth{}
	//	auth.User = proxy2.User.Username()
	//	auth.Password, _ = proxy2.User.Password()
	//}

	proxyUrl := &url.URL{
		Scheme: "http",
		User:   url.UserPassword(_proxyUserName, _proxyPassword),
		Host:   proxyAddr,
	}
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(100))
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}
