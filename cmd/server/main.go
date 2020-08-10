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

	"github.com/ttys3/uap"
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
		Flags: []cli.Flag{
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
				Value:   uap.DftPasswd,
				Usage:   "auth password",
				EnvVars: []string{"UAP_AUTH"},
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {

	log.Println(ctx.App.Version)

	// xdg-open will cal DE open command like gio, and gio may failed with: Error: no DISPLAY environment variable specified
	// but xdg-open will result with exit code 0
	// here we exit and leave systemd restart the app so that we can ensure the DISPLAY environment is set
	if dp := os.Getenv("DISPLAY"); dp == "" {
		// https://wiki.archlinux.org/index.php/Systemd/User#DISPLAY_and_XAUTHORITY
		// https://github.com/systemd/systemd/blob/v229/xorg/50-systemd-user.sh
		//exec.Command("systemctl", "--user", "import-environment", "DISPLAY", "XAUTHORITY")
		//exec.Command("dbus-update-activation-environment", "DISPLAY", "XAUTHORITY")
		// exit and leave systemd restart the app
		log.Fatalf("DISPLAY environment is not set, exit ...")
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/open", func(c *gin.Context) {
		var req uap.UaProxyReq
		if err := c.BindJSON(&req); err != nil {
			log.Printf("invalid req param, req=%+v err=%s", req, err.Error())
			c.JSON(http.StatusBadRequest, uap.UaProxyRsp{RetCode: uap.RetCodeInvalidReq, Msg: err.Error()})
			return
		}

		if err := req.ValidateURL(); err != nil {
			log.Printf("url validate failed, url=%s err=%s", req.Url, err.Error())
			c.JSON(http.StatusBadRequest, uap.UaProxyRsp{RetCode: uap.RetCodeInvalidURL, Msg: err.Error()})
			return
		}

		if req.Auth != ctx.String("auth") {
			c.JSON(http.StatusBadRequest, uap.UaProxyRsp{RetCode: uap.RetCodeAuthFailed, Msg: "auth failed"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		cmd := exec.CommandContext(ctx, "xdg-open", req.Url)
		out, err := cmd.CombinedOutput()

		log.Printf("exec %+v, err=%v, out=%v ", cmd, err, string(out))
		if err != nil {
			c.JSON(http.StatusInternalServerError, uap.UaProxyRsp{
				RetCode: uap.RetCodeExecFailed,
				Msg:     fmt.Sprintf("err: %s, out=%s", err, string(out)),
			})
			return
		}
		c.JSON(http.StatusOK, uap.UaProxyRsp{RetCode: uap.RetCodeOK})
	})

	withPwdEn := ""
	if ctx.String("auth") != "" {
		withPwdEn = "with auth enabled"
	}
	log.Printf("http server listen on %s %s", ctx.String("addr"), withPwdEn)
	return r.Run(ctx.String("addr"))
}
