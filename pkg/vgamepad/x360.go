package vgamepad

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"github.com/CB2Moon/vgamepad-go/internal/vigem"
	"github.com/CB2Moon/vgamepad-go/pkg/commons"
)

// VX360Gamepad represents a virtual Xbox 360 gamepad
type VX360Gamepad struct {
	*BaseGamepad
	report commons.XUSBReport
}

// NewVX360Gamepad creates a new virtual Xbox 360 gamepad
func NewVX360Gamepad() (*VX360Gamepad, error) {
	base, err := NewBaseGamepad(func() (uintptr, error) {
		client, err := vigem.NewViGEmClient()
		if err != nil {
			return 0, err
		}
		return client.TargetX360Alloc()
	})
	if err != nil {
		return nil, err
	}

	gamepad := &VX360Gamepad{
		BaseGamepad: base,
		report:      getDefaultX360Report(),
	}

	// Send initial report
	err = gamepad.Update()
	if err != nil {
		gamepad.Close()
		return nil, err
	}

	return gamepad, nil
}

// getDefaultX360Report returns a default Xbox 360 report
func getDefaultX360Report() commons.XUSBReport {
	return commons.XUSBReport{
		WButtons:      0,
		BLeftTrigger:  0,
		BRightTrigger: 0,
		SThumbLX:      0,
		SThumbLY:      0,
		SThumbRX:      0,
		SThumbRY:      0,
	}
}

// Reset resets the gamepad to default state
func (g *VX360Gamepad) Reset() {
	g.report = getDefaultX360Report()
}

// Update sends the current report to the virtual device
func (g *VX360Gamepad) Update() error {
	return g.client.TargetX360Update(g.busp, g.devicep, g.report)
}

// PressButton presses a button (no effect if already pressed)
func (g *VX360Gamepad) PressButton(button commons.XUSBButton) {
	g.report.WButtons = g.report.WButtons | uint16(button)
}

// ReleaseButton releases a button (no effect if already released)
func (g *VX360Gamepad) ReleaseButton(button commons.XUSBButton) {
	g.report.WButtons = g.report.WButtons &^ uint16(button)
}

// LeftTrigger sets the value (0-255, 0 = trigger released) of the left trigger
func (g *VX360Gamepad) LeftTrigger(value uint8) {
	g.report.BLeftTrigger = value
}

// RightTrigger sets the value (0-255, 0 = trigger released) of the right trigger
func (g *VX360Gamepad) RightTrigger(value uint8) {
	g.report.BRightTrigger = value
}

// LeftTriggerFloat sets the value (0.0-1.0, 0.0 = trigger released) of the left trigger using a float
func (g *VX360Gamepad) LeftTriggerFloat(valueFloat float64) {
	g.LeftTrigger(uint8(math.Round(valueFloat * 255)))
}

// RightTriggerFloat sets the value (0.0-1.0, 0.0 = trigger released) of the right trigger using a float
func (g *VX360Gamepad) RightTriggerFloat(valueFloat float64) {
	g.RightTrigger(uint8(math.Round(valueFloat * 255)))
}

// LeftJoystick sets the values (-32768 to 32768, 0 = neutral position) of the X and Y axis for the left joystick
func (g *VX360Gamepad) LeftJoystick(xValue, yValue int16) {
	g.report.SThumbLX = xValue
	g.report.SThumbLY = yValue
}

// RightJoystick sets the values (-32768 to 32768, 0 = neutral position) of the X and Y axis for the right joystick
func (g *VX360Gamepad) RightJoystick(xValue, yValue int16) {
	g.report.SThumbRX = xValue
	g.report.SThumbRY = yValue
}

// LeftJoystickFloat sets the values (-1.0 to 1.0, 0 = neutral position) of the X and Y axis for the left joystick using floats
func (g *VX360Gamepad) LeftJoystickFloat(xValueFloat, yValueFloat float64) {
	g.LeftJoystick(
		int16(math.Round(xValueFloat*32767)),
		int16(math.Round(yValueFloat*32767)),
	)
}

// RightJoystickFloat sets the values (-1.0 to 1.0, 0 = neutral position) of the X and Y axis for the right joystick using floats
func (g *VX360Gamepad) RightJoystickFloat(xValueFloat, yValueFloat float64) {
	g.RightJoystick(
		int16(math.Round(xValueFloat*32767)),
		int16(math.Round(yValueFloat*32767)),
	)
}

// RegisterNotification registers a callback function for notifications
func (g *VX360Gamepad) RegisterNotification(callback NotificationCallback) error {
	// Create a syscall.Callback from the Go function
	callbackPtr := syscall.NewCallback(func(client, target uintptr, largeMotor, smallMotor, ledNumber uint8, userData uintptr) uintptr {
		callback(client, target, largeMotor, smallMotor, ledNumber, userData)
		return 0
	})

	g.cmpFunc = unsafe.Pointer(&callbackPtr)

	err := g.client.TargetX360RegisterNotification(g.busp, g.devicep, uintptr(g.cmpFunc), 0)
	if err != nil {
		return fmt.Errorf("failed to register notification: %w", err)
	}

	return nil
}

// UnregisterNotification unregisters a previously registered callback function
func (g *VX360Gamepad) UnregisterNotification() {
	g.client.TargetX360UnregisterNotification(g.devicep)
	g.cmpFunc = nil
}
