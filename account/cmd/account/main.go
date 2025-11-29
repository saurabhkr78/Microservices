package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/saurabh/Microservices/account"
	"github.com/tinrab/retry"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	// retry until we get a repo (the callback signature func(int) error)
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return // return err (nil => stop retrying)
	})
	// now repo is non-nil and connected — defer close here
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	// run gRPC server in goroutine so we can shutdown on signal
	errCh := make(chan error, 1)
	go func() {
		if err := account.ListenGRPC(s, 8080); err != nil {
			errCh <- err
		}
	}()

	// graceful shutdown on SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %v, shutting down", sig)
		// if ListenGRPC returns a cancellable server, call its Stop/GracefulStop here
		// repo.Close() will run from deferred call
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}
}
