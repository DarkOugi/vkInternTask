package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vk/internal/db"
	"vk/internal/server"
	"vk/internal/service"

	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

//nolint:gochecknoglobals // тише тише тише
var (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUsername = "vk"
	dbPassword = "0000"
	dbName     = "vkdb"

	serverPort = "8080"
)

func main() {
	var err error
	var pSQL *db.PostgresDB

	ctx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,
	}).Level(zerolog.DebugLevel)

	flag.StringVar(&dbHost, "dbHost", dbHost, "dbHost pgx connect")
	flag.StringVar(&dbPort, "dbPort", dbPort, "dbPort pgx connect")
	flag.StringVar(&dbUsername, "dbUsername", dbUsername, "dbUsername pgx connect")
	flag.StringVar(&dbPassword, "dbPassword", dbPassword, "dbPassword pgx connect")
	flag.StringVar(&dbName, "dbName", dbName, "dbName pgx connect")

	flag.StringVar(&serverPort, "serverPort", serverPort, "server run in this port")
	flag.Parse()

	if pSQL == nil {
		pSQL, err = db.NewPostgresDB(ctx, dbHost, dbPort, dbUsername, dbPassword, dbName)
		if err != nil {
			log.Error().Err(err).Msg("don't create connect with db")
			return
		}
	}
	defer func() {
		pSQL.Close()
	}()

	sv := service.NewService(pSQL)
	sr := server.NewServer(sv)

	r := router.New()
	r.POST("/api/auth", sr.Auth)
	r.POST("/api/auth", sr.Reg)
	r.POST("/api/adv", sr.CreateAdv)
	r.GET("/api/advs", sr.GetAllAdv)

	go func() {
		if errServer := fasthttp.ListenAndServe(fmt.Sprintf(":%s", serverPort), r.Handler); errServer != nil {
			log.Fatal().Err(errServer).Msg("server critical error")
		}
	}()
	log.Info().Msg("SERVER SUCCESS START")
	<-ctx.Done()
}
