package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
	"golang.org/x/sys/windows"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var asm *wca.IAudioSessionManager
	if err = mmd.Activate(wca.IID_IAudioSessionManager, wca.CLSCTX_ALL, nil, &asm); err != nil {
		return
	}
	defer asm.Release()

	/*

		var asc *wca.IAudioSessionControl
		if err = asm.GetAudioSessionControl(nil, 0, &asc); err != nil {
			return
		}
		defer asc.Release()

		var displayName string
		if err = asc.GetDisplayName(&displayName); err != nil {
			return
		}
		fmt.Println(displayName)


	*/

	var asm2 *wca.IAudioSessionManager2
	if err = mmd.Activate(wca.IID_IAudioSessionManager2, wca.CLSCTX_ALL, nil, &asm2); err != nil {
		return
	}
	defer asm2.Release()

	var sessionEnumerator *wca.IAudioSessionEnumerator
	if err = asm2.GetSessionEnumerator(&sessionEnumerator); err != nil {
		return
	}
	defer sessionEnumerator.Release()

	var cnt int
	if err = sessionEnumerator.GetCount(&cnt); err != nil {
		return
	}

	fmt.Println(cnt)

	for i := 0; i < cnt; i++ {
		var session *wca.IAudioSessionControl
		if err = sessionEnumerator.GetSession(i, &session); err != nil {
			return
		}
		defer session.Release()
		var _ids *ole.IDispatch
		var acs2 *wca.IAudioSessionControl2
		if _ids, err = session.QueryInterface(wca.IID_IAudioSessionControl2); err != nil {
			return
		}
		acs2 = (*wca.IAudioSessionControl2)(unsafe.Pointer(_ids))
		var processId uint32
		acs2.GetProcessId(&processId)
		var name string
		acs2.GetDisplayName(&name)
		fmt.Println(processId, name)
		process, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, processId)
		if err != nil {
			return err
		}
		reflect.TypeOf(process)
		var buff []uint16 = make([]uint16, 255)
		var buffPtr *uint16 = &buff[0]
		var buffSize uint32 = uint32(len(buff))
		res := windows.QueryFullProcessImageName(process, 0, buffPtr, &buffSize)
		if res != nil {
			return res
		}
		fmt.Println(windows.UTF16ToString(buff))
	}

	return
}
