//go:build linux
// +build linux

package main

import (
	"context"
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"strings"
	"sync"
	"time"
)

func monitorTemp(
	ctx context.Context,
	wg *sync.WaitGroup,
	t0 time.Time,
	ch chan error,
) {
	_, err := host.Init()
	noErr(err)

	bus, err := i2creg.Open("1")
	noErr(err)

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
