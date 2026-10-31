package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"qixue-api/internal/platform/applog"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/platform/config"
	"qixue-api/internal/platform/db"
	webapprepo "qixue-api/internal/webapp/repo"
	"qixue-api/internal/server"
	"qixue-api/internal/service/captcha"
	"qixue-api/internal/webapp/service/rank"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logWriter, err := applog.Setup(cfg.LogDir, cfg.LogKeepDays)
	if err != nil {
		log.Fatal(err)
	}
	defer logWriter.Close()

	conn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	rdb, err := cache.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slide, err := captcha.New(rdb)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	rank.New(webapprepo.NewRank(conn), rdb).StartTicker(ctx)

	log.Printf("qixue-api listening on %s log=%s keep_days=%d", cfg.Addr, logWriter.Path(), cfg.LogKeepDays)
	if err := server.New(cfg, conn, rdb, slide).Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
