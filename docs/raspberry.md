Connecting OBDII to a Raspberry PI
----------------------------------

### Equipment
- [BAFX OBDII bluetooth adapter for Android][1]
- [Plugable USB Bluetooth 4.0 LE adapter][2]

### Set up
1. Put the Raspberry PI in the car. The bluetooth range is not great.
1. Install the necessary bluetooth utilities for raspbian `sudo apt-get install bluetooth bluez-utils bluez python-gobject`
1. `hcitool scan`
  - This should show a device called "OBDII". Copy the address.
1. Pair with bluez. `sudo bluez-simple-agent hci0 <address from previous step>`
  - You'll be asked to enter your pairing code. If you're using the BAFX bt adapter, it's 1234
1. Create a serial port. `sudo rfcomm bind hci0 <bt address>`
1. Verify that it works.
    - `screen /dev/rfcomm0`. You're now sending commands straight to the OBDII board
    - `at@1` - This should return "OBDII to RS232 Interpreter"
    - `0100` - This should return the list of PIDs supported by the OBD. It only works if the ECU is on, so you'll need to turn the car on (or at least in acc mode).


[1]: http://www.amazon.com/BAFX-Products-Bluetooth-diagnostics-Android/dp/B005NLQAHS
[2]: http://www.amazon.com/Plugable-Bluetooth-Adapter-Raspberry-Compatible/dp/B009ZIILLI
