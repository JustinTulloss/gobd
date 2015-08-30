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
	M1_STATUS
	M1_FREEZE
	M1_FUEL_STATUS
	M1_ENGINE_LOAD
	M1_COOLANT_TEMP
	// Long/short term fuel trim banks 1 and 2
	M1_STFT_1
	M1_LTFT_1
	M1_STFT_2
	M1_LTFT_2
	M1_FUEL_PRESSURE
	M1_MANIFOLD_PRESSURE
	M1_RPM
	M1_SPEED
	M1_TIMING_ADVANCE
	M1_INTAKE_AIR_TEMP
	M1_MAF_AIR_FLOW
	M1_THROTTLE_POSITION
	M1_SECONDARY_AIR_STATUS
	M1_OXYGEN_SENSORS_PRESENT
	// Oxygen sensors, banks 1 and 2, sensors 1-4
	M1_B1_O2_1
	M1_B1_O2_2
	M1_B1_O2_3
	M1_B1_O2_4
	M1_B2_O2_1
	M1_B2_O2_2
	M1_B2_O2_3
	M1_B2_O2_4
	M1_OBD_STANDARDS
	M1_PIDS_2
	M1_DISTANCE_WITH_ENGINE_LIGHT
	M1_FUEL_RAIL_PRESSURE
	M1_FUEL_RAIL_PRESSURE_DIESEL
	// Wide band O2 sensors
	M1_O2_WR_1_VOLTS
	M1_O2_WR_2_VOLTS
	M1_O2_WR_3_VOLTS
	M1_O2_WR_4_VOLTS
	M1_O2_WR_5_VOLTS
	M1_O2_WR_6_VOLTS
	M1_O2_WR_7_VOLTS
	M1_O2_WR_8_VOLTS
	M1_COMMANDED_EGR // https://en.wikipedia.org/wiki/Exhaust_gas_recirculation
	M1_EGR_ERROR
	M1_COMMANDED_EVAPORATIVE_PURGE
	M1_FUEL_LEVEL
	M1_WARM_UPS
	M1_DISTANCE_SINCE_CLEARED
	M1_VAPOR_PRESSURE
	M1_BAROMETRIC_PRESSURE
	M1_O2_WR_1_CURRENT
	M1_O2_WR_2_CURRENT
	M1_O2_WR_3_CURRENT
	M1_O2_WR_4_CURRENT
	M1_O2_WR_5_CURRENT
	M1_O2_WR_6_CURRENT
	M1_O2_WR_7_CURRENT
	M1_O2_WR_8_CURRENT
	M1_CATALYST_TEMP_B1_1
	M1_CATALYST_TEMP_B2_1
	M1_CATALYST_TEMP_B1_2
	M1_CATALYST_TEMP_B2_2
	M1_PIDS_3
	M1_MONITOR_STATUS
	M1_CONTROL_MODULE_VOLTS
	M1_ABSOLUTE_LOAD
	M1_FUEL_AIR_EQUIVALENCE_RATIO
	M1_RELATIVE_THROTTLE
	M1_AMBIENT_AIR_TEMP
	M1_ABSOLUTE_THROTTLE_B
	M1_ABSOLUTE_THROTTLE_C
	M1_ABSOLUTE_THROTTLE_D
	M1_ABSOLUTE_THROTTLE_E
	M1_ABSOLUTE_THROTTLE_F
	M1_COMMANDED_THROTTLE_ACTUATOR
	M1_TIME_WITH_ENGINE_LIGHT
	M1_TIME_SINCE_CLEARED
	M1_MAX_VALUES
	M1_MAX_AIR_FLOW_RATE
	M1_FUEL_TYPE
	M1_ETHANOL_FUEL_PERCENT
	M1_ABSOLUTE_VAPOR_PRESSURE
	M1_ST_SECONDARY_O2_1_3
	M1_LT_SECONDARY_O2_1_3
	M1_ST_SECONDARY_O2_2_4
	M1_LT_SECONDARY_O2_2_4
	M1_FUEL_RAIL_ABSOLUTE_PRESSURE
	M1_RELATIVE_ACCELERATOR
	M1_HYBRID_BATTERY_LIFE
	M1_ENGINE_OIL_TEMP
	M1_EMISSION_REQUIREMENTS
	M1_PIDS_4
	M1_DRIVERS_ENGINE_TORQUE
	M1_ACTUAL_ENGINE_TORQUE
	M1_REFERENCE_ENGINE_TORQUE
	M1_ENGINE_PERCENT_TORQUE
	M1_AUX_INPUT_SUPPORTED
	// There are a lot more, but they're a lot less clear
	// https://en.wikipedia.org/wiki/OBD-II_PIDs
)

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
