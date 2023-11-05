package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := GetRouter()
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen:", err)
		}
	}()

	if err := run(os.Args); err != nil {
		slog.Error("Fatal:", err)
	}

	<-ctx.Done()

	stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: ", err)
	}

	slog.Info("Server exiting")
}

func run(args []string) (err error) {

	/*
		if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
			slog.Error("Fatal:", err)
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
		}*/

	return nil
}
