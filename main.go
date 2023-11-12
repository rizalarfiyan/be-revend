package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/logger"
)

func init() {
	config.Init()
	conf := config.Get()
	logger.Init(conf)
	logger.UpdateLogLevel(conf.Logger.Level)
}

func main() {
	conf := config.Get()
	app := fiber.New(config.FiberConfig())
	logs := logger.Get("main")
	logApi := logger.Get("api").Logger()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logApi,
		Fields: []string{"ip", "method", "path", "url", "latency", "status", "body", "error", "requestId"},
		SkipBody: func(ctx *fiber.Ctx) bool {
			return strings.Contains(string(ctx.Request().Header.ContentType()), "multipart/form-data")
		},
	}))
	app.Use(requestid.New())
	app.Use(cors.New(config.CorsConfig()))
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return response.New(c, http.StatusOK, "Hello World!", nil)
	})

	baseUrl := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server := &http.Server{
		Addr: baseUrl,
	}

	go func() {
		err := app.Listen(baseUrl)
		if err != nil {
			logs.Fatal(err, "Error app serve")
		}
	}()

	handleShutdown(server, logs)
}

func handleShutdown(server *http.Server, logs logger.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logs.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if err = server.Shutdown(ctx); err != nil {
		logs.Fatal(err, "Server forced to shutdown")
	}

	logs.Info("Server exiting")
}
