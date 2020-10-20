package task

import (
	"kool-dev/kool/cmd/shell"

	"github.com/gookit/color"
)

// Runner holds logic for running tasks
type Runner interface {
	Run(string, func() (interface{}, error)) (interface{}, error)
}

// DefaultRunner holds data for running tasks
type DefaultRunner struct {
	out shell.OutputWriter
}

// NewRunner creates a new task runner instance
func NewRunner() Runner {
	return &DefaultRunner{shell.NewOutputWriter()}
}

// Run runs task
func (r *DefaultRunner) Run(message string, closure func() (interface{}, error)) (result interface{}, err error) {
	r.out.Printf("%s ... ", message)

	chError := make(chan error)
	chResult := make(chan interface{})
	defer close(chResult)
	defer close(chError)

	go execClosure(closure, chResult, chError)

	result = <-chResult
	err = <-chError

	if err != nil {
		errorOutput := color.New(color.Red).Sprint("error")
		r.out.Println(errorOutput)
	} else {
		r.out.Success("done")
	}

	return
}

func execClosure(closure func() (interface{}, error), chResult chan interface{}, chError chan error) {
	result, err := closure()
	chResult <- result
	chError <- err
}
