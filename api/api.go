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

// BadAPIServer represents some issue in the API side
type BadAPIServer string

func (err BadAPIServer) Error() string {
	return string(err)
}

// Deploy represents a deployment process, from
// request to finish and retrieving the public URL.
type Deploy struct {
	tarballPath string
}

var (
	apiBaseURL string = "https://kool.dev/api"
)

// SetBaseURL defines the target Kool API URL to be used
// when reaching out endpoints.
func SetBaseURL(url string) {
	apiBaseURL = url
}

// NewDeploy creates a new handler for using the
// Kool Dev API for deploying your application.
func NewDeploy(tarballPath string) (d *Deploy) {
	d = new(Deploy)
	d.tarballPath = tarballPath
	return
}

// SendFile calls deploy/create in the Kool Dev API
func (d *Deploy) SendFile() (id string, err error) {
	var (
		buff           bytes.Buffer
		file           *os.File
		fw             io.Writer
		slug, apiToken string
		resp           *http.Response
		raw            []byte
	)

	apiToken = os.Getenv("KOOL_API_TOKEN")

	if apiToken == "" {
		err = errors.New("missing KOOL_API_TOKEN")
		return
	}

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
	if slug = os.Getenv("KOOL_DEPLOY_SLUG"); slug != "" {
		w.WriteField("slug", slug)
	}
	w.Close()

	req, _ := http.NewRequest("POST", apiBaseURL+"/deploy/create", &buff)
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiToken)
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()

	if raw, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		err = BadAPIServer("unauthorized; please check your KOOL_API_TOKEN")
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		err = BadAPIServer("something went wrong validating the payload")
		fmt.Println(string(raw))
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = BadAPIServer("unexpected return status")
	}

	if err != nil {
		return
	}

	deploy := make(map[string]interface{})

	if err = json.Unmarshal(raw, &deploy); err != nil {
		return
	}

	var (
		ok  bool
		idF float64
	)

	if idF, ok = deploy["id"].(float64); ok {
		id = fmt.Sprintf("%d", int64(idF))
	} else {
		err = errors.New("unexpected API response. Please ask for support")
	}

	return
}
