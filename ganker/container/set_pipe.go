package container

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// newPipe create a pipe
func newPipe() (*os.File, *os.File, error) {
	// create a pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return readPipe, writePipe, nil
}

// sendInitCommand send command to child process
func sendInitCommand(comArray []string, writePipe *os.File) error {
	command := strings.Join(comArray, " ")
	if _, err := writePipe.WriteString(command); err != nil {
		return fmt.Errorf("write pipe error %v", err)
	}
	writePipe.Close()
	return nil
}

// readParentProcess read command from pipe
func readParentProcess() []string {

	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
