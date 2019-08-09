package client

import "testing"

func TestProxyEmailClientAPI(t *testing.T) {
	EmailProxyClientAPI("insomnus@lovec.at")
}

func TestLuminati(t *testing.T) {
	Luminati()
}

func TestEmailDirectClientAPI(t *testing.T) {
	EmailDirectClientAPI("insomnus@lovec.at")
}