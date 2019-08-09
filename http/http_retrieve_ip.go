package http_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

// Proxy Config
var (
	_proxyAddr     = ""
	_proxyUserName = ""
	_proxyPassword = ""
)

// Access Site
var (
	Url = ""
)

//example
//https://api.ipify.org/
//http://lumtest.com/myip.json

// ----------------------------

func RetrieveV1(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to retrieve public IP")
	}
	defer res.Body.Close()
	buf, _ := ioutil.ReadAll(res.Body)
	names, err := net.LookupAddr(string(buf))
	if err != nil {
		return string(buf)
	}
	return strings.TrimSuffix(names[0], ".")
}

func RetrieveV2(url string) string {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("err request", err)
	}
	request.Header.Add("Upgrade-Insecure-Requests", "1")
	request.Header.Add("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	res, err := NewHttpProxyClient(_proxyAddr).Do(request)
	if err != nil {
		log.Fatal("Failed to retrieve public IP :", err)
	}
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("RetrieveV2 read body err", err)
	}
	//avoid memory leak
	//buf := bytes.NewBuffer(make([]byte, 1024))
	//_, err = io.Copy(buf, res.Body)
	//if err != nil {
	//log.Fatal("RetrieveV2 read body err", err)
	//}
	return string(buf)
}

func RetrieveV3() {
	//creating the proxyURL
	proxyURL, err := url.Parse(_proxyAddr)
	if err != nil {
		log.Println(err)
	}

	//creating the URL to be loaded through the proxy
	url2, err := url.Parse(Url)
	if err != nil {
		log.Fatal(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", url2.String(), nil)
	if err != nil {
		log.Println(err)
	}

	//adding proxy authentication
	auth := fmt.Sprintf("%s:%s", _proxyUserName, _proxyPassword)
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	request.Header.Add("Proxy-Authorization", basicAuth)
	//request.SetBasicAuth(proxyUser,ProxyPwd)

	//printing the request to the console
	dump, _ := httputil.DumpRequest(request, false)
	log.Println(string(dump))

	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response.StatusCode)
	log.Println(response.Status)
	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	log.Println(string(data))

}

func InitSock5(socks5Url string) {
	var (
		auth    *proxy.Auth
		pUrl, _ = url.Parse(socks5Url)
	)

	if pUrl.User != nil {
		auth = &proxy.Auth{}
		auth.User = pUrl.User.Username()
		auth.Password, _ = pUrl.User.Password()
	}

	dialer, err := proxy.SOCKS5("tcp", pUrl.Host, auth, proxy.Direct)
	if err != nil {
		log.Fatalf("socks5 proxy err: %v\n", err)
	}
	http.DefaultClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}}
}
