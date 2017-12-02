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

package main

import(
	"fmt"
	"time"
	"github.com/the-sibyl/sysfsGPIO"
)

type Stepper struct {
	PinEna *sysfsGPIO.IOPin;
	PinA *sysfsGPIO.IOPin;
	PinB *sysfsGPIO.IOPin;
	PinC *sysfsGPIO.IOPin;
	PinD *sysfsGPIO.IOPin;
	StepState int;
	PulseDuration int;
}

func InitStepper(enaPin int, pinA int, pinB int, pinC int, pinD int) *Stepper {
	ena, _ := sysfsGPIO.InitPin(enaPin, "out")
	a, _ := sysfsGPIO.InitPin(pinA, "out")
	b, _ := sysfsGPIO.InitPin(pinB, "out")
	c, _ := sysfsGPIO.InitPin(pinC, "out")
	d, _ := sysfsGPIO.InitPin(pinD, "out")

	stepper := Stepper{
		PinEna: ena,
		PinA: a,
		PinB: b,
		PinC: c,
		PinD: d,
		StepState: 0,
		PulseDuration: 15,
	}

	return &stepper
}

func (s *Stepper) Step() {
	switch s.StepState {
		case 0:
			fmt.Println(s.StepState)
			s.PinA.SetHigh()
			s.PinB.SetLow()
			s.PinC.SetLow()
			s.PinD.SetHigh()
			s.StepState = 1
		case 1:
			fmt.Println(s.StepState)
			s.PinA.SetLow()
			s.PinB.SetHigh()
			s.PinC.SetLow()
			s.PinD.SetHigh()
			s.StepState = 2
		case 2:
			fmt.Println(s.StepState)
			s.PinA.SetLow()
			s.PinB.SetHigh()
			s.PinC.SetHigh()
			s.PinD.SetLow()
			s.StepState = 3
		case 3:
			fmt.Println(s.StepState)
			s.PinA.SetHigh()
			s.PinB.SetLow()
			s.PinC.SetHigh()
			s.PinD.SetLow()
			s.StepState = 0
		default:
			s.StepState = 0
	}

	fmt.Println("Enabling")
	s.PinEna.SetHigh()
	time.Sleep(time.Millisecond * time.Duration(s.PulseDuration))
	s.PinEna.SetLow()
}

func main() {
//	stepper1 := InitStepper(8, 9, 7, 0, 2)
//	stepper2 := InitStepper(30, 21, 22, 23, 24)
//	stepper3 := InitStepper(31, 26, 27, 28, 29)

	stepper1 := InitStepper(2, 3, 4, 17, 27)

//	stepper2 := InitStepper(, 5, 6, 13, 19)
//	stepper3 := InitStepper(, 12, 16, 20, 21)

	fmt.Println(stepper1)
//	fmt.Println(stepper2)

	for k := 0; k < 100; k++ {
		stepper1.Step()
		time.Sleep(time.Millisecond * 200)
//		stepper2.Step()
//		stepper3.Step()
	}
}
