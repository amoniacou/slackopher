package main

import (
	"github.com/amoniacou/slackopher/router"
	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v2"
	"io"
	"log"
	"os"
	s "strings"
)

var webCmd = cli.Command{
	Name:    "web",
	Usage:   "starts slaclopher web server daemon",
	Aliases: []string{"s"},
	Action: func(c *cli.Context) {
		if err := server(c); err != nil {
			log.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			EnvVar: "SERVER_ADDR",
			Name:   "bind, b",
			Usage:  "URI to bind to (tcp://, unix://, ssl://)",
		},
		cli.StringFlag{
			EnvVar: "PORT",
			Name:   "port, p",
			Value:  "3030",
			Usage:  "Define the TCP port to bind to (default to 3030)\nUse -b for more advanced options",
		},
	},
}

func server(c *cli.Context) error {
	f, err := os.OpenFile("/var/log/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	bind := c.String("b")
	if bind == "" {
		bind = ":" + c.String("p")
	}

	logWriter := io.MultiWriter(f, os.Stdout)
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = logWriter
	log.SetOutput(logWriter)
	//}
	r := router.WebServer()
	switch {
	case s.Contains(bind, "ssl://"):
		panic("Not implemented yet")
	case s.Contains(bind, "tcp://"):
		return r.Run(s.Split(bind, "tcp://")[0])
	case s.Contains(bind, "unix://"):
		return r.RunUnix(s.Split(bind, "unix://")[0])
	default:
		return r.Run(bind)
	}
}
