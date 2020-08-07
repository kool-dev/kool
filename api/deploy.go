package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// Deploy represents a deployment process, from
// request to finish and retrieving the public URL.
type Deploy struct {
	tarballPath, id, Status, url string
}

// NewDeploy creates a new handler for using the
// Kool Dev API for deploying your application.
func NewDeploy(tarballPath string) (d *Deploy) {
	d = new(Deploy)
	d.tarballPath = tarballPath
	return
}

// GetID returns the ID for the deployment
func (d *Deploy) GetID() string {
	return d.id
}

// SendFile calls deploy/create in the Kool Dev API
func (d *Deploy) SendFile() (err error) {
	var (
		buff   bytes.Buffer
		file   *os.File
		fw     io.Writer
		domain string
		resp   *http.Response
		raw    []byte
	)

	w := multipart.NewWriter(&buff)

	if file, err = os.Open(d.tarballPath); err != nil {
		return
	}

	// fw, err = w.CreateFormFile("deploy", d.tarballPath)
	if fw, err = w.CreateFormFile("deploy", "deploy.tgz"); err != nil {
		return
	}

	if _, err = io.Copy(fw, file); err != nil {
		return
	}

	defer file.Close()
	if domain = os.Getenv("KOOL_DEPLOY_DOMAIN"); domain != "" {
		w.WriteField("domain", domain)
	}
	w.Close()

	req, _ := http.NewRequest("POST", apiBaseURL+"/deploy/create", &buff)
	req.Header.Add("Content-Type", w.FormDataContentType())
	if resp, err = doRequest(req); err != nil {
		return
	}

	defer resp.Body.Close()

	if raw, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		err = ErrPayloadValidation
		fmt.Println(string(raw))
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = ErrBadResponseStatus
		fmt.Println(string(raw))
	}

	if err != nil {
		return
	}

	deploy := make(map[string]interface{})

	if err = json.Unmarshal(raw, &deploy); err != nil {
		return
	}
	raw = nil

	var (
		ok  bool
		idF float64
	)

	if idF, ok = deploy["id"].(float64); ok {
		d.id = fmt.Sprintf("%d", int64(idF))
	} else {
		err = errors.New("unexpected API response. Please ask for support")
	}

	return
}

// GetStatus checks the API for the status of
// the deployment process happening in the
// background.
func (d *Deploy) GetStatus() (err error) {
	var (
		request *http.Request
		resp    *http.Response
		raw     []byte
		ok      bool
	)
	request, _ = http.NewRequest("GET", fmt.Sprintf("%s/deploy/%s/status", apiBaseURL, d.id), nil)

	if resp, err = doRequest(request); err != nil {
		return
	}

	defer resp.Body.Close()

	if raw, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	var data = make(map[string]interface{})

	if err = json.Unmarshal(raw, &data); err != nil {
		return
	}
	raw = nil

	if d.Status, ok = data["status"].(string); !ok {
		err = ErrUnexpectedResponse
		return
	}

	if d.Status == "failed" {
		err = ErrDeployFailed
		return
	}

	if d.Status == "success" {
		if d.url, ok = data["url"].(string); !ok {
			err = ErrUnexpectedResponse
			return
		}
	}

	return
}

// IsSuccessful tells whether the deployment
// process finished successfully.
func (d *Deploy) IsSuccessful() bool {
	return d.Status == "success"
}

// GetURL returns the generated URL for the deployment
// after it finishes successfully.
func (d *Deploy) GetURL() string {
	return d.url
}
