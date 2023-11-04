package main

import (
	"log/slog"
	"os"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func main() {

	if err := run(os.Args); err != nil {
		slog.Error("Fatal:", err)
	}
}

func run(args []string) (err error) {
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var sessionsControls *[]*wca.IAudioSessionControl
	if sessionsControls, err = GetAudioSessionControls(); err != nil {
		return
	}

	for _, session := range *sessionsControls {
		defer session.Release()

		var asc2 *wca.IAudioSessionControl2
		if asc2, err = GetAudioSessionControl2(session); err != nil {
			slog.Warn("", err)
			continue
		}

		var processId uint32
		if err = asc2.GetProcessId(&processId); err != nil {
			slog.Warn("", err)
			continue
		}

		if _, err = GetIconFromPid(processId); err != nil {
			slog.Warn("", err)
			continue
		}
	}

	return nil
}
