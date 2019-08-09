package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/g-airport/tool-box/email/entity"
)

const (
	
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

func Luminati() {
	cli, err := NewHttpProxyClient(_ProxyHttpHost).Get("http://lumtest.com/myip.json")
	if err != nil {
		log.Print("Luminati client err : ", err)
	}
	if cli.StatusCode != http.StatusOK {
		log.Print("Luminati client status : ", cli.StatusCode)
	}
	defer cli.Body.Close()
	buf, err := ioutil.ReadAll(cli.Body)
	if err != nil {
		log.Print("Luminati read body err : ", err)
	}
	ip := &struct {
		IP string `json:"ip"`
	}{}
	err = json.Unmarshal(buf, ip)
	if err != nil {
		log.Print("Luminati Unmarshal err", err)
	}
	log.Print("Luminati", *ip)
}

func EmailProxyClientAPI(email string) *entity.EmailInfo {
start:
	apiUrl := fmt.Sprintf("https://api.trumail.io/v2/lookups/json?email=%s", email)
	cli, err := NewHttpProxyClient(_ProxyHttpHost).Get(apiUrl)
	if err != nil {
		log.Print("EmailProxyClientAPI client err : ", err)
		goto start
	}
	if cli == nil {
		goto start
	}
	if cli.StatusCode != http.StatusOK {
		log.Print("EmailProxyClientAPI client status : ", cli.StatusCode)
	}
	defer cli.Body.Close()
	buf, err := ioutil.ReadAll(cli.Body)
	if err != nil {
		log.Print("EmailProxyClientAPI read body err : ", err)
	}
	e := &entity.EmailInfo{}
	err = json.Unmarshal(buf, e)
	if err != nil {
		log.Print("EmailProxyClientAPI Unmarshal err", err)
	}
	if e.Email == "" {
		goto start
	}
	log.Print("EmailProxyClientAPI", e)
	return e
}

func EmailDirectClientAPI(email string) *entity.EmailInfo {
start:
	apiUrl := fmt.Sprintf("https://api.trumail.io/v2/lookups/json?email=%s", email)
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Print("EmailDirectClientAPI request err : ", err)
	}
	cli := &http.Client{}
	//cli.Timeout = 5 * time.Second
	res, err := cli.Do(request)
	if err != nil {
		log.Print("EmailDirectClientAPI do request err : ", err)
	}
	if res == nil {
		log.Print("EmailDirectClientAPI res nil", email)
		return &entity.EmailInfo{}
	}
	var getErr = err
	if res.StatusCode != http.StatusOK {
		log.Print("EmailDirectClientAPI client status : ", res.StatusCode)
	}
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("EmailDirectClientAPI read body err : ", err)
	}
	e := &entity.EmailInfo{}
	err = json.Unmarshal(buf, e)
	e.Err = getErr
	if err != nil {
		log.Print("EmailDirectClientAPI Unmarshal err", err)
	}
	if e.Email == "" {
		goto start
	}
	log.Print("EmailDirectClientAPI", *e)
	return e
}
