//go:build linux

package reflow

import (
	"context"
	"fmt"
	"github.com/minor-industries/rtgraph"
	"github.com/pkg/errors"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"strings"
	"sync"
	"time"
)

func MonitorTemp(
	ctx context.Context,
	gr *rtgraph.Graph,
	wg *sync.WaitGroup,
	t0 time.Time,
	errCh chan error,
	profile Schedule,
) {
	_, err := host.Init()
	if err != nil {
		errCh <- errors.Wrap(err, "init host")
		return
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		errCh <- errors.Wrap(err, "open i2c bus")
		return
	}

	log := func(s string) {
		fmt.Println(s)
	}

	tcs := []*Thermocouple{
		NewThermocouple(log, bus, 0x67, "probe0"),
		//NewThermocouple(log, bus, 0x60, "probe1"),
	}

	cook := gpioreg.ByName("GPIO16")
	if cook == nil {
		panic("no gpio")
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	data := make([][]Point, len(tcs))

	for {
		select {
		case t1 := <-ticker.C:
			t := t1.Sub(t0)

			target := profile.Val(t)
			parts := []string{
				fmt.Sprintf("t=%3.2f", t.Seconds()),
				fmt.Sprintf("target=%3.2f", target),
			}

			for i, tc := range tcs {
				temp, err := tc.Temperature()
				if err != nil {
					errCh <- errors.Wrap(err, "read temperature")
					return
				}

				data[i] = append(data[i], Point{
					Time: Duration(t),
					Val:  temp,
				})

				parts = append(parts, fmt.Sprintf("%s=%3.2f", tc.Description, temp))
			}

			// use the value of the first probe for control
			temp := data[0][len(data[0])-1].Val

			if err := gr.CreateValue("reflowoven_temperature", t1, temp); err != nil {
				fmt.Println("error:", errors.Wrap(err, "adding sample to graph"))
			}

			on := target > temp

			parts = append(parts, fmt.Sprintf("on=%v", on))
			fmt.Println(strings.Join(parts, " "))

			err := cook.Out(gpio.Level(on))
			if err != nil {
				errCh <- errors.Wrap(err, "gpio level")
				return
			}

		case <-ctx.Done():
			ticker.Stop()
			_ = cook.Out(gpio.Low)
			fmt.Println("done")
			plot_svg(profile, tcs, data)
			wg.Done()
			return
		}
	}
}
