package api

// DeployExec holds data and logic for consuming the "exec" endpoint
type DeployExec struct {
	Endpoint
}

// DeployExecResponse holds data from the "exec" endpoint
type DeployExecResponse struct {
	Server    string `json:"server"`
	Namespace string `json:"namespace"`
	Path      string `json:"path"`
	Token     string `json:"token"`
	CA        string `json:"ca.crt"`
}

// NewDeployExec creates a new caller for Deploy API exec endpoint
func NewDeployExec() (e *DeployExec) {
	e = &DeployExec{
		Endpoint: NewDefaultEndpoint("POST"),
	}

	e.SetPath("deploy/exec")

	return e
}

// Call performs the request to the endpoint
func (s *DeployExec) Call() (resp *DeployExecResponse, err error) {
	resp = &DeployExecResponse{}
	s.SetResponseReceiver(resp)
	err = s.DoCall()
	return
}
