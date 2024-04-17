package main

import (
	"context"
	"github.com/minor-industries/codelab/cmd/reflowoven/html"
	"github.com/minor-industries/rtgraph"
	"github.com/minor-industries/rtgraph/database/inmem"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var profile2 = NewScheduleRelativeDurations([]Point{
	{Duration(0 * time.Second), 45},
	{Duration(40 * time.Second), 45}, // preheat the element

	{Duration(0 * time.Second), 25},

	{Duration(30 * time.Second), 100},
	{Duration(90 * time.Second), 150},
	{Duration(30 * time.Second), 183},
	{Duration(60 * time.Second), 235},
	{Duration(5 * time.Second), 235},
	{Duration(25 * time.Second), 183},
	{Duration(30 * time.Second), 25},
	{Duration(30 * time.Second), 25},
})

var profile1 = NewScheduleRelativeDurations([]Point{
	{Duration(0 * time.Second), 60},
	{Duration(40 * time.Second), 60}, // preheat the element

	{Duration(25 * time.Second), 100},
	{Duration(90 * time.Second), 150},
	{Duration(30 * time.Second), 183},
	{Duration(60 * time.Second), 235},
	{Duration(5 * time.Second), 235},
	{Duration(25 * time.Second), 183},
	{Duration(30 * time.Second), 25},
	{Duration(30 * time.Second), 25},
})

func main() {
	profile := profile1

	t0 := time.Now()

	errCh := make(chan error)
	be := &backend{
		t0:            t0,
		normalBackend: inmem.NewBackend(),
		profile:       profile,
	}
	gr, err := rtgraph.New(be, errCh, rtgraph.Opts{}, []string{"reflowoven_temperature"})
	noErr(err)

	gr.StaticFiles(html.FS,
		"index.html", "text/html",
	)

	go func() {
		noErr(<-errCh)
	}()

	go func() {
		errCh <- gr.RunServer("0.0.0.0:8080")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go monitorTemp(ctx, gr, &wg, t0, errCh, profile)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signals
		cancel()
	}()

	wg.Wait()
}

func graph(
	profile Schedule,
	tcs []*Thermocouple,
	data [][]Point,
) {
	p := plot.New()

	p.Title.Text = "Reflow Oven Temperature"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Temp"

	var vs []interface{}
	{
		vs = append(vs, "profile")
		var pts plotter.XYs
		for _, d := range profile {
			pts = append(pts, plotter.XY{
				d.T().Seconds(),
				d.Val,
			})
		}
		vs = append(vs, pts)
	}

	for i, tc := range tcs {
		vs = append(vs, tc.Description)
		pts := plotter.XYs{}
		for _, d := range data[i] {
			pts = append(pts, plotter.XY{
				d.T().Seconds(),
				d.Val,
			})
		}
		vs = append(vs, pts)
	}

	err := plotutil.AddLines(p, vs...)
	noErr(err)

	err = p.Save(16*vg.Inch, 4*vg.Inch, "cook.svg")
	noErr(err)
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
