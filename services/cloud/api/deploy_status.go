package api

import "fmt"

// DeployCreate consumes the API endpoint to create a new deployment
type DeployStatus struct {
	Endpoint
}

// DeployStatusService holds a service listed in the deploy status
type DeployStatusService struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
}

// DeployStatusResponse holds data returned from the deploy status endpoint.
// Error is only set once a deploy has failed or timed out; Hint and Detail
// are optional extras explaining the failure, printed when present.
type DeployStatusResponse struct {
	ID          int                   `json:"id"`
	Status      string                `json:"status"`
	Type        string                `json:"type"`
	Url         string                `json:"url"`
	Error       string                `json:"error"`
	Hint        string                `json:"hint"`
	Detail      string                `json:"detail"`
	StartedAt   string                `json:"started_at"`
	CompletedAt string                `json:"completed_at"`
	Services    []DeployStatusService `json:"services"`
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
