package main

import (
	"github.com/area3001/goira"
	"github.com/area3001/goira/comm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
)

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

	app := goira.NewApp(nc)

	if err := app.Devices.Service.Ping(); err != nil {
		log.Panicln(err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

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
