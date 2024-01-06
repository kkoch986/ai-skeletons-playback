package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hypebeast/go-osc/osc"
	"github.com/juju/zaputil/zapctx"
	"github.com/kkoch986/ai-skeletons-playback/command/internal/common"
	"github.com/kkoch986/ai-skeletons-playback/output"
	"github.com/kkoch986/ai-skeletons-playback/server"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var Command = &cli.Command{
	Name:  "server",
	Usage: "start the http server to handle output generation requests",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Usage:   "the hostname to listen on, default to an empty string to listen on all interfaces",
			Value:   "",
			EnvVars: []string{"HOST"},
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "the port for the server to listen on",
			EnvVars: []string{"PORT"},
			Value:   3000,
		},
		&cli.StringFlag{
			Name:    "osc-host",
			Usage:   "the hostname to send OSC messages to",
			Value:   "localhost",
			EnvVars: []string{"OSC_HOST"},
		},
		&cli.IntFlag{
			Name:    "osc-port",
			Usage:   "the port to send OSC messages on",
			EnvVars: []string{"OSC_PORT"},
			Value:   7771,
		},
	},
	Action: func(c *cli.Context) error {
		ctx, cancel := common.InitContext()
		defer cancel()

		serverHost := c.String("host")
		serverPort := c.Int("port")
		oscHost := c.String("osc-host")
		oscPort := c.Int("osc-port")

		ctx, flush, err := common.Logger(ctx, c.String("log-level"))
		if err != nil {
			return err
		}
		defer flush()
		logger := zapctx.Logger(ctx)

		logger.Info("starting web server", zap.String("host", serverHost), zap.Int("port", serverPort))
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		r.Use(ginzap.RecoveryWithZap(logger, true))

		oscClient := osc.NewClient(oscHost, oscPort)

		w := output.NewMultiplexer("main", []output.HardwareWriter{
			&output.LogWriter{},
			output.NewOSCWriter("test-osc", oscClient),
		})

		// The playback action
		r.POST("/playback", func(c *gin.Context) {
			r := &server.PlaybackRequest{
				SequenceData: output.Sequence{},
			}
			err := c.BindJSON(r)
			if err != nil {
				c.JSON(http.StatusBadRequest, &ErrorResponse{err.Error()})
				return
			}
			resp, err := server.HandlePOSTPlayback(ctx, r, w)
			if err != nil {
				c.JSON(http.StatusInternalServerError, &ErrorResponse{err.Error()})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		})

		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", serverHost, serverPort),
			Handler: r,
		}

		go server.ListenAndServe()

		<-ctx.Done()
		logger.Info("context cancelled, shutting server down")
		err = server.Shutdown(context.Background())
		if err != nil {
			logger.Error("error shutting down server", zap.Error(err))
		}
		return err
	},
}
