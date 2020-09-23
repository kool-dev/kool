package network

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
)

// Handler defines network handler
type Handler interface {
	HandleGlobalNetwork(string) error
}

// DefaultHandler holds docker network command
type DefaultHandler struct {
	CheckNetworkCmd  builder.Runner
	CreateNetworkCmd builder.Runner
}

// NewHandler initializes handler
func NewHandler() *DefaultHandler {
	var checkNetCmd, createNetCmd *builder.DefaultCommand

	checkNetCmd = builder.NewCommand("docker", "network", "ls", "-q", "-f")
	createNetCmd = builder.NewCommand("docker", "network", "create", "--attachable")

	return &DefaultHandler{checkNetCmd, createNetCmd}
}

// HandleGlobalNetwork handles global network
func (h *DefaultHandler) HandleGlobalNetwork(networkName string) error {
	if networkID, err := h.CheckNetworkCmd.Exec(fmt.Sprintf("NAME=^%s$", networkName)); err != nil || networkID != "" {
		return err
	}

	return h.CreateNetworkCmd.Interactive(networkName)
}
