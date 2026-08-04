package ssmclient

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mmmorris1975/ssm-session-client/datachannel"
)

func ShellSession(cfg aws.Config, target string, initCmd ...io.Reader) error {
	c := new(datachannel.SsmDataChannel)
	if err := c.Open(cfg, &ssm.StartSessionInput{Target: aws.String(target)}); err != nil {
		return err
	}
	defer c.Close()

	if err := initialize(c); err != nil {
		return err
	}
	defer cleanup() //nolint:errcheck

	errCh := make(chan error, 5)
	go func() {
		if _, err := io.Copy(c, os.Stdin); err != nil {
			errCh <- err
		}
	}()

	for _, cmd := range initCmd {
		_, _ = io.Copy(c, cmd)
	}

	if _, err := io.Copy(os.Stdout, c); err != nil {
		if !errors.Is(err, io.EOF) {
			errCh <- err
		}
	}
	close(errCh)

	return <-errCh
}

func updateTermSize(c datachannel.DataChannel) error {
	rows, cols, err := getWinSize()
	if err != nil {
		cols = 132
		rows = 45
		log.Printf("Could not get size of the terminal: %s, using width %d height %d\n", err, cols, rows)
	}

	return c.SetTerminalSize(rows, cols)
}
