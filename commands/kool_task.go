package commands

import (
	"bufio"
	"fmt"
	"io"
	"kool-dev/kool/core/shell"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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
	message     string
	actualOut   shell.Shell
	frameOutput bool
	originalOut io.Writer
}

// NewKoolTask creates a new kool task
func NewKoolTask(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{
		KoolService: service,
		message:     message,
		actualOut:   shell.NewShell(),
		frameOutput: true,
	}
}

// Run runs task
func (t *DefaultKoolTask) Run(args []string) (err error) {
	if !t.Shell().IsTerminal() {
		return t.Execute(args)
	}

	t.originalOut = t.Shell().OutStream()
	t.actualOut.SetOutStream(t.originalOut)
	pipeReader, pipeWriter := io.Pipe()

	t.Shell().SetOutStream(pipeWriter)
	origErr := t.Shell().ErrStream()
	t.Shell().SetErrStream(pipeWriter)
	defer t.Shell().SetOutStream(t.originalOut)
	defer t.Shell().SetErrStream(origErr)

	lines := make(chan string)

	readServiceOutput(pipeReader, lines)
	donePrinting := t.printServiceOutput(lines)

	err = <-t.execService(args)
	pipeWriter.Close()
	<-donePrinting
	var statusMessage string
	if err != nil {
		statusMessage = color.New(color.Red).Sprint("error")
	} else {
		statusMessage = color.New(color.Green).Sprint("done")
	}

	t.actualOut.Printf("\r" + strings.Repeat(" ", 100) + "\r")
	t.actualOut.Println(fmt.Sprintf("[%s] %s", statusMessage, t.message))

	return
}

// SetFrameOutput
func (t *DefaultKoolTask) SetFrameOutput(frame bool) {
	t.frameOutput = frame
}

func (t *DefaultKoolTask) execService(args []string) <-chan error {
	var err = make(chan error)

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
			if line, err = bufReader.ReadString('\n'); strings.TrimSpace(line) != "" {
				lines <- strings.TrimSpace(line)
			}
		}
	}()
}

func (t *DefaultKoolTask) printServiceOutput(lines chan string) <-chan bool {
	var (
		donePrinting = make(chan bool)
		loading      = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(t.originalOut))
	)

	loading.Prefix = " "
	loading.Suffix = " " + t.message
	loading.Start()

	go func() {
		defer close(donePrinting)

		for line := range lines {
			loading.Lock()
			t.actualOut.Printf("\r" + strings.Repeat(" ", 100) + "\r")
			if t.frameOutput {
				t.actualOut.Println(">", line)
			} else {
				t.actualOut.Println(line)
			}
			loading.Unlock()
		}

		loading.Stop()
		donePrinting <- true
	}()

	return donePrinting
}
