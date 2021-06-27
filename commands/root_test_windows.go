package commands

import (
	"bytes"
	"fmt"
	"io"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"os"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
)

func assertExecGotError(t *testing.T, cmd *cobra.Command, partialErr string) {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	if err := cmd.Execute(); err == nil {
		t.Errorf("should have got an error - %s", partialErr)
	} else if !strings.Contains(err.Error(), partialErr) {
		t.Errorf("unexpected error executing command; '%s' but got error: %v", partialErr, err)
	}
}
