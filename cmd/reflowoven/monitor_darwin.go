//go:build darwin

package main

import (
	"context"
	"github.com/minor-industries/rtgraph"
	"sync"
	"time"
)

func monitorTemp(
	ctx context.Context,
	gr *rtgraph.Graph,
	wg *sync.WaitGroup,
	t0 time.Time,
	errCh chan error,
	profile Schedule,
) {
	ticker := time.NewTicker(250 * time.Millisecond)

	for t := range ticker.C {
		_ = t
		err := gr.CreateValue("reflowoven_temperature", t, 25.0)
		if err != nil {
			errCh <- err
			return
		}
	}
}
