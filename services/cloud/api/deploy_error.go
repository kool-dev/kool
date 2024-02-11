package api

import "fmt"

// DeployError consumes the API endpoint to inform a build error for a deployment
type DeployError struct {
	Endpoint
}

// NewDeployError creates a new DeployError instance
func NewDeployError(created *DeployCreateResponse, err error) (c *DeployError) {
	c = &DeployError{
		Endpoint: NewDefaultEndpoint("POST"),
	}

	c.SetPath("deploy/error")
	c.Body().Set("id", fmt.Sprintf("%d", created.Deploy.ID))
	c.Body().Set("err", err.Error())

	return
}

// Run calls deploy/error in the Kool Dev API
func (c *DeployError) Run() (err error) {
	var resp interface{}
	c.SetResponseReceiver(&resp)
	err = c.DoCall()
	return
}
