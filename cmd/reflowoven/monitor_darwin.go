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

	for {
		select {
		case t := <-ticker.C:
			err := gr.CreateValue("reflowoven_temperature", t, 25.0)
			if err != nil {
				errCh <- err
				return
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}
