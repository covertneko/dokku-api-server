package dokku

import (
	"fmt"
	"strings"
	"os/exec"
	"bufio"
)

func Exec(args... string) (output []string, err error) {
	cmd := exec.Command("dokku", args...)
	cmdString := strings.Join(cmd.Args, " ")

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("Could not connect to Dokku: %s", err)
		return
	}

	// Send output line-by-line through output channel
	scanner := bufio.NewScanner(cmdOut)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			output = append(output, scanner.Text())
		}
		done <- true
	}()

	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("Error while executing %q: %s", cmdString, err)
		return
	}

	if err = cmd.Wait(); err != nil {
		err = fmt.Errorf("Error while executing %q: %s", cmdString, err)
		return
	}

	<-done
	return
}

