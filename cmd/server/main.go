package main

import (
	"github.com/gin-gonic/gin"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	ua_proxy "github.com/ttys3/ua-proxy"
)

var Version = "dev"
var CommitSHA = "dev"
var BuildDate = "unkown"

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"v"},
		Usage: "print only the version",
	}

	app := &cli.App{
		Version: "1.0.0",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Value:   ":18080",
				Usage:   "listen addr",
				EnvVars: []string{"UAP_ADDR"},
			},
			&cli.StringFlag{
				Name:    "auth",
				Aliases: []string{"p"},
				Value:   ua_proxy.DftPasswd,
				Usage:   "auth password",
				EnvVars: []string{"UAP_AUTH"},
			},
		},
		Action:run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/open", func(c *gin.Context) {
		var req ua_proxy.UaProxyReq
		c.BindJSON(&req)
		if req.Url == "" ||
			(!strings.HasPrefix(req.Url, "ftp") &&
				!strings.HasPrefix(req.Url, "http") &&
				!strings.HasPrefix(req.Url, "https")) {
			c.JSON(http.StatusBadRequest, ua_proxy.UaProxyRsp{RetCode: 1, Msg: "invalid protocol"})
			return
		}

		if req.Auth != ctx.String("auth") {
			c.JSON(http.StatusBadRequest, ua_proxy.UaProxyRsp{RetCode: 1, Msg: "auth failed"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
		defer cancel()
		cmd := exec.CommandContext(ctx, "xdg-open", req.Url)
		out, err := cmd.CombinedOutput()

		if err != nil {
			c.JSON(http.StatusInternalServerError, ua_proxy.UaProxyRsp{RetCode: 2, Msg: string(out)})
			return
		}
		c.JSON(http.StatusOK, ua_proxy.UaProxyRsp{RetCode: 0})
	})

	withPwdEn := ""
	if ctx.String("auth") != "" {
		withPwdEn = "with auth enabled"
	}
	log.Printf("http server listen on %s %s", ctx.String("addr"), withPwdEn)
	return r.Run(ctx.String("addr"))
}