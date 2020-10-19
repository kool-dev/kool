package task

import (
	"kool-dev/kool/cmd/shell"
	"sync"

	"github.com/gookit/color"
)

// Runner holds logic for running tasks
type Runner interface {
	Run(string, func() error) error
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
func (r *DefaultRunner) Run(message string, closure func() error) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)

	r.out.Printf("%s ... ", message)
	chError := make(chan error)

	go func() {
		chError <- closure()
		wg.Done()
	}()

	err = <-chError
	wg.Wait()
	close(chError)

	if err != nil {
		errorOutput := color.New(color.Red).Sprint("error")
		r.out.Println(errorOutput)
	} else {
		r.out.Success("done")
	}

	return
}
