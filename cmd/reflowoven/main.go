package main

import (
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

var profile Schedule = Schedule{
	{Duration(0 * time.Second), 25},
	{Duration(30 * time.Second), 100},
	{Duration(120 * time.Second), 150},
	{Duration(150 * time.Second), 183},
	{Duration(210 * time.Second), 235},
	{Duration(240 * time.Second), 183},
	{Duration(300 * time.Second), 25},
	{Duration(330 * time.Second), 25},
}

func main() {
	t0 := time.Now()
	fmt.Println(t0)

	_, err := host.Init()
	noErr(err)

	bus, err := i2creg.Open("1")
	noErr(err)

	log := func(s string) {
		fmt.Println(s)
	}

	tc := NewThermocouple(log, bus, 0x60)

	cook := gpioreg.ByName("GPIO16")
	if cook == nil {
		panic("no gpio")
	}

	go monitorTemp(t0, tc, cook)

	select {}
}

func monitorTemp(t0 time.Time, tc *Thermocouple, cook gpio.PinIO) {
	ticker := time.NewTicker(250 * time.Millisecond)
	for range ticker.C {
		t1 := time.Now()
		t := t1.Sub(t0)

		target := profile.Val(t)

		temp, err := tc.Temperature()
		noErr(err)

		on := target > temp
		fmt.Println(t, target, temp, on)

		err = cook.Out(gpio.Level(on))
		noErr(err)
	}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
