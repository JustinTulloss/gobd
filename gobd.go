// Package gobd provides an interface for dealing with a vehicle's On board
// diagnostsic system. Supports OBDII.
package gobd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/term"
)

const (
	_ = iota
	MODE_CURRENT
	MODE_FREEZE_FRAME
	MODE_TROUBLE_CODES
	MODE_CLEAR_TROUBLE_CODES
	MODE_TEST_RESULTS_O2
	MODE_TEST_RESULTS_NCT
	MODE_PENDING_TROUBLE_CODES
	MODE_SCM // Special control mode
	MODE_INFO
	MODE_PERMANENT_TROUBLE_CODES
)

const (
	M1_PIDS = iota
)

var RETURN_SIZES = map[byte]byte{
	M1_PIDS: 4,
}

type OBD struct {
	port  *term.Term
	ready bool
}

func (obd *OBD) waitTilReady() {
	if obd.ready {
		return
	}
	respReader := bufio.NewReader(obd.port)
	// Wait for the obd to signal that it's ready
	_, err := respReader.ReadBytes('>')
	if err != nil && err != io.EOF {
		panic("Error occurred waiting for ready: " + err.Error())
	}
	obd.ready = true
}

func (obd *OBD) SendCommand(command []byte) error {
	fmt.Println("Sending command", string(command))
	obd.waitTilReady()
	err := obd.port.Flush()
	if err != nil {
		return err
	}
	message := []byte{}
	message = append(message, command...)
	message = append(message, '\r')
	n, err := obd.port.Write(message)
	if err != nil {
		return err
	}
	obd.ready = false
	fmt.Printf("Wrote %d bytes\n", n)
	return nil
}

func (obd *OBD) ReadResult() ([]byte, error) {
	respReader := bufio.NewReader(obd.port)
	response, err := respReader.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	response = response[:len(response)-1] // Drop the carriage return
	fmt.Printf("read: '%s'\n", string(response))
	return response, err
}

func (obd *OBD) Reset() error {
	err := obd.SendCommand([]byte{'a', 't', 'z'})
	obd.ReadResult()
	return err
}

func (obd *OBD) Close() error {
	obd.Reset()
	return obd.port.Close()
}

func NewOBD(serialPortName string) (*OBD, error) {
	port, err := term.Open(serialPortName, term.Speed(9600), term.RawMode)
	if err != nil {
		return nil, err
	}
	obd := &OBD{port, true}
	err = obd.Reset()
	if err != nil {
		return nil, err
	}
	err = obd.SendCommand([]byte("ate0")) // Turns off echo
	obd.ReadResult()
	return obd, err
}
