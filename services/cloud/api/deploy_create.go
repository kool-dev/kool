package api

// DeployCreateResponse holds data returned from the deploy endpoint
type DeployCreateResponse struct {
	Deploy struct {
		ID          int    `json:"id"`
		Project     string `json:"project"`
		Url         string `json:"url"`
		Environment struct {
			Name string      `json:"name"`
			Env  interface{} `json:"env"`
		} `json:"environment"`
		Cluster struct {
			Region string `json:"region"`
		} `json:"cluster"`
	} `json:"deploy"`

	Config struct {
		ImagePrefix     string `json:"image_prefix"`
		ImageRepository string `json:"image_repository"`
		ImageTag        string `json:"image_tag"`
	} `json:"stuff"`

	Docker struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"docker"`

	LogsUrl string `json:"logs_url"`
}

// DeployCreate consumes the API endpoint to create a new deployment
type DeployCreate struct {
	Endpoint
}

// NewDeployCreate creates a new DeployCreate instance
func NewDeployCreate() (c *DeployCreate) {
	c = &DeployCreate{
		Endpoint: NewDefaultEndpoint("POST"),
	}

	c.SetPath("deploy/create")
	_ = c.PostField("is_local", "1")

	return
}

// Run calls deploy/create in the Kool Dev API
func (c *DeployCreate) Run() (resp *DeployCreateResponse, err error) {
	resp = &DeployCreateResponse{}

	c.SetResponseReceiver(resp)

	err = c.DoCall()

	return
}
