package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/urfave/cli/v2"

	ua_proxy "github.com/ttys3/ua-proxy"
)

var Version = "dev"
var CommitSHA = "dev"
var BuildDate = "unkown"

const appName = "uap-server"

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"v"},
		Usage: "print only the version",
	}

	app := &cli.App{
		Version: fmt.Sprintf("%s %s %s %s", appName, Version, CommitSHA, BuildDate),
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

	log.Println(ctx.App.Version)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/open", func(c *gin.Context) {
		var req ua_proxy.UaProxyReq
		c.BindJSON(&req)

		if err := req.ValidateURL(); err != nil {
			c.JSON(http.StatusBadRequest, ua_proxy.UaProxyRsp{RetCode: 1, Msg: err.Error()})
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
			c.JSON(http.StatusInternalServerError, ua_proxy.UaProxyRsp{
				RetCode: 2,
				Msg: fmt.Sprintf("err: %s, out=%s", err, string(out)),
			})
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
