package cmd

import (
	"bufio"
	"fmt"
	"io"
	"kool-dev/kool/cmd/shell"
	"strings"
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
	message   string
	taskShell shell.Shell
}

// NewKoolTask creates a new kool task
func NewKoolTask(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{service, message, shell.NewShell()}
}

// Run runs task
func (t *DefaultKoolTask) Run(args []string) (err error) {
	if !t.IsTerminal() {
		return t.Execute(args)
	}

	originalOutput := t.OutStream()
	t.taskShell.SetOutStream(originalOutput)
	pipeReader, pipeWriter := io.Pipe()

	t.SetOutStream(pipeWriter)
	defer t.SetOutStream(originalOutput)

	startMessage := fmt.Sprintf("%s ...", t.message)
	t.taskShell.Println(startMessage)
	t.taskShell.Println(strings.Repeat("=", len(startMessage)))

	lines := make(chan string)

	readServiceOutput(pipeReader, lines)
	donePrinting := t.printServiceOutput(lines)

	err = <-t.execService(args)
	pipeWriter.Close()
	<-donePrinting

	var statusMessage string
	if err != nil {
		statusMessage = fmt.Sprintf("... %s", color.New(color.Red).Sprint("error"))
	} else {
		statusMessage = fmt.Sprintf("... %s", color.New(color.Green).Sprint("done"))
	}

	t.taskShell.Printf("\r")
	t.taskShell.Println(statusMessage)

	return
}

func (t *DefaultKoolTask) execService(args []string) <-chan error {
	err := make(chan error)

	go func() {
		defer close(err)
		err <- t.Execute(args)
	}()

	return err
}

func readServiceOutput(reader io.Reader, lines chan string) {
	bufReader := bufio.NewReader(reader)

	go func() {
		defer func() {
			close(lines)
		}()

		var (
			line string
			err  error
		)

		for err == nil {
			if line, err = bufReader.ReadString('\n'); line != "" {
				lines <- strings.TrimSpace(line)
			}
		}
	}()
}

func (t *DefaultKoolTask) printServiceOutput(lines chan string) <-chan bool {
	donePrinting := make(chan bool)
	spinChars := []byte{'-', '/', '|', '\\'}
	spinPos := 0
	currentSpin := spinChars[spinPos : spinPos+1]

	go func() {
		defer close(donePrinting)

	OutputPrint:
		for {
			select {
			case line, ok := <-lines:
				if ok {
					t.taskShell.Printf("\r")
					t.taskShell.Println(">", line)
					t.taskShell.Printf("... %s", currentSpin)
				} else {
					t.taskShell.Printf("\r")
					t.taskShell.Printf("... %s", currentSpin)
					break OutputPrint
				}
			case <-time.After(100 * time.Millisecond):
				spinPos = (spinPos + 1) % 4
				currentSpin = spinChars[spinPos : spinPos+1]
				t.taskShell.Printf("\r... %s", currentSpin)
			}
		}

		donePrinting <- true
	}()

	return donePrinting
}
