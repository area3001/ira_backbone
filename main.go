package main

import (
	"github.com/area3001/goira/app"
	"github.com/area3001/goira/comm"
	_ "github.com/area3001/goira/docs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"os"
)

// @title Fri3d IRA API
// @version 1.0
// @description This is the Fri3d IRA ReST API
// @termsOfService http://swagger.io/terms/

// @contact.name code@fri3d.be
// @contact.url http://www.swagger.io/support
// @contact.email code@fri3d.be

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1323
// @BasePath /
// @schemes http
func main() {
	loadEnv()

	opts := &comm.NatsClientOpts{
		Root:             "area3001",
		NatsUrl:          os.Getenv("GOIRA_NATS_URL"),
		NatsOptions:      []nats.Option{},
		JetStreamOptions: []nats.JSOpt{},
	}

	nc, err := comm.Dial(opts)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		_ = nc.Close()
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	a := app.NewApp(nc, e)

	if err := a.Devices.Service.Ping(); err != nil {
		log.Panicln(err)
	}

	e.Logger.Fatal(e.Start(":1323"))
}

func loadEnv() {
	env := os.Getenv("GOIRA_ENV")
	if "" == env {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")
	if "test" != env {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load() // The Original .env
}
