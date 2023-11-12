package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/rizalarfiyan/be-revend/docs"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
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

// @title BE Revend API
// @version 1.0
// @termsOfService http://swagger.io/terms/
// @contact.name Rizal Arfiyan
// @contact.url https://rizalrfiyan.com
// @contact.email rizal.arfiyan.23@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @description This is a API documentation of BE Revend
// @BasePath /
// @securityDefinitions.apikey AccessToken
// @in header
// @name Authorization
func main() {
	conf := config.Get()
	app := fiber.New(config.FiberConfig())
	logs := logger.Get("main")
	logApi := logger.Get("api").Logger()
	app.Use(fiberzerolog.New(config.FiberZerolog(logApi)))
	app.Use(requestid.New())
	app.Use(cors.New(config.CorsConfig()))
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(recover.New())

	app.Get("/swagger/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			conf.Swagger.Username: conf.Swagger.Password,
		},
	}), swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  true,
		DocExpansion: "list",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return response.New(c, http.StatusOK, "Hello World!", nil)
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		if true {
			panic(response.NewErrorMessage(http.StatusBadRequest, "invalid input", nil))
		}

		return response.New(c, http.StatusOK, "Success", nil)
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
