package api

import "fmt"

// DeployPingResponse holds data returned from the deploy endpoint
type DeployPingResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// DeployPing consumes the API endpoint to create a new deployment
type DeployPing struct {
	Endpoint
}

// NewDeployPing creates a new DeployStart instance
func NewDeployPing(created *DeployCreateResponse) (c *DeployPing) {
	c = &DeployPing{
		Endpoint: NewDefaultEndpoint("POST"),
	}

	c.SetPath("deploy/ping")
	c.Body().Set("id", fmt.Sprintf("%d", created.Deploy.ID))

	return
}

// Run calls deploy/ping in the Kool Dev API
func (c *DeployPing) Run() (resp *DeployPingResponse, err error) {
	resp = &DeployPingResponse{}
	c.SetResponseReceiver(resp)
	err = c.DoCall()
	return
}
