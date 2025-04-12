package vgamepad

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"github.com/CB2Moon/vgamepad-go/internal/vigem"
	"github.com/CB2Moon/vgamepad-go/pkg/commons"
)

// VDS4Gamepad represents a virtual DualShock 4 gamepad
type VDS4Gamepad struct {
	*BaseGamepad
	report commons.DS4Report
}

// NewVDS4Gamepad creates a new virtual DualShock 4 gamepad
func NewVDS4Gamepad() (*VDS4Gamepad, error) {
	base, err := NewBaseGamepad(func() (uintptr, error) {
		client, err := vigem.NewViGEmClient()
		if err != nil {
			return 0, err
		}
		return client.TargetDS4Alloc()
	})
	if err != nil {
		return nil, err
	}

	gamepad := &VDS4Gamepad{
		BaseGamepad: base,
		report:      getDefaultDS4Report(),
	}

	// Send initial report
	err = gamepad.Update()
	if err != nil {
		gamepad.Close()
		return nil, err
	}

	return gamepad, nil
}

// getDefaultDS4Report returns a default DualShock 4 report
func getDefaultDS4Report() commons.DS4Report {
	report := commons.DS4Report{
		BThumbLX:  0,
		BThumbLY:  0,
		BThumbRX:  0,
		BThumbRY:  0,
		WButtons:  0,
		BSpecial:  0,
		BTriggerL: 0,
		BTriggerR: 0,
	}
	commons.DS4ReportInit(&report)
	return report
}

// Reset resets the gamepad to default state
func (g *VDS4Gamepad) Reset() {
	g.report = getDefaultDS4Report()
}

// Update sends the current report to the virtual device
func (g *VDS4Gamepad) Update() error {
	return g.client.TargetDS4Update(g.busp, g.devicep, g.report)
}

// PressButton presses a button (no effect if already pressed)
func (g *VDS4Gamepad) PressButton(button commons.DS4Button) {
	g.report.WButtons = g.report.WButtons | uint16(button)
}

// ReleaseButton releases a button (no effect if already released)
func (g *VDS4Gamepad) ReleaseButton(button commons.DS4Button) {
	g.report.WButtons = g.report.WButtons &^ uint16(button)
}

// PressSpecialButton presses a special button (no effect if already pressed)
func (g *VDS4Gamepad) PressSpecialButton(specialButton commons.DS4SpecialButton) {
	g.report.BSpecial = g.report.BSpecial | uint8(specialButton)
}

// ReleaseSpecialButton releases a special button (no effect if already released)
func (g *VDS4Gamepad) ReleaseSpecialButton(specialButton commons.DS4SpecialButton) {
	g.report.BSpecial = g.report.BSpecial &^ uint8(specialButton)
}

// LeftTrigger sets the value (0-255, 0 = trigger released) of the left trigger
func (g *VDS4Gamepad) LeftTrigger(value uint8) {
	g.report.BTriggerL = value
}

// RightTrigger sets the value (0-255, 0 = trigger released) of the right trigger
func (g *VDS4Gamepad) RightTrigger(value uint8) {
	g.report.BTriggerR = value
}

// LeftTriggerFloat sets the value (0.0-1.0, 0.0 = trigger released) of the left trigger using a float
func (g *VDS4Gamepad) LeftTriggerFloat(valueFloat float64) {
	g.LeftTrigger(uint8(math.Round(valueFloat * 255)))
}

// RightTriggerFloat sets the value (0.0-1.0, 0.0 = trigger released) of the right trigger using a float
func (g *VDS4Gamepad) RightTriggerFloat(valueFloat float64) {
	g.RightTrigger(uint8(math.Round(valueFloat * 255)))
}

// LeftJoystick sets the values (0-255, 128 = neutral position) of the X and Y axis for the left joystick
func (g *VDS4Gamepad) LeftJoystick(xValue, yValue uint8) {
	g.report.BThumbLX = xValue
	g.report.BThumbLY = yValue
}

// RightJoystick sets the values (0-255, 128 = neutral position) of the X and Y axis for the right joystick
func (g *VDS4Gamepad) RightJoystick(xValue, yValue uint8) {
	g.report.BThumbRX = xValue
	g.report.BThumbRY = yValue
}

// LeftJoystickFloat sets the values (-1.0 to 1.0, 0 = neutral position) of the X and Y axis for the left joystick using floats
func (g *VDS4Gamepad) LeftJoystickFloat(xValueFloat, yValueFloat float64) {
	g.LeftJoystick(
		uint8(128+math.Round(xValueFloat*127)),
		uint8(128+math.Round(yValueFloat*127)),
	)
}

// RightJoystickFloat sets the values (-1.0 to 1.0, 0 = neutral position) of the X and Y axis for the right joystick using floats
func (g *VDS4Gamepad) RightJoystickFloat(xValueFloat, yValueFloat float64) {
	g.RightJoystick(
		uint8(128+math.Round(xValueFloat*127)),
		uint8(128+math.Round(yValueFloat*127)),
	)
}

// DirectionalPad sets the direction of the directional pad (hat)
func (g *VDS4Gamepad) DirectionalPad(direction commons.DS4DPadDirection) {
	commons.DS4SetDPad(&g.report, direction)
}

// UpdateExtendedReport enables using DS4_REPORT_EX instead of DS4_REPORT (advanced users only)
func (g *VDS4Gamepad) UpdateExtendedReport(extendedReport *commons.DS4ReportEx) error {
	return g.client.TargetDS4UpdateExPtr(g.busp, g.devicep, extendedReport)
}

// RegisterNotification registers a callback function for notifications
func (g *VDS4Gamepad) RegisterNotification(callback NotificationCallback) error {
	// Create a syscall.Callback from the Go function
	callbackPtr := syscall.NewCallback(func(client, target uintptr, largeMotor, smallMotor, ledNumber uint8, userData uintptr) uintptr {
		callback(client, target, largeMotor, smallMotor, ledNumber, userData)
		return 0
	})

	g.cmpFunc = unsafe.Pointer(&callbackPtr)

	err := g.client.TargetDS4RegisterNotification(g.busp, g.devicep, uintptr(g.cmpFunc), 0)
	if err != nil {
		return fmt.Errorf("failed to register notification: %w", err)
	}

	return nil
}

// UnregisterNotification unregisters a previously registered callback function
func (g *VDS4Gamepad) UnregisterNotification() {
	g.client.TargetDS4UnregisterNotification(g.devicep)
	g.cmpFunc = nil
}
