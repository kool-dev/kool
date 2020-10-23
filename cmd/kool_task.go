package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
)

// KoolTask holds logic for running kool service as a long task
type KoolTask interface {
	KoolService
	Run([]string) error
}

// DefaultKoolTask holds data for running kool service as a long task
type DefaultKoolTask struct {
	KoolService
	message string
}

// NewKoolTask creates a new kool task
func NewKoolTask(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{service, message}
}

// Run runs task
func (t *DefaultKoolTask) Run(args []string) (err error) {
	if !t.IsTerminal() {
		return t.Execute(args)
	}

	t.Println(t.message, "...")

	originalWriter := t.GetWriter()

	_ = bytes.NewBuffer([]byte{})

	r, w, _ := os.Pipe()
	t.SetWriter(w)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	done := make(chan bool)
	lines := make(chan string)

	go func() {
		defer close(lines)

		for range time.Tick(100 * time.Millisecond) {
			select {
			case <-done:
				return
			default:
				if scanner.Scan() {
					lines <- scanner.Text()
				}
			}
		}
	}()

	go func() {
		spinChars := []byte{'-', '/', '|', '\\'}
		spinPos := 0

		for range time.Tick(100 * time.Millisecond) {
			select {
			case <-done:
				return
			case line := <-lines:
				spinPos = (spinPos + 1) % 4

				fmt.Fprint(originalWriter, "\r")

				fmt.Fprintln(originalWriter, line)

				fmt.Fprintf(originalWriter, "Status: %s", spinChars[spinPos:spinPos+1])
			default:
				spinPos = (spinPos + 1) % 4
				fmt.Fprintf(originalWriter, "\rStatus: %s", spinChars[spinPos:spinPos+1])
			}
		}
	}()

	err = <-t.execService(args)

	done <- true
	close(done)

	t.SetWriter(originalWriter)

	if err != nil {
		t.Printf("\rStatus: %s\n", color.New(color.Red).Sprint("error"))
	} else {
		t.Printf("\rStatus: %s\n", color.New(color.Green).Sprint("done"))
	}

	return
}

func (t *DefaultKoolTask) execService(args []string) <-chan error {
	chError := make(chan error)

	go func() {
		defer func() {
			close(chError)
		}()
		chError <- t.Execute(args)
	}()

	return chError
}
