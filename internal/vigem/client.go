package vigem

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/CB2Moon/vgamepad-go/pkg/commons"
)

// ViGEmBus version
const VIGEMBUS_VERSION = "1.17.333.0"

// ViGEmClient represents the DLL interface
type ViGEmClient struct {
	dll                                   *syscall.DLL
	vigemAlloc                            *syscall.Proc
	vigemFree                             *syscall.Proc
	vigemConnect                          *syscall.Proc
	vigemDisconnect                       *syscall.Proc
	vigemTargetX360Alloc                  *syscall.Proc
	vigemTargetDS4Alloc                   *syscall.Proc
	vigemTargetFree                       *syscall.Proc
	vigemTargetAdd                        *syscall.Proc
	vigemTargetRemove                     *syscall.Proc
	vigemTargetSetVid                     *syscall.Proc
	vigemTargetSetPid                     *syscall.Proc
	vigemTargetGetVid                     *syscall.Proc
	vigemTargetGetPid                     *syscall.Proc
	vigemTargetX360Update                 *syscall.Proc
	vigemTargetDS4Update                  *syscall.Proc
	vigemTargetDS4UpdateExPtr             *syscall.Proc
	vigemTargetGetIndex                   *syscall.Proc
	vigemTargetGetType                    *syscall.Proc
	vigemTargetIsAttached                 *syscall.Proc
	vigemTargetX360GetUserIndex           *syscall.Proc
	vigemTargetX360RegisterNotification   *syscall.Proc
	vigemTargetX360UnregisterNotification *syscall.Proc
	vigemTargetDS4RegisterNotification    *syscall.Proc
	vigemTargetDS4UnregisterNotification  *syscall.Proc
}

// NotificationCallback is the function signature for notification callbacks
type NotificationCallback func(client, target uintptr, largeMotor, smallMotor, ledNumber uint8, userData uintptr)

// NewViGEmClient creates a new ViGEmClient
func NewViGEmClient() (*ViGEmClient, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("ViGEmClient is only supported on Windows")
	}

	// Determine architecture
	var arch string
	if runtime.GOARCH == "amd64" {
		arch = "x64"
	} else if runtime.GOARCH == "386" {
		arch = "x86"
	} else {
		return nil, fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Check if ViGEmBus is installed and install if needed
	err := ensureViGEmBusInstalled(arch)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure ViGEmBus is installed: %w", err)
	}

	// Load the DLL
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	dllPath := filepath.Join(filepath.Dir(exePath), "ViGEmClient.dll")
	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		// If DLL is not in the executable directory, use the embedded one
		dllPath = filepath.Join(".", "vigem", "client", arch, "ViGEmClient.dll")
	}

	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ViGEmClient.dll: %w", err)
	}

	client := &ViGEmClient{
		dll: dll,
	}

	// Load all the procedures
	client.vigemAlloc, err = dll.FindProc("vigem_alloc")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_alloc: %w", err)
	}

	client.vigemFree, err = dll.FindProc("vigem_free")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_free: %w", err)
	}

	client.vigemConnect, err = dll.FindProc("vigem_connect")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_connect: %w", err)
	}

	client.vigemDisconnect, err = dll.FindProc("vigem_disconnect")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_disconnect: %w", err)
	}

	client.vigemTargetX360Alloc, err = dll.FindProc("vigem_target_x360_alloc")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_x360_alloc: %w", err)
	}

	client.vigemTargetDS4Alloc, err = dll.FindProc("vigem_target_ds4_alloc")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_ds4_alloc: %w", err)
	}

	client.vigemTargetFree, err = dll.FindProc("vigem_target_free")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_free: %w", err)
	}

	client.vigemTargetAdd, err = dll.FindProc("vigem_target_add")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_add: %w", err)
	}

	client.vigemTargetRemove, err = dll.FindProc("vigem_target_remove")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_remove: %w", err)
	}

	client.vigemTargetSetVid, err = dll.FindProc("vigem_target_set_vid")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_set_vid: %w", err)
	}

	client.vigemTargetSetPid, err = dll.FindProc("vigem_target_set_pid")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_set_pid: %w", err)
	}

	client.vigemTargetGetVid, err = dll.FindProc("vigem_target_get_vid")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_get_vid: %w", err)
	}

	client.vigemTargetGetPid, err = dll.FindProc("vigem_target_get_pid")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_get_pid: %w", err)
	}

	client.vigemTargetX360Update, err = dll.FindProc("vigem_target_x360_update")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_x360_update: %w", err)
	}

	client.vigemTargetDS4Update, err = dll.FindProc("vigem_target_ds4_update")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_ds4_update: %w", err)
	}

	client.vigemTargetDS4UpdateExPtr, err = dll.FindProc("vigem_target_ds4_update_ex_ptr")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_ds4_update_ex_ptr: %w", err)
	}

	client.vigemTargetGetIndex, err = dll.FindProc("vigem_target_get_index")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_get_index: %w", err)
	}

	client.vigemTargetGetType, err = dll.FindProc("vigem_target_get_type")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_get_type: %w", err)
	}

	client.vigemTargetIsAttached, err = dll.FindProc("vigem_target_is_attached")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_is_attached: %w", err)
	}

	client.vigemTargetX360GetUserIndex, err = dll.FindProc("vigem_target_x360_get_user_index")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_x360_get_user_index: %w", err)
	}

	client.vigemTargetX360RegisterNotification, err = dll.FindProc("vigem_target_x360_register_notification")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_x360_register_notification: %w", err)
	}

	client.vigemTargetX360UnregisterNotification, err = dll.FindProc("vigem_target_x360_unregister_notification")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_x360_unregister_notification: %w", err)
	}

	client.vigemTargetDS4RegisterNotification, err = dll.FindProc("vigem_target_ds4_register_notification")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_ds4_register_notification: %w", err)
	}

	client.vigemTargetDS4UnregisterNotification, err = dll.FindProc("vigem_target_ds4_unregister_notification")
	if err != nil {
		return nil, fmt.Errorf("failed to find vigem_target_ds4_unregister_notification: %w", err)
	}

	return client, nil
}

