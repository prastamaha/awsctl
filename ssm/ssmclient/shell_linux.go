//go:build linux

package ssmclient

import (
	"golang.org/x/sys/unix"
	"os"
)

func cleanup() error {
	if origTermios != nil {
		return unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETSF, origTermios)
	}
	return nil
}

func configureStdin() (err error) {
	origTermios, err = unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	if err != nil {
		return err
	}

	newTermios := *origTermios
	newTermios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	newTermios.Oflag &^= unix.OPOST
	newTermios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	newTermios.Cflag &^= unix.CSIZE | unix.PARENB
	newTermios.Cflag |= unix.CS8
	newTermios.Cc[unix.VMIN] = 1
	newTermios.Cc[unix.VTIME] = 0

	return unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETSF, &newTermios)
}
