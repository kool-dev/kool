package api

import "fmt"

// DeployCreate consumes the API endpoint to create a new deployment
type DeployStatus struct {
	Endpoint
}

// DeployCreateResponse holds data returned from the deploy endpoint
type DeployStatusResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// NewDeployCreate creates a new DeployCreate instance
func NewDeployStatus(created *DeployCreateResponse) (c *DeployStatus) {
	c = &DeployStatus{
		Endpoint: NewDefaultEndpoint("GET"),
	}

	c.SetPath(fmt.Sprintf("deploy/%d/status", created.Deploy.ID))

	return
}

// Run calls deploy/?/status in the Kool Dev API
func (c *DeployStatus) Run() (resp *DeployStatusResponse, err error) {
	resp = &DeployStatusResponse{}

	c.SetResponseReceiver(resp)

	err = c.DoCall()

	return
}
