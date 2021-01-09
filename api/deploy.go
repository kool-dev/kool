package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kool-dev/kool/environment"
	"mime/multipart"
	"net/http"
	"os"
)

// Deploy represents a deployment process, from
// request to finish and retrieving the public URL.
type Deploy struct {
	tarballPath, id, url string

	Status *StatusResponse
}

// NewDeploy creates a new handler for using the
// Kool Dev API for deploying your application.
func NewDeploy(tarballPath string) *Deploy {
	return &Deploy{
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
		buff         bytes.Buffer
		file         *os.File
		fw           io.Writer
		domain       string
		domainExtras string
		wwwRedirect  string
		resp         *http.Response
		raw          []byte
	)

	w := multipart.NewWriter(&buff)

	if file, err = os.Open(d.tarballPath); err != nil {
		return
	}

	fi, _ := file.Stat()
	fmt.Printf("Release tarball got %.2fMBs...\n", float64(fi.Size())/1024/1024)

	// fw, err = w.CreateFormFile("deploy", d.tarballPath)
	if fw, err = w.CreateFormFile("deploy", "deploy.tgz"); err != nil {
		return
	}

	if _, err = io.Copy(fw, file); err != nil {
		return
	}

	defer file.Close()

	if domain = environment.NewEnvStorage().Get("KOOL_DEPLOY_DOMAIN"); domain != "" {
		if err = w.WriteField("domain", domain); err != nil {
			return
		}
	}

	if domainExtras = environment.NewEnvStorage().Get("KOOL_DEPLOY_DOMAIN_EXTRAS"); domainExtras != "" {
		if err = w.WriteField("domain_extras", domainExtras); err != nil {
			return
		}
	}

	if wwwRedirect = environment.NewEnvStorage().Get("KOOL_DEPLOY_WWW_REDIRECT"); wwwRedirect != "" {
		if err = w.WriteField("www_redirect", wwwRedirect); err != nil {
			return
		}
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

	var (
		ok  bool
		idF float64
	)

	if idF, ok = deploy["id"].(float64); ok {
		d.id = fmt.Sprintf("%d", int64(idF))
	} else {
		err = errors.New("unexpected API response, please ask for support")
	}

	return
}

// FetchLatestStatus checks the API for the status of the deployment process
// happening in the background
func (d *Deploy) FetchLatestStatus() (err error) {
	if d.Status, err = NewStatusCall(d.id).Call(); err != nil {
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
