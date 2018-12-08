package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/ASeegull/secrets-vault/api"
	"github.com/ASeegull/secrets-vault/env"
	"github.com/ASeegull/secrets-vault/storage"
	"github.com/ASeegull/secrets-vault/storage/migrations"
)

func main() {
	var (
		err  error
		once sync.Once
		st   storage.DB
		cfg  env.Config
	)

	envconfig.MustProcess("", &cfg)

	once.Do(func() {
		st, err = storage.InitDB(
			cfg.DB_HOST, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_SSLMODE,
		)
	})

	if err != nil {
		log.Fatal("DB is unavailable: ", err)
	}

	if err = migrations.Run(st.Conn.DB); err != nil {
		log.Fatal("failed to setup schema: ", err)
	}

	srv := api.New(st, cfg.HOST, cfg.PORT)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(errors.Wrap(err, "unexpected server shutdown"))
			os.Exit(1)
		}
	}()

	log.WithFields(log.Fields{
		"Addr": srv.Addr,
	}).Info("Server is listening")

	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	<-sign

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Warn(errors.Wrap(err, "some oppened connections were interrupted"))
	}

	log.Info("Application has stopped")
}
