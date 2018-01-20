/*
Copyright (c) 2017 Forrest Sibley <My^Name^Without^The^Surname@ieee.org>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package softStepper

import(
	"fmt"
	"time"
	"github.com/the-sibyl/sysfsGPIO"
)

type Stepper struct {
	// Enable output pin to the stepper driver IC
	pinEna *sysfsGPIO.IOPin;
	// Channel A output pin to the stepper driver IC
	pinA *sysfsGPIO.IOPin;
	// Channel B output pin to the stepper driver IC
	pinB *sysfsGPIO.IOPin;
	// Channel C output pin to the stepper driver IC
	pinC *sysfsGPIO.IOPin;
	// Channel D output pin to the stepper driver IC
	pinD *sysfsGPIO.IOPin;
	// Internal flag for the stepper's rotational direction
	stepDirectionForward bool;
	// Stepper state. The valid range is 0 to 3.
	stepState int;
	// Pulse duration in milliseconds
	pulseDuration time.Duration;
	// Internal flag for enabling a hold on the stepper. This flag leaves the driver IC enabled and providing
	// current.
	holdEnable bool;
}

// Create a Stepper struct with enough data to drive a stepper. Less critical values like HoldEnable may be changed by
// the user after initialization.
func InitStepper(enaPin int, pinA int, pinB int, pinC int, pinD int, pulseDuration time.Duration) *Stepper {
	ena, err := sysfsGPIO.InitPin(enaPin, "out")
	initStepperGpioErrorHandler(err)
	a, err := sysfsGPIO.InitPin(pinA, "out")
	initStepperGpioErrorHandler(err)
	b, err := sysfsGPIO.InitPin(pinB, "out")
	initStepperGpioErrorHandler(err)
	c, err := sysfsGPIO.InitPin(pinC, "out")
	initStepperGpioErrorHandler(err)
	d, err := sysfsGPIO.InitPin(pinD, "out")
	initStepperGpioErrorHandler(err)

	stepper := Stepper{
		pinEna: ena,
		pinA: a,
		pinB: b,
		pinC: c,
		pinD: d,
		stepDirectionForward: true,
		stepState: 0,
		pulseDuration: pulseDuration,
		holdEnable: false,
	}

	return &stepper
}

// Helper function for InitStepper to debug
func initStepperGpioErrorHandler(err error) {
	if err != nil {
		fmt.Println("GPIO error while initializing stepper:", err)
	}
}

// Internal generalized step method. The public stepping methods shall call this method.
func (s *Stepper) step(numSteps int) {
	// There are four states total. States will overflow forward or backward, e.g. incrementing state 3 results in
	// state 0, and decrementing state 0 results in state 3.

	if numSteps < 1 {
		fmt.Println("Warning: An invalid number of steps was specified.")
		return
	}

	for k := 0; k < numSteps; k++ {

		// Forward direction case
		if s.stepDirectionForward == true {
			if s.stepState < 3 {
				s.stepState++
			} else {
				s.stepState = 0
			}

		// Reverse direction case
		} else {
			if s.stepState > 0 {
				s.stepState--
			} else {
				s.stepState = 3
			}
		}

		// Set the pin outputs based on the new state
		switch s.stepState {
			case 0:
				s.pinA.SetHigh()
				s.pinB.SetLow()
				s.pinC.SetHigh()
				s.pinD.SetLow()
			case 1:
				s.pinA.SetLow()
				s.pinB.SetHigh()
				s.pinC.SetHigh()
				s.pinD.SetLow()
			case 2:
				s.pinA.SetLow()
				s.pinB.SetHigh()
				s.pinC.SetLow()
				s.pinD.SetHigh()
			case 3:
				s.pinA.SetHigh()
				s.pinB.SetLow()
				s.pinC.SetLow()
				s.pinD.SetHigh()
			default:
				panic("Code error: default stepper state was reached.")
		}

		// Now that the new stepper state is driven on the output pins, assert the enable signal so that the driver IC
		// will provide current to the motor.
		s.pinEna.SetHigh()
		time.Sleep(s.pulseDuration)
	}

	// If the stepper is in holding mode, keep the enable pin asserted so that the coils continue to be driven with
	// the present state.
	if s.holdEnable == false {
		s.pinEna.SetLow()
	}
}

// Run the stepper forward by one step
func (s *Stepper) StepForward() {
	s.stepDirectionForward = true
	s.step(1)
}

// Run the stepper backward by one step
func (s *Stepper) StepBackward() {
	s.stepDirectionForward = false
	s.step(1)
}

// Run the stepper forward by a specified number of steps
func (s *Stepper) StepForwardMulti(numSteps int) {
	s.stepDirectionForward = true
	s.step(numSteps)
}

// Run the stepper backward by a specified number of steps
func (s *Stepper) StepBackwardMulti(numSteps int) {
	s.stepDirectionForward = false
	s.step(numSteps)
}

// Enable stepper holding. Note: this usually consumes a large amount of energy.
func (s *Stepper) EnableHold() {
	s.holdEnable = true
	s.pinEna.SetHigh()
}

// Disable stepper holding
func (s *Stepper) DisableHold() {
	s.holdEnable = false
	s.pinEna.SetLow()
}

