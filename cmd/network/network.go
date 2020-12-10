package network

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
)

// Handler defines network handler
type Handler interface {
	HandleGlobalNetwork(string) error
}

// DefaultHandler holds docker network command
type DefaultHandler struct {
	CheckNetworkCmd  builder.Command
	CreateNetworkCmd builder.Command
	shell            shell.Shell
}

// NewHandler initializes handler
func NewHandler(s shell.Shell) *DefaultHandler {
	var checkNetCmd, createNetCmd *builder.DefaultCommand

	checkNetCmd = builder.NewCommand("docker", "network", "ls", "-q", "-f")
	createNetCmd = builder.NewCommand("docker", "network", "create", "--attachable")

	return &DefaultHandler{checkNetCmd, createNetCmd, s}
}

// HandleGlobalNetwork handles global network
func (h *DefaultHandler) HandleGlobalNetwork(networkName string) error {
	if networkID, err := h.shell.Exec(h.CheckNetworkCmd, fmt.Sprintf("NAME=^%s$", networkName)); err != nil || networkID != "" {
		return err
	}

	return h.shell.Interactive(h.CreateNetworkCmd, networkName)
}
