package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minor-industries/reflow-oven/html"
	"github.com/minor-industries/reflow-oven/reflow"
	"github.com/minor-industries/rtgraph"
	"github.com/minor-industries/rtgraph/database/inmem"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func run() error {
	gin.SetMode(gin.ReleaseMode)

	profile := reflow.Profile1

	t0 := time.Now()

	errCh := make(chan error)
	be := &reflow.Backend{
		T0:            t0,
		NormalBackend: inmem.NewBackend(),
		Profile:       profile,
	}

	gr, err := rtgraph.New(
		be,
		errCh,
		rtgraph.Opts{},
		[]string{
			"reflowoven_temperature",
		},
	)
	if err != nil {
		return errors.Wrap(err, "new rtgraph")
	}

	gr.StaticFiles(html.FS,
		"index.html", "text/html",
	)

	go func() {
		errCh <- gr.RunServer("0.0.0.0:8080")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)

	go reflow.MonitorTemp(ctx, gr, &wg, t0, errCh, profile)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			<-signals
			cancel()
		}()

		wg.Wait()
		errCh <- nil
	}()

	return <-errCh
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
