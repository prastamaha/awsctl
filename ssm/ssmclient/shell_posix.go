//go:build !windows && !js

package ssmclient

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mmmorris1975/ssm-session-client/datachannel"
	"golang.org/x/sys/unix"
)

const (
	ResizeSleepInterval = time.Millisecond * 500
)

var origTermios *unix.Termios

func initialize(c datachannel.DataChannel) error {
	installSignalHandlers(c) <- unix.SIGWINCH
	handleTerminalResize(c)
	return configureStdin()
}

func installSignalHandlers(c datachannel.DataChannel) chan os.Signal {
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt, unix.SIGQUIT, unix.SIGTERM, unix.SIGWINCH)

	go func() {
		switch <-sigCh {
		case unix.SIGWINCH:
			_ = updateTermSize(c)
		case os.Interrupt, unix.SIGQUIT, unix.SIGTERM:
			log.Print("exiting")
			_ = cleanup()
			_ = c.Close()
			os.Exit(0)
		}
	}()

	return sigCh
}

func getWinSize() (rows, cols uint32, err error) {
	var sz *unix.Winsize

	sz, err = unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}

	return uint32(sz.Row), uint32(sz.Col), nil
}

func handleTerminalResize(c datachannel.DataChannel) {
	go func() {
		for {
			_ = updateTermSize(c)
			time.Sleep(ResizeSleepInterval)
		}
	}()
}
