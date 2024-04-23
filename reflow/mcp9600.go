package reflow

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"periph.io/x/conn/v3/i2c"
)

type Thermocouple struct {
	log         func(string)
	i2c         i2c.Bus
	addr        byte
	Description string
}

func NewThermocouple(
	log func(string),
	i2c i2c.Bus,
	addr byte,
	description string,
) *Thermocouple {
	return &Thermocouple{
		log:         log,
		i2c:         i2c,
		addr:        addr,
		Description: description,
	}
}

func (t *Thermocouple) Temperature() (float64, error) {
	var data [2]byte

	const register = 0x00
	err := t.i2c.Tx(uint16(t.addr), []byte{register}, data[:])

	if err != nil {
		return 0, errors.Wrap(err, "read register")
	} else {
		t := float64(binary.BigEndian.Uint16(data[:])) * 0.0625
		if data[0]&0x80 != 0 {
			t -= 4096
		}
		return t, nil
	}
}
