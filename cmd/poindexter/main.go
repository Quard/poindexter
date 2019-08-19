package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jessevdk/go-flags"

	"github.com/Quard/poindexter/internal/internal_api"
	"github.com/Quard/poindexter/internal/storage"
)

var opts struct {
	Bind      string `short:"b" long:"bind" env:"BIND" default:"localhost:5010"`
	MongoURI  string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb://localhost:27017/poindexter"`
	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN"`
}

func main() {
	parser := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	sentry.Init(sentry.ClientOptions{
		Dsn: opts.SentryDSN,
	})
	defer sentry.Flush(time.Second * 5)

	stor := storage.NewMongoStorage(opts.MongoURI)
	srv := internal_api.NewInternalAPIServer(opts.Bind, stor)
	srv.Run()
}
