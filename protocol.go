package ua_proxy

import (
	"errors"
	"net/url"
)

const DftPasswd = "passw0rd"

const RepoURL = "https://github.com/ttys3/ua-proxy"

var ErrUnsupportedProto = errors.New("unsupported protocol")

const (
	RetCodeOK = iota
	RetCodeInvalidReq
	RetCodeInvalidURL
	RetCodeAuthFailed
	RetCodeExecFailed
)

var supportedProtos = map[string]bool{
	"http":   true,
	"https":  true,
	"ftp":    true,
	"irc":    true,
	"mailto": true,
	"mms":    true,
	"news":   true,
	"nntp":   true,
	"sms":    true,
	"smsto":  true,
	"snews":  true,
	"tel":    true,
	"urn":    true,
	"webcal": true,
}

type UaProxyReq struct {
	Auth        string `json: "auth"`
	FromMachine string `json: "from_machine"`
	Url         string `json: "url"`
	ReqTs       int64  `json: "req_ts"`
}

type UaProxyRsp struct {
	RetCode int    `json: "ret_code"`
	Msg     string `json: "msg"`
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
