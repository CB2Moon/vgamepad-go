# Virtual Gamepad for Go

The Go implementation of [Virtual XBox360 and DualShock4 gamepads in python](https://github.com/yannbouteiller/vgamepad).

---

Virtual Gamepad for Go (`vgamepad-go`) is a Go library that emulates XBox360 and DualShock4 gamepads on your system.
It enables controlling applications that require gamepad input (such as video games) directly from your Go code.

On Windows, `vgamepad-go` uses the [Virtual Gamepad Emulation](https://github.com/nefarius/ViGEmBus) C++ framework, providing Go bindings and a user-friendly interface.

---

__Development status:__

|  Windows  |  Linux  |
|:---------:|:-------:|
| *Stable.* | *Not yet supported* |

## Quick links
- [Installation](#installation)
- [Getting started](#getting-started)
  - [XBox360 gamepad](#xbox360-gamepad)
  - [DualShock4 gamepad](#dualshock4-gamepad)
  - [Rumble and LEDs](#rumble-and-leds)
- [Local Development](#local-development)
- [Publishing](#publishing)
- [Contribute](#contribute)

---

## Installation

### Prerequisites:

1. Install the ViGEmBus driver from [here](https://github.com/nefarius/ViGEmBus/releases).
   - Download the latest release
   - Run the installer
   - Accept the license agreement
   - Allow the installer to modify your PC
   - Wait for completion and click "Finish"

2. Install Go (version 1.16 or later) from [golang.org](https://golang.org/dl/).

### Installing the library:

```bash
go get github.com/CB2Moon/vgamepad-go
```

---

## Getting started

**You need to run as admin**

`vgamepad-go` provides two main Go types: `VX360Gamepad`, which emulates an XBox360 gamepad, and `VDS4Gamepad`, which emulates a DualShock4 gamepad.

The state of a virtual gamepad (e.g., pressed buttons, joystick values) is called a report.
To modify the report, a number of user-friendly API functions are provided by `vgamepad-go`.
When the report is modified as desired, it must be sent to the computer using the `Update()` method.

### XBox360 gamepad

The following Go code creates a virtual XBox360 gamepad:

```go
import "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"

gamepad, err := vgamepad.NewVX360Gamepad()
if err != nil {
    // Handle error
}
defer gamepad.Close()
```

As soon as the `VX360Gamepad` object is created, the virtual gamepad is connected to your system via the ViGEmBus driver, and will remain connected until the object is destroyed.

Buttons can be pressed and released through `PressButton` and `ReleaseButton`:

```go
import (
    "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"
    "github.com/CB2Moon/vgamepad-go/internal/vigem"
)

// Press the A button
gamepad.PressButton(vigem.XUSB_GAMEPAD_A)
// Press the left hat button
gamepad.PressButton(vigem.XUSB_GAMEPAD_DPAD_LEFT)

// Send the updated state to the computer
gamepad.Update()

// (...) A and left hat are pressed...

// Release the A button
gamepad.ReleaseButton(vigem.XUSB_GAMEPAD_A)

// Send the updated state to the computer
gamepad.Update()

// (...) left hat is still pressed...
```

All available buttons are defined in the `vigem` package as `XUSBButton` constants.

To control the triggers (1 axis each) and the joysticks (2 axis each), two options are provided by the API.

It is possible to input raw integer values directly:

```go
// Left trigger: value between 0 and 255
gamepad.LeftTrigger(100)
// Right trigger: value between 0 and 255
gamepad.RightTrigger(255)
// Left joystick: values between -32768 and 32767
gamepad.LeftJoystick(-10000, 0)
// Right joystick: values between -32768 and 32767
gamepad.RightJoystick(-32768, 15000)

gamepad.Update()
```

Or to input float values:

```go
// Left trigger: value between 0.0 and 1.0
gamepad.LeftTriggerFloat(0.5)
// Right trigger: value between 0.0 and 1.0
gamepad.RightTriggerFloat(1.0)
// Left joystick: values between -1.0 and 1.0
gamepad.LeftJoystickFloat(-0.5, 0.0)
// Right joystick: values between -1.0 and 1.0
gamepad.RightJoystickFloat(-1.0, 0.8)

gamepad.Update()
```

Reset to default state:

```go
gamepad.Reset()
gamepad.Update()
```

Full example:

```go
package main

import (
    "time"

    "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"
    "github.com/CB2Moon/vgamepad-go/internal/vigem"
)

func main() {
    gamepad, err := vgamepad.NewVX360Gamepad()
    if err != nil {
        panic(err)
    }
    defer gamepad.Close()

    // Press a button to wake the device up
    gamepad.PressButton(vigem.XUSB_GAMEPAD_A)
    gamepad.Update()
    time.Sleep(500 * time.Millisecond)
    gamepad.ReleaseButton(vigem.XUSB_GAMEPAD_A)
    gamepad.Update()
    time.Sleep(500 * time.Millisecond)

    // Press buttons and things
    gamepad.PressButton(vigem.XUSB_GAMEPAD_A)
    gamepad.PressButton(vigem.XUSB_GAMEPAD_LEFT_SHOULDER)
    gamepad.PressButton(vigem.XUSB_GAMEPAD_DPAD_DOWN)
    gamepad.PressButton(vigem.XUSB_GAMEPAD_DPAD_LEFT)
    gamepad.LeftTriggerFloat(0.5)
    gamepad.RightTriggerFloat(0.5)
    gamepad.LeftJoystickFloat(0.0, 0.2)
    gamepad.RightJoystickFloat(-1.0, 1.0)

    gamepad.Update()

    time.Sleep(1 * time.Second)

    // Release buttons and things
    gamepad.ReleaseButton(vigem.XUSB_GAMEPAD_A)
    gamepad.ReleaseButton(vigem.XUSB_GAMEPAD_DPAD_LEFT)
    gamepad.RightTriggerFloat(0.0)
    gamepad.RightJoystickFloat(0.0, 0.0)

    gamepad.Update()

    time.Sleep(1 * time.Second)

    // Reset gamepad to default state
    gamepad.Reset()
    gamepad.Update()

    time.Sleep(1 * time.Second)
}
```

### DualShock4 gamepad

Using a virtual DS4 gamepad is similar to X360:

```go
import "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"

gamepad, err := vgamepad.NewVDS4Gamepad()
if err != nil {
    // Handle error
}
defer gamepad.Close()
```

Press and release buttons:

```go
import (
    "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"
    "github.com/CB2Moon/vgamepad-go/internal/vigem"
)

gamepad.PressButton(vigem.DS4_BUTTON_TRIANGLE)
gamepad.Update()

// (...)

gamepad.ReleaseButton(vigem.DS4_BUTTON_TRIANGLE)
gamepad.Update()
```

Available buttons are defined in the `vigem` package as `DS4Button` constants.

Press and release special buttons:

```go
gamepad.PressSpecialButton(vigem.DS4_SPECIAL_BUTTON_PS)
gamepad.Update()

// (...)

gamepad.ReleaseSpecialButton(vigem.DS4_SPECIAL_BUTTON_PS)
gamepad.Update()
```

Special buttons are defined in the `vigem` package as `DS4SpecialButton` constants.

Triggers and joysticks (integer values):

```go
// Left trigger: value between 0 and 255
gamepad.LeftTrigger(100)
// Right trigger: value between 0 and 255
gamepad.RightTrigger(255)
// Left joystick: value between 0 and 255 (128 is center)
gamepad.LeftJoystick(0, 128)
// Right joystick: value between 0 and 255 (128 is center)
gamepad.RightJoystick(0, 255)

gamepad.Update()
```

Triggers and joysticks (float values):

```go
// Left trigger: value between 0.0 and 1.0
gamepad.LeftTriggerFloat(0.5)
// Right trigger: value between 0.0 and 1.0
gamepad.RightTriggerFloat(1.0)
// Left joystick: values between -1.0 and 1.0
gamepad.LeftJoystickFloat(-0.5, 0.0)
// Right joystick: values between -1.0 and 1.0
gamepad.RightJoystickFloat(-1.0, 0.8)

gamepad.Update()
```

Directional pad (hat):

```go
gamepad.DirectionalPad(vigem.DS4_BUTTON_DPAD_NORTHWEST)
gamepad.Update()
```

Directions for the directional pad are defined in the `vigem` package as `DS4DPadDirection` constants.

Reset to default state:

```go
gamepad.Reset()
gamepad.Update()
```

Full example:

```go
package main

import (
    "time"

    "github.com/CB2Moon/vgamepad-go/pkg/vgamepad"
    "github.com/CB2Moon/vgamepad-go/internal/vigem"
)

func main() {
    gamepad, err := vgamepad.NewVDS4Gamepad()
    if err != nil {
        panic(err)
    }
    defer gamepad.Close()

    // Press a button to wake the device up
    gamepad.PressButton(vigem.DS4_BUTTON_TRIANGLE)
    gamepad.Update()
    time.Sleep(500 * time.Millisecond)
    gamepad.ReleaseButton(vigem.DS4_BUTTON_TRIANGLE)
    gamepad.Update()
    time.Sleep(500 * time.Millisecond)

    // Press buttons and things
    gamepad.PressButton(vigem.DS4_BUTTON_TRIANGLE)
    gamepad.PressButton(vigem.DS4_BUTTON_CIRCLE)
    gamepad.PressButton(vigem.DS4_BUTTON_THUMB_RIGHT)
    gamepad.PressButton(vigem.DS4_BUTTON_TRIGGER_LEFT)
    gamepad.PressSpecialButton(vigem.DS4_SPECIAL_BUTTON_TOUCHPAD)
    gamepad.LeftTriggerFloat(0.5)
    gamepad.RightTriggerFloat(0.5)
    gamepad.LeftJoystickFloat(0.0, 0.2)
    gamepad.RightJoystickFloat(-1.0, 1.0)

    gamepad.Update()

    time.Sleep(1 * time.Second)

    // Release buttons and things
    gamepad.ReleaseButton(vigem.DS4_BUTTON_TRIANGLE)
    gamepad.RightTriggerFloat(0.0)
    gamepad.RightJoystickFloat(0.0, 0.0)

    gamepad.Update()

    time.Sleep(1 * time.Second)

    // Reset gamepad to default state
    gamepad.Reset()
    gamepad.Update()

    time.Sleep(1 * time.Second)
}
```

### Rumble and LEDs:

`vgamepad-go` enables registering custom callback functions to handle updates of the rumble motors and the LED ring.

Custom callback functions require the following signature:

```go
func myCallback(client, target uintptr, largeMotor, smallMotor, ledNumber uint8, userData uintptr) {
    // Do your things here. For instance:
    fmt.Printf("Received notification for client %v, target %v\n", client, target)
    fmt.Printf("large motor: %d, small motor: %d\n", largeMotor, smallMotor)
    fmt.Printf("led number: %d\n", ledNumber)
}
```

The callback function needs to be registered as follows:

```go
err := gamepad.RegisterNotification(myCallback)
if err != nil {
    // Handle error
}
```

Each time the state of the gamepad is changed (for example by a video game that sends rumbling requests), the callback function will be called.

If not needed anymore, the callback function can be unregistered:

```go
gamepad.UnregisterNotification()
```
