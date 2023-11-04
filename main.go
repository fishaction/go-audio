package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
	"golang.org/x/sys/windows"
)

func main() {
	log.SetFlags(0)

	if err := run(os.Args); err != nil {
		log.Panic(err)
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

	SHGetFileInfoW := windows.NewLazyDLL("shell32.dll").NewProc("SHGetFileInfoW")
	if SHGetFileInfoW.Find() != nil {
		return SHGetFileInfoW.Find()
	}

	BitmapFromHICON := windows.NewLazyDLL("gdiplus.dll").NewProc("GdipCreateBitmapFromHICON")
	if err := BitmapFromHICON.Find(); err != nil {
		return err
	}

	SaveImageToFile := windows.NewLazyDLL("gdiplus.dll").NewProc("GdipSaveImageToFile")
	if err := SaveImageToFile.Find(); err != nil {
		return err
	}

	ExtractIconLibrary, err := windows.LoadDLL("./lib/ExtractIconLibrary.dll")
	ExtractIconV3, err := ExtractIconLibrary.FindProc("ExtractIconV3")
	if err != nil {
		return err
	}

	for i := 0; i < cnt; i++ {
		var session *wca.IAudioSessionControl
		if err = sessionEnumerator.GetSession(i, &session); err != nil {
			return err
		}
		defer session.Release()
		var _ids *ole.IDispatch
		var acs2 *wca.IAudioSessionControl2
		if _ids, err = session.QueryInterface(wca.IID_IAudioSessionControl2); err != nil {
			return err
		}
		acs2 = (*wca.IAudioSessionControl2)(unsafe.Pointer(_ids))
		var processId uint32
		acs2.GetProcessId(&processId)
		var name string
		acs2.GetDisplayName(&name)
		process, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, processId)
		if processId == 0 {
			continue
		}
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
		filePath := windows.UTF16ToString(buff)
		expPath := strconv.Itoa(int(processId)) + ".png"

		r, r2, err := ExtractIconV3.Call(
			(uintptr)(unsafe.Pointer(windows.StringToUTF16Ptr(filePath))),
			(uintptr)(unsafe.Pointer(windows.StringToUTF16Ptr(expPath))),
			0x4,
		)
		fmt.Println("r:", r)
		fmt.Println("r2:", r2)
		fmt.Println("err:", err)

		fmt.Println("-----------")
	}

	return
}

const SHGFI_ICON uint32 = 0x000000100
const SHGFI_LARGEICON uint32 = 0x000000000
const GMEM_FIXED uint32 = 0x0000
