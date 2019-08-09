package http_client

const (
	_ProxyHttpHost = ""
	_ProxyHttpUsername = ""
	_ProxyHttpPassword = ""
)

type BasicAuth struct {
	_ProxyHttpUsername string
	_ProxyHttpPassword string
}

func NewBasicAuth(username, password string) BasicAuth {
	return BasicAuth{
		_ProxyHttpUsername: username,
		_ProxyHttpPassword: password,
	}
}
