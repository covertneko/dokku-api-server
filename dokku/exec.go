package dokku

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Output struct {
	Lines chan string
	Err   error
}

const (
	OUTPUT_CHANNEL_SIZE = 100
	DOKKU_COMMAND       = "dokku"
)

func followCmd(command string, args ...string) (*Output, error) {
	cmd := exec.Command(command, args...)
	cmdString := strings.Join(cmd.Args, " ")

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("Error constructing command %q: %s", cmdString, err)
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("Error starting %q: %s", cmdString, err)
		return nil, err
	}

	output := &Output{
		Lines: make(chan string, OUTPUT_CHANNEL_SIZE),
		Err:   nil,
	}

	// Send output line-by-line through output channel
	scanner := bufio.NewScanner(cmdOut)
	go func() {
		for scanner.Scan() {
			output.Lines <- scanner.Text()
		}

		if err = cmd.Wait(); err != nil {
			log.Printf("Error while executing %q: %s\n", cmdString, err)
		}

		close(output.Lines)
	}()

	return output, nil
}

func execCmd(command string, args ...string) (output []string, err error) {
	cmd := exec.Command(command, args...)
	cmdString := strings.Join(cmd.Args, " ")

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("Error constructing command %q: %s", cmdString, err)
		return
	}

	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("Error starting %q: %s", cmdString, err)
		return
	}

	scanner := bufio.NewScanner(cmdOut)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			output = append(output, scanner.Text())
		}
		done <- true
	}()

	if err = cmd.Wait(); err != nil {
		err = fmt.Errorf("Error while executing %q: %s", cmdString, err)
		return
	}

	<-done
	return
}

func Exec(args ...string) (output []string, err error) {
	return execCmd(DOKKU_COMMAND, args...)
}
