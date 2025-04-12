package main

import (
	"fmt"
	"time"

	"github.com/CB2Moon/vgamepad-go/pkg/commons"
	"github.com/CB2Moon/vgamepad-go/pkg/vgamepad"
)

func main() {
	// Example with Xbox 360 gamepad
	fmt.Println("Creating Xbox 360 gamepad...")
	x360Gamepad, err := vgamepad.NewVX360Gamepad()
	if err != nil {
		fmt.Printf("Failed to create Xbox 360 gamepad: %v\n", err)
		return
	}
	defer x360Gamepad.Close()

	// Press a button to wake the device up
	fmt.Println("Pressing A button...")
	x360Gamepad.PressButton(commons.XUSB_GAMEPAD_A)
	x360Gamepad.Update()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Releasing A button...")
	x360Gamepad.ReleaseButton(commons.XUSB_GAMEPAD_A)
	x360Gamepad.Update()
	time.Sleep(500 * time.Millisecond)

	// Press buttons and things
	fmt.Println("Pressing multiple buttons and setting joysticks...")
	x360Gamepad.PressButton(commons.XUSB_GAMEPAD_A)
	x360Gamepad.PressButton(commons.XUSB_GAMEPAD_LEFT_SHOULDER)
	x360Gamepad.PressButton(commons.XUSB_GAMEPAD_DPAD_DOWN)
	x360Gamepad.PressButton(commons.XUSB_GAMEPAD_DPAD_LEFT)
	x360Gamepad.LeftTriggerFloat(0.5)
	x360Gamepad.RightTriggerFloat(0.5)
	x360Gamepad.LeftJoystickFloat(0.0, 0.2)
	x360Gamepad.RightJoystickFloat(-1.0, 1.0)
	x360Gamepad.Update()
	time.Sleep(1 * time.Second)

	// Release buttons and things
	fmt.Println("Releasing some buttons...")
	x360Gamepad.ReleaseButton(commons.XUSB_GAMEPAD_A)
	x360Gamepad.ReleaseButton(commons.XUSB_GAMEPAD_DPAD_LEFT)
	x360Gamepad.RightTriggerFloat(0.0)
	x360Gamepad.RightJoystickFloat(0.0, 0.0)
	x360Gamepad.Update()
	time.Sleep(1 * time.Second)

	// Reset gamepad to default state
	fmt.Println("Resetting gamepad...")
	x360Gamepad.Reset()
	x360Gamepad.Update()
	time.Sleep(1 * time.Second)

	// Example with DualShock 4 gamepad
	fmt.Println("\nCreating DualShock 4 gamepad...")
	ds4Gamepad, err := vgamepad.NewVDS4Gamepad()
	if err != nil {
		fmt.Printf("Failed to create DualShock 4 gamepad: %v\n", err)
		return
	}
	defer ds4Gamepad.Close()

	// Press a button to wake the device up
	fmt.Println("Pressing Triangle button...")
	ds4Gamepad.PressButton(commons.DS4_BUTTON_TRIANGLE)
	ds4Gamepad.Update()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Releasing Triangle button...")
	ds4Gamepad.ReleaseButton(commons.DS4_BUTTON_TRIANGLE)
	ds4Gamepad.Update()
	time.Sleep(500 * time.Millisecond)

	// Press buttons and things
	fmt.Println("Pressing multiple buttons and setting joysticks...")
	ds4Gamepad.PressButton(commons.DS4_BUTTON_TRIANGLE)
	ds4Gamepad.PressButton(commons.DS4_BUTTON_CIRCLE)
	ds4Gamepad.PressButton(commons.DS4_BUTTON_THUMB_RIGHT)
	ds4Gamepad.PressButton(commons.DS4_BUTTON_TRIGGER_LEFT)
	ds4Gamepad.PressSpecialButton(commons.DS4_SPECIAL_BUTTON_TOUCHPAD)
	ds4Gamepad.LeftTriggerFloat(0.5)
	ds4Gamepad.RightTriggerFloat(0.5)
	ds4Gamepad.LeftJoystickFloat(0.0, 0.2)
	ds4Gamepad.RightJoystickFloat(-1.0, 1.0)
	ds4Gamepad.Update()
	time.Sleep(1 * time.Second)

	// Release buttons and things
	fmt.Println("Releasing some buttons...")
	ds4Gamepad.ReleaseButton(commons.DS4_BUTTON_TRIANGLE)
	ds4Gamepad.RightTriggerFloat(0.0)
	ds4Gamepad.RightJoystickFloat(0.0, 0.0)
	ds4Gamepad.Update()
	time.Sleep(1 * time.Second)

	// Reset gamepad to default state
	fmt.Println("Resetting gamepad...")
	ds4Gamepad.Reset()
	ds4Gamepad.Update()
	time.Sleep(1 * time.Second)

	fmt.Println("Example completed successfully!")
}
