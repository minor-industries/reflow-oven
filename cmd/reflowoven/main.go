package main

import (
	"context"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"os"
	"os/signal"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"strings"
	"sync"
	"syscall"
	"time"
)

//var profile = NewScheduleRelativeDurations([]Point{
//	{Duration(0 * time.Second), 25},
//	//{Duration(40 * time.Second), 45}, // preheat the element
//	{Duration(30 * time.Second), 100},
//	{Duration(90 * time.Second), 150},
//	{Duration(30 * time.Second), 183},
//	{Duration(60 * time.Second), 235},
//	{Duration(30 * time.Second), 235},
//	{Duration(25 * time.Second), 183},
//	{Duration(60 * time.Second), 25},
//	{Duration(30 * time.Second), 25},
//})

var profile = NewScheduleRelativeDurations([]Point{
	//{Duration(40 * time.Second), 45}, // preheat the element
	{Duration(0 * time.Second), 238},
	{Duration(8 * 60 * time.Second), 238},
	{Duration(30 * time.Second), 25},
})

func main() {
	for _, p := range profile {
		fmt.Println(p.T().Seconds(), p.Val)
	}

	t0 := time.Now()
	fmt.Println(t0)

	_, err := host.Init()
	noErr(err)

	bus, err := i2creg.Open("1")
	noErr(err)

	log := func(s string) {
		fmt.Println(s)
	}

	tcs := []*Thermocouple{
		NewThermocouple(log, bus, 0x67, "probe0"),
		NewThermocouple(log, bus, 0x60, "probe1"),
	}

	cook := gpioreg.ByName("GPIO16")
	if cook == nil {
		panic("no gpio")
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go monitorTemp(ctx, &wg, t0, tcs, cook)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signals
		cancel()
	}()

	wg.Wait()
}

func monitorTemp(
	ctx context.Context,
	wg *sync.WaitGroup,
	t0 time.Time,
	tcs []*Thermocouple,
	cook gpio.PinIO,
) {
	ticker := time.NewTicker(250 * time.Millisecond)
	data := make([][]Point, len(tcs))

	for {
		select {
		case <-ticker.C:
			t1 := time.Now()
			t := t1.Sub(t0)

			target := profile.Val(t)
			parts := []string{
				fmt.Sprintf("t=%3.2f", t.Seconds()),
				fmt.Sprintf("target=%3.2f", target),
			}

			for i, tc := range tcs {
				temp, err := tc.Temperature()
				noErr(err)

				data[i] = append(data[i], Point{
					Time: Duration(t),
					Val:  temp,
				})

				parts = append(parts, fmt.Sprintf("%s=%3.2f", tc.Description, temp))
			}

			// use the value of the first probe for control
			temp := data[0][len(data[0])-1].Val
			on := target > temp

			parts = append(parts, fmt.Sprintf("on=%v", on))
			fmt.Println(strings.Join(parts, " "))

			err := cook.Out(gpio.Level(on))
			noErr(err)
		case <-ctx.Done():
			ticker.Stop()
			_ = cook.Out(gpio.Low)
			fmt.Println("done")
			graph(profile, tcs, data)
			wg.Done()
			return
		}
	}
}

func graph(
	schedule Schedule,
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
		for _, d := range schedule {
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
