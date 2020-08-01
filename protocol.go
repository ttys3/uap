package ua_proxy

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
