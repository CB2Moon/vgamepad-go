package vgamepad

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/CB2Moon/vgamepad-go/internal/vigem"
	"github.com/CB2Moon/vgamepad-go/pkg/commons"
)

// NotificationCallback is the function signature for notification callbacks
// client: vigem bus ID
// target: vigem device ID
// largeMotor: integer in [0, 255] representing the state of the large motor
// smallMotor: integer in [0, 255] representing the state of the small motor
// ledNumber: integer in [0, 255] representing the state of the LED ring
// userData: placeholder, do not use
type NotificationCallback func(client, target uintptr, largeMotor, smallMotor, ledNumber uint8, userData uintptr)

// VBus represents a virtual USB bus (ViGEmBus)
type VBus struct {
	client *vigem.ViGEmClient
	busp   uintptr
	mu     sync.Mutex
}

var (
	// Global VBus instance for all controllers
	globalVBus     *VBus
	globalVBusOnce sync.Once
)

// GetVBus returns the global VBus instance (singleton)
func GetVBus() (*VBus, error) {
	var err error
	globalVBusOnce.Do(func() {
		globalVBus, err = newVBus()
	})
	if err != nil {
		return nil, err
	}
	return globalVBus, nil
}

// newVBus creates a new VBus instance
func newVBus() (*VBus, error) {
	client, err := vigem.NewViGEmClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create ViGEmClient: %w", err)
	}

	busp, err := client.Alloc()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate ViGEm bus: %w", err)
	}

	err = client.Connect(busp)
	if err != nil {
		client.Free(busp)
		return nil, fmt.Errorf("failed to connect to ViGEm bus: %w", err)
	}

	return &VBus{
		client: client,
		busp:   busp,
	}, nil
}

// Close closes the VBus
func (v *VBus) Close() {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.busp != 0 {
		v.client.Disconnect(v.busp)
		v.client.Free(v.busp)
		v.busp = 0
	}
}

// Gamepad is the interface for all gamepad types
type Gamepad interface {
	// Update sends the current report to the virtual device
	Update() error

	// Reset resets the gamepad to default state
	Reset()

	// Close closes the gamepad and removes it from the bus
	Close()

	// GetVID returns the vendor ID of the virtual device
	GetVID() uint16

	// GetPID returns the product ID of the virtual device
	GetPID() uint16

	// SetVID sets the vendor ID of the virtual device
	SetVID(vid uint16)

	// SetPID sets the product ID of the virtual device
	SetPID(pid uint16)

	// GetIndex returns the internally used index of the target device
	GetIndex() uint32

	// GetType returns the type of the object
	GetType() commons.ViGEmTargetType

	// RegisterNotification registers a callback function for notifications
	RegisterNotification(callback NotificationCallback) error

	// UnregisterNotification unregisters a previously registered callback function
	UnregisterNotification()
}

// BaseGamepad contains common functionality for all gamepad types
type BaseGamepad struct {
	vbus    *VBus
	client  *vigem.ViGEmClient
	busp    uintptr
	devicep uintptr
	cmpFunc unsafe.Pointer // Keep reference to callback function
}

// NewBaseGamepad creates a new BaseGamepad
func NewBaseGamepad(targetAlloc func() (uintptr, error)) (*BaseGamepad, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("vgamepad is only supported on Windows")
	}

	vbus, err := GetVBus()
	if err != nil {
		return nil, err
	}

	devicep, err := targetAlloc()
	if err != nil {
		return nil, err
	}

	err = vbus.client.TargetAdd(vbus.busp, devicep)
	if err != nil {
		vbus.client.TargetFree(devicep)
		return nil, err
	}

	if !vbus.client.TargetIsAttached(devicep) {
		vbus.client.TargetFree(devicep)
		return nil, fmt.Errorf("the virtual device could not connect to ViGEmBus")
	}

	return &BaseGamepad{
		vbus:    vbus,
		client:  vbus.client,
		busp:    vbus.busp,
		devicep: devicep,
	}, nil
}

// Close closes the gamepad and removes it from the bus
func (g *BaseGamepad) Close() {
	if g.devicep != 0 {
		g.client.TargetRemove(g.busp, g.devicep)
		g.client.TargetFree(g.devicep)
		g.devicep = 0
	}
}

// GetVID returns the vendor ID of the virtual device
func (g *BaseGamepad) GetVID() uint16 {
	return g.client.TargetGetVid(g.devicep)
}

// GetPID returns the product ID of the virtual device
func (g *BaseGamepad) GetPID() uint16 {
	return g.client.TargetGetPid(g.devicep)
}

// SetVID sets the vendor ID of the virtual device
func (g *BaseGamepad) SetVID(vid uint16) {
	g.client.TargetSetVid(g.devicep, vid)
}

// SetPID sets the product ID of the virtual device
func (g *BaseGamepad) SetPID(pid uint16) {
	g.client.TargetSetPid(g.devicep, pid)
}

// GetIndex returns the internally used index of the target device
func (g *BaseGamepad) GetIndex() uint32 {
	return g.client.TargetGetIndex(g.devicep)
}

// GetType returns the type of the object
func (g *BaseGamepad) GetType() commons.ViGEmTargetType {
	return g.client.TargetGetType(g.devicep)
}
