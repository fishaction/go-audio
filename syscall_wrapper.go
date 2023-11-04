package main

import (
	"strconv"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
	"golang.org/x/sys/windows"
)

func GetAudioSessionControl2(asc *wca.IAudioSessionControl) (*wca.IAudioSessionControl2, error) {
	var err error
	var _ids *ole.IDispatch
	var asc2 *wca.IAudioSessionControl2
	if _ids, err = asc.QueryInterface(wca.IID_IAudioSessionControl2); err != nil {
		return nil, err
	}
	asc2 = (*wca.IAudioSessionControl2)(unsafe.Pointer(_ids))
	return asc2, nil
}

func GetIconFromPid(pid uint32) (string, error) {
	ExtractIconLibrary, err := windows.LoadDLL("./lib/ExtractIconLibrary.dll")
	ExtractIconV3, err := ExtractIconLibrary.FindProc("ExtractIconV3")
	if err != nil {
		return string([]byte{0}), err
	}

	process, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return string([]byte{0}), err
	}
	var buff []uint16 = make([]uint16, 255)
	var buffPtr *uint16 = &buff[0]
	var buffSize uint32 = uint32(len(buff))
	if err = windows.QueryFullProcessImageName(process, 0, buffPtr, &buffSize); err != nil {
		return string([]byte{0}), err
	}

	filePath := windows.UTF16ToString(buff)
	expPath := strconv.Itoa(int(pid)) + ".png"

	var resultPtr uintptr
	resultPtr, _, err = ExtractIconV3.Call(
		(uintptr)(unsafe.Pointer(windows.StringToUTF16Ptr(filePath))),
		(uintptr)(unsafe.Pointer(windows.StringToUTF16Ptr(expPath))),
		0x4,
	)
	result := (*int32)(unsafe.Pointer(&resultPtr))
	if *result != 0 {
		return string([]byte{0}), err
	}

	return strconv.Itoa(int(pid)) + ".png", nil
}

func GetAudioSessionControls() (*[]*wca.IAudioSessionControl, error) {

	var err error
	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return nil, err
	}
	defer mmd.Release()

	var asm *wca.IAudioSessionManager
	if err = mmd.Activate(wca.IID_IAudioSessionManager, wca.CLSCTX_ALL, nil, &asm); err != nil {
		return nil, err
	}
	defer asm.Release()

	var asm2 *wca.IAudioSessionManager2
	if err = mmd.Activate(wca.IID_IAudioSessionManager2, wca.CLSCTX_ALL, nil, &asm2); err != nil {
		return nil, err
	}
	defer asm2.Release()

	var sessionEnumerator *wca.IAudioSessionEnumerator
	if err = asm2.GetSessionEnumerator(&sessionEnumerator); err != nil {
		return nil, err
	}
	defer sessionEnumerator.Release()

	var cnt int
	if err = sessionEnumerator.GetCount(&cnt); err != nil {
		return nil, err
	}

	sessions := []*wca.IAudioSessionControl{}

	for i := 0; i < cnt; i++ {
		var session *wca.IAudioSessionControl
		if err = sessionEnumerator.GetSession(i, &session); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return &sessions, nil
}
