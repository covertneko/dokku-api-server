package dokku

import (
	"fmt"
	"log"
	"time"
	"strings"
	"os/exec"
	"bufio"
)

func Exec(args... string) (output <-chan string, err error) {
	cmd := exec.Command("dokku", args...)
	out := make(chan string, 1)
	cmdString := strings.Join(cmd.Args, " ")

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("Could not connect to Dokku: %s", err)
		return
	}

	// Send output line-by-line through output channel
	scanner := bufio.NewScanner(cmdOut)
	go func() {
		defer close(out)
		log.Printf("Sending command output.")
		for scanner.Scan() {
			select {
			case out <- scanner.Text():
			case <-time.After(time.Second * 3):
				log.Printf("Sending output of command %q timed out.", cmdString)
				return
			}
		}
	} ()

	if err = cmd.Start(); err != nil {
		err = fmt.Errorf("Error while executing %q: %s", cmdString, err)
		return
	}

	if err = cmd.Wait(); err != nil {
		err = fmt.Errorf("Error while executing %q: %s", cmdString, err)
		return
	}

	output = out
	return
}

