package api

// DeployDestroy holds data and logic for consuming the "destroy" endpoint
type DeployDestroy struct {
	Endpoint
}

// DeployDestroyResponse holds data from the "destroy" endpoint
type DeployDestroyResponse struct {
	Environment struct {
		ID int `json:"id"`
	} `json:"environment"`
}

// NewDeployDestroy creates a new caller for Deploy API exec endpoint
func NewDeployDestroy() (d *DeployDestroy) {
	d = &DeployDestroy{
		Endpoint: NewDefaultEndpoint("DELETE"),
	}

	d.SetPath("deploy/destroy")

	return
}

// Call performs the request to the endpoint
func (s *DeployDestroy) Call() (resp *DeployDestroyResponse, err error) {
	resp = &DeployDestroyResponse{}
	s.SetResponseReceiver(resp)
	err = s.DoCall()
	return
}
