// +build windows

package disk

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	RootWindows = "C:\\"
)

var (
	ErrDllLoad   = "Kernel32.dll not found"
	ErrCalcSpace = "Error calculating disk space"
	ErrCalcFree  = "Error calculating free disk space"
)

// Space returns total and free bytes available in a directory`C:\`.
func Space() (int64, int64, error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return -1, -1, errors.Wrap(err, ErrDllLoad)
	}
	defer syscall.FreeLibrary(kernel32)

	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return -1, -1, errors.Wrap(err, ErrCalcSpace)
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	r1, _, e1 := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(RootWindows))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return -1, -1, errors.Wrap(err, ErrCalcFree)
	}
	return lpTotalNumberOfBytes, lpFreeBytesAvailable, nil
}
