package network

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"os"
)

// Handler defines network handler
type Handler interface {
	HandleGlobalNetwork() error
}

// DefaultHandler holds docker network command
type DefaultHandler struct {
	CheckNetworkCmd  builder.Runner
	CreateNetworkCmd builder.Runner
}

// NewHandler initializes handler
func NewHandler() *DefaultHandler {
	var checkNetCmd, createNetCmd *builder.DefaultCommand

	globalNetworkName := os.Getenv("KOOL_GLOBAL_NETWORK")

	checkNetCmd = builder.NewCommand("docker", "network", "ls", "-q", "-f", fmt.Sprintf("NAME=^%s$", globalNetworkName))
	createNetCmd = builder.NewCommand("docker", "network", "create", "--attachable", globalNetworkName)

	return &DefaultHandler{checkNetCmd, createNetCmd}
}

// HandleGlobalNetwork handles global network
func (h *DefaultHandler) HandleGlobalNetwork() error {
	if networkID, err := h.CheckNetworkCmd.Exec(); err != nil || networkID != "" {
		return err
	}

	return h.CreateNetworkCmd.Interactive()
}