// ensureViGEmBusInstalled checks if ViGEmBus is installed and installs it if needed
func ensureViGEmBusInstalled(arch string) error {
	// Check if ViGEmBus is installed
	cmd := exec.Command("reg", "query", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, "/s")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to query registry: %w", err)
	}

	// Check if ViGEmBus is in the registry output
	if !strings.Contains(strings.ToLower(string(output)), "nefarius virtual gamepad emulation bus driver") {
		// Install ViGEmBus
		msiPath := filepath.Join(".", "vigem", "install", arch, fmt.Sprintf("ViGEmBusSetup_%s.msi", arch))
		cmd = exec.Command("msiexec", "/i", msiPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install ViGEmBus: %w", err)
		}
	}

	return nil
}

// Alloc allocates an object representing a driver connection
func (c *ViGEmClient) Alloc() (uintptr, error) {
	ret, _, _ := c.vigemAlloc.Call()
	if ret == 0 {
		return 0, fmt.Errorf("failed to allocate ViGEm client")
	}
	return ret, nil
}

// Free frees up memory used by the driver connection object
func (c *ViGEmClient) Free(client uintptr) {
	c.vigemFree.Call(client)
}

// Connect initializes the driver object and establishes a connection to the emulation bus driver
func (c *ViGEmClient) Connect(client uintptr) error {
	ret, _, _ := c.vigemConnect.Call(client)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// Disconnect disconnects from the bus device and resets the driver object state
func (c *ViGEmClient) Disconnect(client uintptr) {
	c.vigemDisconnect.Call(client)
}

// TargetX360Alloc allocates an object representing an Xbox 360 Controller device
func (c *ViGEmClient) TargetX360Alloc() (uintptr, error) {
	ret, _, _ := c.vigemTargetX360Alloc.Call()
	if ret == 0 {
		return 0, fmt.Errorf("failed to allocate Xbox 360 target")
	}
	return ret, nil
}

// TargetDS4Alloc allocates an object representing a DualShock 4 Controller device
func (c *ViGEmClient) TargetDS4Alloc() (uintptr, error) {
	ret, _, _ := c.vigemTargetDS4Alloc.Call()
	if ret == 0 {
		return 0, fmt.Errorf("failed to allocate DualShock 4 target")
	}
	return ret, nil
}

// TargetFree frees up memory used by the target device object
func (c *ViGEmClient) TargetFree(target uintptr) {
	c.vigemTargetFree.Call(target)
}

// TargetAdd adds a provided target device to the bus driver
func (c *ViGEmClient) TargetAdd(client, target uintptr) error {
	ret, _, _ := c.vigemTargetAdd.Call(client, target)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// TargetRemove removes a provided target device from the bus driver
func (c *ViGEmClient) TargetRemove(client, target uintptr) error {
	ret, _, _ := c.vigemTargetRemove.Call(client, target)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// TargetSetVid overrides the default Vendor ID value with the provided one
func (c *ViGEmClient) TargetSetVid(target uintptr, vid uint16) {
	c.vigemTargetSetVid.Call(target, uintptr(vid))
}

// TargetSetPid overrides the default Product ID value with the provided one
func (c *ViGEmClient) TargetSetPid(target uintptr, pid uint16) {
	c.vigemTargetSetPid.Call(target, uintptr(pid))
}

// TargetGetVid returns the Vendor ID of the provided target device object
func (c *ViGEmClient) TargetGetVid(target uintptr) uint16 {
	ret, _, _ := c.vigemTargetGetVid.Call(target)
	return uint16(ret)
}

// TargetGetPid returns the Product ID of the provided target device object
func (c *ViGEmClient) TargetGetPid(target uintptr) uint16 {
	ret, _, _ := c.vigemTargetGetPid.Call(target)
	return uint16(ret)
}

// TargetX360Update sends a state report to the provided target device
func (c *ViGEmClient) TargetX360Update(client, target uintptr, report commons.XUSBReport) error {
	ret, _, _ := c.vigemTargetX360Update.Call(client, target, uintptr(unsafe.Pointer(&report)))
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// TargetDS4Update sends a state report to the provided target device
func (c *ViGEmClient) TargetDS4Update(client, target uintptr, report commons.DS4Report) error {
	ret, _, _ := c.vigemTargetDS4Update.Call(client, target, uintptr(unsafe.Pointer(&report)))
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// TargetDS4UpdateExPtr sends a full size state report to the provided target device
func (c *ViGEmClient) TargetDS4UpdateExPtr(client, target uintptr, reportPtr *commons.DS4ReportEx) error {
	ret, _, _ := c.vigemTargetDS4UpdateExPtr.Call(client, target, uintptr(unsafe.Pointer(reportPtr)))
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// TargetGetIndex returns the internal index (serial number) the bus driver assigned to the provided target device object
func (c *ViGEmClient) TargetGetIndex(target uintptr) uint32 {
	ret, _, _ := c.vigemTargetGetIndex.Call(target)
	return uint32(ret)
}

// TargetGetType returns the type of the provided target device object
func (c *ViGEmClient) TargetGetType(target uintptr) commons.ViGEmTargetType {
	ret, _, _ := c.vigemTargetGetType.Call(target)
	return commons.ViGEmTargetType(ret)
}

// TargetIsAttached returns true if the provided target device object is currently attached to the bus
func (c *ViGEmClient) TargetIsAttached(target uintptr) bool {
	ret, _, _ := c.vigemTargetIsAttached.Call(target)
	return ret != 0
}

// TargetX360GetUserIndex returns the user index of the emulated Xenon device
func (c *ViGEmClient) TargetX360GetUserIndex(client, target, index uintptr) error {
	ret, _, _ := c.vigemTargetX360GetUserIndex.Call(client, target, index)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// Registers a function which gets called, when LED index or vibration state changes
// occur on the provided target device. This function fails if the provided
// target device isn't fully operational or in an erroneous state
func (c *ViGEmClient) TargetX360RegisterNotification(client, target, notification, userData uintptr) error {
	ret, _, _ := c.vigemTargetX360RegisterNotification.Call(client, target, notification, userData)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// Removes a previously registered callback function from the provided target object
func (c *ViGEmClient) TargetX360UnregisterNotification(target uintptr) {
	c.vigemTargetX360UnregisterNotification.Call(target)
}

// Registers a function which gets called, when LightBar or vibration state changes
// occur on the provided target device. This function fails if the provided
// target device isn't fully operational or in an erroneous state
func (c *ViGEmClient) TargetDS4RegisterNotification(client, target, notification, userData uintptr) error {
	ret, _, _ := c.vigemTargetDS4RegisterNotification.Call(client, target, notification, userData)
	if commons.ViGEmError(ret) != commons.VIGEM_ERROR_NONE {
		return commons.ViGEmError(ret)
	}
	return nil
}

// Removes a previously registered callback function from the provided target object
func (c *ViGEmClient) TargetDS4UnregisterNotification(target uintptr) {
	c.vigemTargetDS4UnregisterNotification.Call(target)
}
