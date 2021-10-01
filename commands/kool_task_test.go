package commands

import (
	"bytes"
	"errors"
	"fmt"
	"kool-dev/kool/core/shell"
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
		t.Shell().Println(t.MockOutput)
	}

	return t.MockError
}

func newKoolServiceTest() *DefaultKoolService {
	service := &DefaultKoolService{
		shell.NewShell(),
	}
	buf := bytes.NewBufferString("")
	service.Shell().SetOutStream(buf)
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

func newKoolTaskTest(message string, service KoolService) *DefaultKoolTask {
	return &DefaultKoolTask{
		KoolService: service,
		message:     message,
		actualOut:   &shell.FakeShell{},
		frameOutput: true,
	}
}

func TestNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	task := NewKoolTask("testing", service)

	message := task.message

	if message != "testing" {
		t.Errorf("expecting message 'testing' on KoolTask, got '%s'", message)
	}

	if _, ok := task.actualOut.(*shell.DefaultShell); !ok {
		t.Error("unexpected shell.Shell on KoolTask.actualOut")
	}
}

func TestRunNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.Fake()
	service.shell.(*shell.FakeShell).MockIsTerminal = false
	task := newKoolTaskTest("testing", service)

	_ = task.Run([]string{})

	if !service.shell.(*shell.FakeShell).CalledIsTerminal {
		t.Error("did not call IsTerminal on task KoolService")
	}

	if !service.CalledExecute {
		t.Error("did not call Execute on task KoolService")
	}

	outputLines := task.actualOut.(*shell.FakeShell).OutLines

	if len(outputLines) >= 1 && !strings.HasSuffix(outputLines[0], " testing") {
		t.Errorf("expecting message '[done] testing', got %s", outputLines[0])
	}

	task = newKoolTaskTest("testing", service)
	task.SetFrameOutput(false)

	_ = task.Run([]string{})

	if !service.shell.(*shell.FakeShell).CalledIsTerminal {
		t.Error("did not call IsTerminal on task KoolService")
	}

	if !service.CalledExecute {
		t.Error("did not call Execute on task KoolService")
	}

	outputLines = task.actualOut.(*shell.FakeShell).OutLines

	if len(outputLines) >= 1 && !strings.HasSuffix(outputLines[0], " testing") {
		t.Errorf("expecting message '[done] testing', got %s", outputLines[0])
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

	outputLines := task.actualOut.(*shell.FakeShell).OutLines

	expected := fmt.Sprintf("... %s", color.New(color.Red).Sprint("error"))
	if len(outputLines) >= 3 && outputLines[2] != expected {
		t.Errorf("expecting task status '%s', got %s", expected, outputLines[2])
	}
}

func TestRunNonTtyNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.Fake()
	service.shell.(*shell.FakeShell).MockIsTerminal = false
	service.shell.(*shell.FakeShell).MockOutStream = bytes.NewBufferString("")
	task := newKoolTaskTest("testing", service)

	_ = task.Run([]string{})

	if outputLines := task.actualOut.(*shell.FakeShell).OutLines; len(outputLines) > 0 {
		t.Error("should not print out task output")
	}
}

func TestRunOutputNewKoolTask(t *testing.T) {
	service := &koolTaskServiceTest{
		*newKoolServiceTest(),
		false,
		nil,
		"testing output",
	}

	service.Fake()
	service.shell.(*shell.FakeShell).MockIsTerminal = false
	task := newKoolTaskTest("testing", service)
	task.actualOut = service.Shell()

	_ = task.Run([]string{})

	output := strings.Join(service.shell.(*shell.FakeShell).OutLines, "\n")

	if !strings.Contains(output, "testing output") {
		t.Error("did not print KoolService output")
	}
}
