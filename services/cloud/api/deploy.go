package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"kool-dev/kool/core/environment"
	"mime/multipart"
	"net/http"
	"os"
)

// Deploy represents a deployment process, from
// request to finish and retrieving the public URL.
type Deploy struct {
	Endpoint

	tarballPath, id string

	env environment.EnvStorage

	Status *StatusResponse
}

// DeployResponse holds data returned from the deploy endpoint
type DeployResponse struct {
	ID int `json:"id"`
}

// NewDeploy creates a new handler for using the
// Kool Dev API for deploying your application.
func NewDeploy(tarballPath string) *Deploy {
	return &Deploy{
		Endpoint:    NewDefaultEndpoint("POST"),
		env:         environment.NewEnvStorage(),
		tarballPath: tarballPath,
	}
}

// GetID returns the ID for the deployment
func (d *Deploy) GetID() string {
	return d.id
}

// SendFile calls deploy/create in the Kool Dev API
func (d *Deploy) SendFile() (err error) {
	var (
		body io.Reader
		resp = &DeployResponse{}
	)

	if body, err = d.getPayload(); err != nil {
		return
	}

	d.SetPath("deploy/create")
	d.SetRawBody(body)
	d.SetResponseReceiver(resp)
	if err = d.DoCall(); err != nil {
		if errAPI, is := err.(*ErrAPI); is {
			// override the error for a better message
			if errAPI.Status == http.StatusUnauthorized {
				err = ErrUnauthorized
			} else if errAPI.Status == http.StatusUnprocessableEntity {
				fmt.Println(errAPI.Message)
				fmt.Println(errAPI.Errors)
				err = ErrPayloadValidation
			} else if errAPI.Status != http.StatusOK && errAPI.Status != http.StatusCreated {
				err = ErrBadResponseStatus
			}
		}
		return
	}

	d.id = fmt.Sprintf("%d", resp.ID)
	if d.id == "0" {
		err = errors.New("unexpected API response, please reach out for support on Slack or Github")
	}

	return
}

func (d *Deploy) getPayload() (body io.Reader, err error) {
	var (
		buff         bytes.Buffer
		file         *os.File
		fw           io.Writer
		cluster      string
		domain       string
		domainExtras string
		wwwRedirect  string
	)

	w := multipart.NewWriter(&buff)

	if file, err = os.Open(d.tarballPath); err != nil {
		return
	}

	fi, _ := file.Stat()
	fmt.Printf("Release tarball got %.2fMBs...\n", float64(fi.Size())/1024/1024)

	if fw, err = w.CreateFormFile("deploy", "deploy.tgz"); err != nil {
		return
	}

	if _, err = io.Copy(fw, file); err != nil {
		return
	}

	defer file.Close()

	if cluster = d.env.Get("KOOL_DEPLOY_CLUSTER"); cluster != "" {
		if err = w.WriteField("cluster", cluster); err != nil {
			return
		}
	}

	if domain = d.env.Get("KOOL_DEPLOY_DOMAIN"); domain != "" {
		if err = w.WriteField("domain", domain); err != nil {
			return
		}
	}

	if domainExtras = d.env.Get("KOOL_DEPLOY_DOMAIN_EXTRAS"); domainExtras != "" {
		if err = w.WriteField("domain_extras", domainExtras); err != nil {
			return
		}
	}

	if wwwRedirect = d.env.Get("KOOL_DEPLOY_WWW_REDIRECT"); wwwRedirect != "" {
		if err = w.WriteField("www_redirect", wwwRedirect); err != nil {
			return
		}
	}

	d.SetContentType(w.FormDataContentType())
	w.Close()

	body = &buff
	return
}

// FetchLatestStatus checks the API for the status of the deployment process
// happening in the background
func (d *Deploy) FetchLatestStatus() (err error) {
	if d.Status, err = NewDefaultStatusCall(d.id).Call(); err != nil {
		return
	}

	if d.Status.Status == "failed" {
		err = ErrDeployFailed
		return
	}

	return
}

// IsSuccessful tells whether the deployment process finished successfully
func (d *Deploy) IsSuccessful() bool {
	return d.Status.Status == "success"
}

// GetURL returns the generated URL for the deployment after it finishes successfully
func (d *Deploy) GetURL() string {
	return d.Status.URL
}
