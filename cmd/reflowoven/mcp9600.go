package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"periph.io/x/conn/v3/i2c"
)

type Thermocouple struct {
	log  func(string)
	i2c  i2c.Bus
	addr byte
}

func NewThermocouple(
	log func(string),
	i2c i2c.Bus,
	addr byte,
) *Thermocouple {
	return &Thermocouple{
		log:  log,
		i2c:  i2c,
		addr: addr,
	}
}

func (t *Thermocouple) Temperature() (float64, error) {
	var data [2]byte

	const register = 0x00
	err := t.i2c.Tx(uint16(t.addr), []byte{register}, data[:])

	if err != nil {
		return 0, errors.Wrap(err, "read register")
	} else {
		t.log("read: " + hex.Dump(data[:]))
		t := float64(binary.BigEndian.Uint16(data[:])) * 0.0625
		if data[0]&0x80 != 0 {
			t -= 4096
		}
		fmt.Println("mpc t", t)
		return t, nil
	}
}
