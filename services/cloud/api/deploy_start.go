package api

import "fmt"

// DeployStartResponse holds data returned from the deploy endpoint
type DeployStartResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// DeployStart consumes the API endpoint to create a new deployment
type DeployStart struct {
	Endpoint
}

// NewDeployStart creates a new DeployStart instance
func NewDeployStart(created *DeployCreateResponse) (c *DeployStart) {
	c = &DeployStart{
		Endpoint: NewDefaultEndpoint("POST"),
	}

	c.SetPath("deploy/start")
	c.Body().Set("id", fmt.Sprintf("%d", created.Deploy.ID))

	return
}

// Run calls deploy/create in the Kool Dev API
func (c *DeployStart) Run() (resp *DeployStartResponse, err error) {
	resp = &DeployStartResponse{}
	c.SetResponseReceiver(resp)
	err = c.DoCall()
	return
}
