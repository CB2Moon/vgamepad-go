package commons

import (
	"fmt"
)

// ViGEmTargetType represents the desired target type for the emulated device
type ViGEmTargetType int

const (
	Xbox360Wired    ViGEmTargetType = 0 // Microsoft Xbox 360 Controller (wired)
	DualShock4Wired ViGEmTargetType = 2 // Sony DualShock 4 (wired)
)

// XUSBButton represents possible XUSB report buttons
type XUSBButton uint16

const (
	XUSB_GAMEPAD_DPAD_UP        XUSBButton = 0x0001
	XUSB_GAMEPAD_DPAD_DOWN      XUSBButton = 0x0002
	XUSB_GAMEPAD_DPAD_LEFT      XUSBButton = 0x0004
	XUSB_GAMEPAD_DPAD_RIGHT     XUSBButton = 0x0008
	XUSB_GAMEPAD_START          XUSBButton = 0x0010
	XUSB_GAMEPAD_BACK           XUSBButton = 0x0020
	XUSB_GAMEPAD_LEFT_THUMB     XUSBButton = 0x0040
	XUSB_GAMEPAD_RIGHT_THUMB    XUSBButton = 0x0080
	XUSB_GAMEPAD_LEFT_SHOULDER  XUSBButton = 0x0100
	XUSB_GAMEPAD_RIGHT_SHOULDER XUSBButton = 0x0200
	XUSB_GAMEPAD_GUIDE          XUSBButton = 0x0400
	XUSB_GAMEPAD_A              XUSBButton = 0x1000
	XUSB_GAMEPAD_B              XUSBButton = 0x2000
	XUSB_GAMEPAD_X              XUSBButton = 0x4000
	XUSB_GAMEPAD_Y              XUSBButton = 0x8000
)

// XUSBReport represents an XINPUT_GAMEPAD-compatible report structure
type XUSBReport struct {
	WButtons      uint16
	BLeftTrigger  uint8
	BRightTrigger uint8
	SThumbLX      int16
	SThumbLY      int16
	SThumbRX      int16
	SThumbRY      int16
}

// DS4Button represents DualShock 4 digital buttons
type DS4Button uint16

const (
	DS4_BUTTON_THUMB_RIGHT    DS4Button = 1 << 15
	DS4_BUTTON_THUMB_LEFT     DS4Button = 1 << 14
	DS4_BUTTON_OPTIONS        DS4Button = 1 << 13
	DS4_BUTTON_SHARE          DS4Button = 1 << 12
	DS4_BUTTON_TRIGGER_RIGHT  DS4Button = 1 << 11
	DS4_BUTTON_TRIGGER_LEFT   DS4Button = 1 << 10
	DS4_BUTTON_SHOULDER_RIGHT DS4Button = 1 << 9
	DS4_BUTTON_SHOULDER_LEFT  DS4Button = 1 << 8
	DS4_BUTTON_TRIANGLE       DS4Button = 1 << 7
	DS4_BUTTON_CIRCLE         DS4Button = 1 << 6
	DS4_BUTTON_CROSS          DS4Button = 1 << 5
	DS4_BUTTON_SQUARE         DS4Button = 1 << 4
)

// DS4SpecialButton represents DualShock 4 special buttons
type DS4SpecialButton uint8

const (
	DS4_SPECIAL_BUTTON_PS       DS4SpecialButton = 1 << 0
	DS4_SPECIAL_BUTTON_TOUCHPAD DS4SpecialButton = 1 << 1 // Windows only, no effect on Linux
)

// DS4DPadDirection represents DualShock 4 directional pad (HAT) values
type DS4DPadDirection uint8

const (
	DS4_BUTTON_DPAD_NONE      DS4DPadDirection = 0x8
	DS4_BUTTON_DPAD_NORTHWEST DS4DPadDirection = 0x7
	DS4_BUTTON_DPAD_WEST      DS4DPadDirection = 0x6
	DS4_BUTTON_DPAD_SOUTHWEST DS4DPadDirection = 0x5
	DS4_BUTTON_DPAD_SOUTH     DS4DPadDirection = 0x4
	DS4_BUTTON_DPAD_SOUTHEAST DS4DPadDirection = 0x3
	DS4_BUTTON_DPAD_EAST      DS4DPadDirection = 0x2
	DS4_BUTTON_DPAD_NORTHEAST DS4DPadDirection = 0x1
	DS4_BUTTON_DPAD_NORTH     DS4DPadDirection = 0x0
)

// DS4Report represents DualShock 4 HID Input report
type DS4Report struct {
	BThumbLX  uint8
	BThumbLY  uint8
	BThumbRX  uint8
	BThumbRY  uint8
	WButtons  uint16
	BSpecial  uint8
	BTriggerL uint8
	BTriggerR uint8
}

// DS4LightbarColor represents the color value (RGB) of a DualShock 4 Lightbar
type DS4LightbarColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// DS4Touch represents DualShock 4 HID Touchpad structure
type DS4Touch struct {
	BPacketCounter    uint8
	BIsUpTrackingNum1 uint8
	BTouchData1       [3]uint8
	BIsUpTrackingNum2 uint8
	BTouchData2       [3]uint16
}

