package api

// DestroyCall interface represents logic for consuming the DELETE /deploy API endpoint
type DestroyCall interface {
	Endpoint

	Call() (*DestroyResponse, error)
}

// DefaultDestroyCall holds data and logic for consuming the "destroy" endpoint
type DefaultDestroyCall struct {
	Endpoint
}

// DestroyResponse holds data from the "destroy" endpoint
type DestroyResponse struct {
	Environment struct {
		ID int `json:"id"`
	} `json:"environment"`
}

// NewDefaultDestroyCall creates a new caller for Deploy API exec endpoint
func NewDefaultDestroyCall() *DefaultDestroyCall {
	return &DefaultDestroyCall{
		Endpoint: newDefaultEndpoint("DELETE"),
	}
}

// Call performs the request to the endpoint
func (s *DefaultDestroyCall) Call() (r *DestroyResponse, err error) {
	r = &DestroyResponse{}

	s.Endpoint.SetPath("deploy")
	s.Endpoint.SetResponseReceiver(r)

	err = s.Endpoint.DoCall()

	return
}
