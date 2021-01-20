package api

import (
	"kool-dev/kool/environment"
	"testing"
)

func TestNewDeploy(t *testing.T) {
	var tarball = "tarball"
	d := NewDeploy(tarball)

	if _, ok := d.env.(*environment.DefaultEnvStorage); !ok {
		t.Error("unexpected default environment.EnvStorage")
	}

	if _, ok := d.Endpoint.(*DefaultEndpoint); !ok {
		t.Error("unexpected default Endpoint")
	}

	if tarball != d.tarballPath {
		t.Error("failed setting tarballPath")
	}

	var id = "id"
	d.id = id

	if id != d.GetID() {
		t.Error("failed setting id")
	}

	var url = "url"
	d.Status = &StatusResponse{Status: "success", URL: url}

	if !d.IsSuccessful() {
		t.Error("failed asserting success")
	}

	if url != d.GetURL() {
		t.Error("failed getting URL")
	}

	// still in need of testing
	// d.FetchLatestStatus()
	// d.SendFile()
	// d.getPayload()
}
