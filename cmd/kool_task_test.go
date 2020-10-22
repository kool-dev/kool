package cmd

import (
	"bytes"
	"errors"
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

func newKoolTaskServiceTest() *koolTaskServiceTest {
	return &koolTaskServiceTest{
		*newFakeKoolService(),
		false,
		nil,
		"",
	}
}

func newKoolTaskServiceTestWithOutput() *koolTaskServiceTest {
	baseservice := &DefaultKoolService{
		&shell.FakeExiter{},
		shell.NewOutputWriter(),
		&shell.FakeInputReader{},
		&shell.FakeTerminalChecker{MockIsTerminal: true},
	}

	return &koolTaskServiceTest{
		*baseservice,
		false,
		nil,
		"testing output",
	}
}

func TestNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	task := NewKoolTask("testing", service)

	_ = task.Run([]string{})

	if !service.term.(*shell.FakeTerminalChecker).CalledIsTerminal {
		t.Error("did not call IsTerminal on task KoolService")
	}

	if !service.CalledExecute {
		t.Error("did not call Execute on task KoolService")
	}

	fOutput := service.out.(*shell.FakeOutputWriter).FOutput

	if fOutput != "testing ... " {
		t.Errorf("expecting message 'testing ... ', got %s", fOutput)
	}

	output := service.out.(*shell.FakeOutputWriter).OutLines

	expected := color.New(color.Green).Sprint("done")
	if len(output) > 0 && output[0] != expected {
		t.Errorf("expecting task status '%s', got %s", expected, output[0])
	}
}

func TestFailingNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.MockError = errors.New("error execute")
	task := NewKoolTask("testing", service)

	err := task.Run([]string{})

	if err == nil {
		t.Errorf("expecting Run to return an error, got none")
	} else if err.Error() != service.MockError.Error() {
		t.Errorf("expecting Run to return the error '%v', got '%v'", service.MockError, err)
	}

	output := service.out.(*shell.FakeOutputWriter).OutLines

	expected := color.New(color.Red).Sprint("error")
	if len(output) > 0 && output[0] != expected {
		t.Errorf("expecting task status '%s', got %s", expected, output[0])
	}
}

func TestNonTtyNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTest()
	service.term.(*shell.FakeTerminalChecker).MockIsTerminal = false
	task := NewKoolTask("testing", service)

	_ = task.Run([]string{})

	if service.out.(*shell.FakeOutputWriter).CalledPrintf {
		t.Error("should not call Printf for task message")
	}

	if service.out.(*shell.FakeOutputWriter).CalledPrintln {
		t.Error("should not call Println for task status")
	}
}

func TestOutputNewKoolTask(t *testing.T) {
	service := newKoolTaskServiceTestWithOutput()
	task := NewKoolTask("testing", service)

	buf := bytes.NewBufferString("")
	service.SetWriter(buf)

	_ = task.Run([]string{})

	bufBytes, err := ioutil.ReadAll(buf)

	if err != nil {
		t.Fatal(err)
	}

	output := strings.TrimSpace(string(bufBytes))

	if !strings.Contains(output, "testing output") {
		t.Error("did not called Println with KoolService output")
	}
}
