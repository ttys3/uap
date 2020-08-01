package ua_proxy

const DftPasswd = "passw0rd"

const RepoURL = "https://github.com/ttys3/ua-proxy"

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
