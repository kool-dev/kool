package api

// ExecCall interface represents logic for consuming the deploy/exec API endpoint
type ExecCall interface {
	Endpoint

	Call() (*ExecResponse, error)
}

// DefaultExecCall holds data and logic for consuming the "exec" endpoint
type DefaultExecCall struct {
	Endpoint

	domain string
}

// ExecResponse holds data from the "exec" endpoint
type ExecResponse struct {
	Server    string `json:"server"`
	Namespace string `json:"namespace"`
	Path      string `json:"path"`
	Token     string `json:"token"`
	CA        string `json:"ca.crt"`
}

// NewDefaultExecCall creates a new caller for Deploy API exec endpoint
func NewDefaultExecCall() *DefaultExecCall {
	return &DefaultExecCall{
		Endpoint: newDefaultEndpoint("POST"),
	}
}

// Call performs the request to the endpoint
func (s *DefaultExecCall) Call() (r *ExecResponse, err error) {
	r = &ExecResponse{}

	s.Endpoint.SetPath("deploy/exec")
	s.Endpoint.SetResponseReceiver(r)

	err = s.Endpoint.DoCall()

	return
}
