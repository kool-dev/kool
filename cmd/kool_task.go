package cmd

import (
	"bytes"
	"github.com/gookit/color"
	"io/ioutil"
	"strings"
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

	var output string

	t.Printf("%s ... ", t.message)

	originalWriter := t.GetWriter()
	buf := bytes.NewBufferString("")
	t.SetWriter(buf)

	err = <-t.execService(args)

	t.SetWriter(originalWriter)
	bufBytes, _ := ioutil.ReadAll(buf)
	output = strings.TrimSpace(string(bufBytes))

	if err != nil {
		t.Println(color.New(color.Red).Sprint("error"))
	} else {
		t.Println(color.New(color.Green).Sprint("done"))
	}

	if output != "" {
		t.Println(output)
	}

	return
}

func (t *DefaultKoolTask) execService(args []string) <-chan error {
	chError := make(chan error)

	go func() {
		defer close(chError)
		chError <- t.Execute(args)
	}()

	return chError
}
