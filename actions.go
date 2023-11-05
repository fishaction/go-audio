package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func GetAllAudioProcesses(c *gin.Context) {
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		slog.Warn("", err)
	}
	defer ole.CoUninitialize()

	arr := []gin.H{}
	var err error
	var sessionsControls *[]*wca.IAudioSessionControl
	if sessionsControls, err = GetAudioSessionControls(); err != nil {
		slog.Warn("", err)
	}

	for _, session := range *sessionsControls {
		defer session.Release()

		var asc2 *wca.IAudioSessionControl2
		if asc2, err = GetAudioSessionControl2(session); err != nil {
			slog.Warn("", err)
			continue
		}

		err = nil
		var processId uint32
		if err = asc2.GetProcessId(&processId); err != nil {
			slog.Warn("", err)
			continue
		}

		arr = append(arr, gin.H{
			"processId": processId,
		})

	}

	c.JSON(200, gin.H{
		"audioProcesses": arr,
	})
}
