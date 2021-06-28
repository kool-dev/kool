package commands

import (
	"errors"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"testing"
)

func TestShareDefaults(t *testing.T) {
	share := NewKoolShare()

	if share.Flags.Service != "app" {
		t.Errorf("bad default service; expected app but got %s", share.Flags.Service)
	}

	if _, ok := share.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("bad default environment.EnvStorage implementation")
	}

	if len(share.share.Args()) != 3 || share.share.Cmd() != "docker" {
		t.Error("bad default builder.Command for sharing")
	}
}

func newFakeShareService() *KoolShare {
	return &KoolShare{
		*newFakeKoolService(),
		&KoolShareFlags{"default-service", "default-subdomain", 0},
		environment.NewFakeEnvStorage(),
		newFakeKoolStatus(),
		&builder.FakeCommand{},
	}
}

func TestFlagParseServiceURI(t *testing.T) {
	f := &KoolShareFlags{"service", "", 10}

	if f.parseServiceURI() != "service:10" {
		t.Errorf("bad service URI generated from flags; expected service:10 but got: %s", f.parseServiceURI())
	}

	f.Port = 0

	if f.parseServiceURI() != "service" {
		t.Errorf("bad service URI generated from flags; expected service but got: %s", f.parseServiceURI())
	}
}

func TestShareCommand(t *testing.T) {
	share := newFakeShareService()
	share.status.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	share.status.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewShareCommand(share)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error on sharing: %v", err)
	}
}

func TestShareCommandBadDomain(t *testing.T) {
	share := newFakeShareService()
	share.status.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	share.status.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewShareCommand(share)
	cmd.SetArgs([]string{"--subdomain", "-sub"})
	assertExecGotError(t, cmd, "invalid subdomain")
}

func TestShareCommandServiceNotRunning(t *testing.T) {
	share := newFakeShareService()

	cmd := NewShareCommand(share)
	assertExecGotError(t, cmd, "is not running")
}

func TestShareCommandServiceDoesNotExist(t *testing.T) {
	share := newFakeShareService()
	share.status.getServiceIDCmd.(*builder.FakeCommand).MockExecError = errors.New("fake error")

	cmd := NewShareCommand(share)
	assertExecGotError(t, cmd, "fake error")
}

func TestShareCommandSetFlags(t *testing.T) {
	share := newFakeShareService()
	share.status.getServiceIDCmd.(*builder.FakeCommand).MockExecOut = "100"
	share.status.getServiceStatusPortCmd.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewShareCommand(share)
	cmd.SetArgs([]string{"--subdomain", "sub", "--service", "foo"})
	if err := cmd.Execute(); err != nil {
		t.Error("unexpected error on running command")
	} else if err = share.shell.(*shell.FakeShell).Err; err != nil {
		t.Error("unexpected error")
	}
	args := share.share.(*builder.FakeCommand).ArgsAppend
	if args[4] != "foo" {
		t.Error("failed setting service")
	}
	if args[8] != "sub" {
		t.Error("failed setting subdomain")
	}
}
