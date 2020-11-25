package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kool-dev/kool/cmd/shell"
	"strings"
	"testing"

	"github.com/gookit/color"
)

type koolTaskServiceTest struct {
	DefaultKoolService
	CalledExecute bool
	MockError     error
	MockOutput    string
}

func (t *koolTaskServiceTest) Execute(args []string) error {
	t.CalledExecute = true

	if t.MockOutput != "" {
		t.Println(t.MockOutput)
	}

	return t.MockError
}

func newKoolServiceTest() *DefaultKoolService {
	service := &DefaultKoolService{
		&shell.FakeExiter{},
		shell.NewOutputWriter(),
		&shell.FakeInputReader{},
		&shell.FakeTerminalChecker{MockIsTerminal: true},
		&shell.FakeShell{},
	}
	buf := bytes.NewBufferString("")
	service.SetWriter(buf)
	return service
}

func newKoolTaskServiceTest() *koolTaskServiceTest {
	return &koolTaskServiceTest{
		*newKoolServiceTest(),
		false,
		nil,
		"",
	}
}

func newKoolTaskServiceTestWithOutput() *koolTaskServiceTest {
	return &koolTaskServiceTest{
		*newKoolServiceTest(),
		false,
		nil,
		"testing output",
	}
}

func newKoolTaskTest(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{service, message, &shell.FakeOutputWriter{}}
}

func TestNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	task := NewKoolTask("testing", service)

	message := task.message

	if message != "testing" {
		t.Errorf("expecting message 'testing' on KoolTask, got '%s'", message)
	}

	if _, ok := task.taskOut.(*shell.DefaultOutputWriter); !ok {
		t.Error("unexpected shell.OutputWriter on KoolTask.taskOut")
	}
}

func TestRunNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	task := newKoolTaskTest("testing", service)

	_ = task.Run([]string{})

	if !service.term.(*shell.FakeTerminalChecker).CalledIsTerminal {
		t.Error("did not call IsTerminal on task KoolService")
	}

	if !service.CalledExecute {
		t.Error("did not call Execute on task KoolService")
	}

	outputLines := task.taskOut.(*shell.FakeOutputWriter).OutLines

	if len(outputLines) >= 1 && outputLines[0] != "testing ..." {
		t.Errorf("expecting message 'testing ...', got %s", outputLines[0])
	}

	expected := fmt.Sprintf("... %s", color.New(color.Green).Sprint("done"))
	if len(outputLines) >= 3 && outputLines[2] != expected {
		t.Errorf("expecting task status '%s', got %s", expected, outputLines[2])
	}
}

func TestRunFailingNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.MockError = errors.New("error execute")
	task := newKoolTaskTest("testing", service)

	err := task.Run([]string{})

	if err == nil {
		t.Errorf("expecting Run to return an error, got none")
	} else if err.Error() != service.MockError.Error() {
		t.Errorf("expecting Run to return the error '%v', got '%v'", service.MockError, err)
	}

	outputLines := task.taskOut.(*shell.FakeOutputWriter).OutLines

	expected := fmt.Sprintf("... %s", color.New(color.Red).Sprint("error"))
	if len(outputLines) >= 3 && outputLines[2] != expected {
		t.Errorf("expecting task status '%s', got %s", expected, outputLines[2])
	}
}

func TestRunNonTtyNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.term.(*shell.FakeTerminalChecker).MockIsTerminal = false
	task := newKoolTaskTest("testing", service)

	_ = task.Run([]string{})

	if outputLines := task.taskOut.(*shell.FakeOutputWriter).OutLines; len(outputLines) > 0 {
		t.Error("should not print out task output")
	}
}

func TestRunOutputNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTestWithOutput()
	task := newKoolTaskTest("testing", service)
	task.taskOut = shell.NewOutputWriter()

	_ = task.Run([]string{})

	bufBytes, err := ioutil.ReadAll(task.taskOut.GetWriter().(io.Reader))

	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(string(bufBytes))

	if !strings.Contains(output, "testing output") {
		t.Error("did not printed KoolService output")
	}
}
