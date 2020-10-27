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
	message string
	taskOut shell.OutputWriter
}

// NewKoolTask creates a new kool task
func NewKoolTask(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{service, message, shell.NewOutputWriter()}
}

// Run runs task
func (t *DefaultKoolTask) Run(args []string) (err error) {
	if !t.IsTerminal() {
		return t.Execute(args)
	}

	t.taskOut.SetWriter(t.GetWriter())
	pipeReader, pipeWriter := io.Pipe()

	t.SetWriter(pipeWriter)
	defer t.SetWriter(t.taskOut.GetWriter())

	startMessage := fmt.Sprintf("%s ...", t.message)
	t.taskOut.Println(startMessage)
	t.taskOut.Println(strings.Repeat("=", len(startMessage)))

	lines := make(chan string)

	doneScanning := startServiceOutputScanner(pipeReader, lines)
	donePrinting := t.startServiceOutputPrinter(lines, doneScanning)

	err = <-t.execService(args)
	pipeWriter.Close()

	<-doneScanning
	<-donePrinting

	var statusMessage string
	if err != nil {
		statusMessage = fmt.Sprintf("... %s", color.New(color.Red).Sprint("error"))
	} else {
		statusMessage = fmt.Sprintf("... %s", color.New(color.Green).Sprint("done"))
	}

	t.taskOut.Printf("\r")
	t.taskOut.Println(statusMessage)

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

func startServiceOutputScanner(reader io.Reader, lines chan string) <-chan bool {
	doneScanning := make(chan bool)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	go func() {
		defer func() {
			close(lines)
			close(doneScanning)
		}()

		for range time.Tick(100 * time.Millisecond) {
			if !scanner.Scan() {
				break
			}

			lines <- scanner.Text()
		}

		doneScanning <- true
	}()

	return doneScanning
}

func (t *DefaultKoolTask) startServiceOutputPrinter(lines chan string, doneScanning <-chan bool) <-chan bool {
	donePrinting := make(chan bool)
	spinChars := []byte{'-', '/', '|', '\\'}
	spinPos := 0

	go func() {
		defer close(donePrinting)

	OutputPrint:
		for range time.Tick(100 * time.Millisecond) {
			spinPos = (spinPos + 1) % 4
			currentSpin := spinChars[spinPos : spinPos+1]

			select {
			case <-doneScanning:
				t.taskOut.Printf("\r")

				// remaining lines
				for line := range lines {
					if line != "" {
						t.taskOut.Println(">", line)
					}
				}

				t.taskOut.Printf("... %s", currentSpin)
				break OutputPrint
			case line := <-lines:
				spinPos = (spinPos + 1) % 4

				t.taskOut.Printf("\r")

				if line != "" {
					t.taskOut.Println(">", line)
				}

				t.taskOut.Printf("... %s", currentSpin)
			default:
				t.taskOut.Printf("\r... %s", currentSpin)
			}
		}

		donePrinting <- true
	}()

	return donePrinting
}