// DS4ReportEx represents DualShock 4 v1 complete HID Input report
type DS4ReportEx struct {
	Report       DS4SubReportEx
	ReportBuffer [63]uint8
}

// DS4SubReportEx represents the detailed structure of a DS4 extended report
type DS4SubReportEx struct {
	BThumbLX           uint8
	BThumbLY           uint8
	BThumbRX           uint8
	BThumbRY           uint8
	WButtons           uint16
	BSpecial           uint8
	BTriggerL          uint8
	BTriggerR          uint8
	WTimestamp         uint16
	BBatteryLvl        uint8
	WGyroX             int16
	WGyroY             int16
	WGyroZ             int16
	WAccelX            int16
	WAccelY            int16
	WAccelZ            int16
	BUnknown1          [5]uint8
	BBatteryLvlSpecial uint8
	BUnknown2          [2]uint8
	BTouchPacketsN     uint8
	SCurrentTouch      DS4Touch
	SPreviousTouch     [2]DS4Touch
}

// ViGEmError represents values that represent ViGEm errors
type ViGEmError uint32

const (
	VIGEM_ERROR_NONE                        ViGEmError = 0x20000000
	VIGEM_ERROR_BUS_NOT_FOUND               ViGEmError = 0xE0000001
	VIGEM_ERROR_NO_FREE_SLOT                ViGEmError = 0xE0000002
	VIGEM_ERROR_INVALID_TARGET              ViGEmError = 0xE0000003
	VIGEM_ERROR_REMOVAL_FAILED              ViGEmError = 0xE0000004
	VIGEM_ERROR_ALREADY_CONNECTED           ViGEmError = 0xE0000005
	VIGEM_ERROR_TARGET_UNINITIALIZED        ViGEmError = 0xE0000006
	VIGEM_ERROR_TARGET_NOT_PLUGGED_IN       ViGEmError = 0xE0000007
	VIGEM_ERROR_BUS_VERSION_MISMATCH        ViGEmError = 0xE0000008
	VIGEM_ERROR_BUS_ACCESS_FAILED           ViGEmError = 0xE0000009
	VIGEM_ERROR_CALLBACK_ALREADY_REGISTERED ViGEmError = 0xE0000010
	VIGEM_ERROR_CALLBACK_NOT_FOUND          ViGEmError = 0xE0000011
	VIGEM_ERROR_BUS_ALREADY_CONNECTED       ViGEmError = 0xE0000012
	VIGEM_ERROR_BUS_INVALID_HANDLE          ViGEmError = 0xE0000013
	VIGEM_ERROR_XUSB_USERINDEX_OUT_OF_RANGE ViGEmError = 0xE0000014
	VIGEM_ERROR_INVALID_PARAMETER           ViGEmError = 0xE0000015
	VIGEM_ERROR_NOT_SUPPORTED               ViGEmError = 0xE0000016
)

// Error returns a string representation of the ViGEmError
func (e ViGEmError) Error() string {
	switch e {
	case VIGEM_ERROR_NONE:
		return "No error"
	case VIGEM_ERROR_BUS_NOT_FOUND:
		return "Bus not found"
	case VIGEM_ERROR_NO_FREE_SLOT:
		return "No free slot"
	case VIGEM_ERROR_INVALID_TARGET:
		return "Invalid target"
	case VIGEM_ERROR_REMOVAL_FAILED:
		return "Removal failed"
	case VIGEM_ERROR_ALREADY_CONNECTED:
		return "Already connected"
	case VIGEM_ERROR_TARGET_UNINITIALIZED:
		return "Target uninitialized"
	case VIGEM_ERROR_TARGET_NOT_PLUGGED_IN:
		return "Target not plugged in"
	case VIGEM_ERROR_BUS_VERSION_MISMATCH:
		return "Bus version mismatch"
	case VIGEM_ERROR_BUS_ACCESS_FAILED:
		return "Bus access failed"
	case VIGEM_ERROR_CALLBACK_ALREADY_REGISTERED:
		return "Callback already registered"
	case VIGEM_ERROR_CALLBACK_NOT_FOUND:
		return "Callback not found"
	case VIGEM_ERROR_BUS_ALREADY_CONNECTED:
		return "Bus already connected"
	case VIGEM_ERROR_BUS_INVALID_HANDLE:
		return "Bus invalid handle"
	case VIGEM_ERROR_XUSB_USERINDEX_OUT_OF_RANGE:
		return "XUSB user index out of range"
	case VIGEM_ERROR_INVALID_PARAMETER:
		return "Invalid parameter"
	case VIGEM_ERROR_NOT_SUPPORTED:
		return "Not supported"
	default:
		return fmt.Sprintf("Unknown error: %d", e)
	}
}

// DS4SetDPad sets the directional pad value in a DS4Report
func DS4SetDPad(report *DS4Report, dpad DS4DPadDirection) {
	report.WButtons &= ^uint16(0xF)
	report.WButtons |= uint16(dpad)
}

// DS4ReportInit initializes a DS4Report with default values
func DS4ReportInit(report *DS4Report) {
	report.BThumbLX = 0x80
	report.BThumbLY = 0x80
	report.BThumbRX = 0x80
	report.BThumbRY = 0x80
	DS4SetDPad(report, DS4_BUTTON_DPAD_NONE)
}
