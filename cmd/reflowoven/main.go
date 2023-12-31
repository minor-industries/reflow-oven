package main

import (
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

func main() {
	_, err := host.Init()
	noErr(err)

	all := gpioreg.All()
	fmt.Println(all)

	bus, err := i2creg.Open("1")
	noErr(err)

	tc := NewThermocouple(func(s string) {
		fmt.Println(s)
	},
		bus,
		0x60,
	)

	go monitorTemp(tc)

	for _, io := range all {
		fmt.Println(io.Name())
	}

	cook := gpioreg.ByName("GPIO16")
	if cook == nil {
		panic("no gpio")
	}

	for {
		err := cook.Out(gpio.High)
		noErr(err)

		time.Sleep(500 * time.Millisecond)

		err = cook.Out(gpio.Low)
		noErr(err)

		time.Sleep(4500 * time.Millisecond)
	}
}

func monitorTemp(tc *Thermocouple) {
	ticker := time.NewTicker(250 * time.Millisecond)
	for range ticker.C {
		t, err := tc.Temperature()
		noErr(err)

		fmt.Println(t)
	}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
