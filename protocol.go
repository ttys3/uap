package ua_proxy

import (
	"errors"
	"net/url"
)

const DftPasswd = "passw0rd"

const RepoURL = "https://github.com/ttys3/ua-proxy"

var ErrUnsupportedProto = errors.New("unsupported protocol")

var supportedProtos = map[string]bool{
	"http": true,
	"https": true,
	"ftp": true,
	"irc": true,
	"mailto": true,
	"mms": true,
	"news": true,
	"nntp": true,
	"sms": true,
	"smsto": true,
	"snews": true,
	"tel": true,
	"urn": true,
	"webcal": true,
}

type UaProxyReq struct {
	Auth        string
	FromMachine string
	Url         string
	ReqTs       int64
}

type UaProxyRsp struct {
	RetCode int
	Msg     string
}

func (r *UaProxyReq) ValidateURL() error {
	u, err := url.Parse(r.Url)
	if err != nil {
		return err
	}
	if _, ok := supportedProtos[u.Scheme]; !ok {
		return ErrUnsupportedProto
	}
	return nil
}