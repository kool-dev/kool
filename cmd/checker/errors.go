package checker

import "errors"

// ErrDockerNotFound happens when docker is not installed
var ErrDockerNotFound = errors.New("docker doesn't seem to be installed, install it first and retry")

// IsDockerNotFoundError tells whether the given error is checker.ErrDockerNotFound
func IsDockerNotFoundError(err error) bool {
	return err.Error() == ErrDockerNotFound.Error()
}

// ErrDockerComposeNotFound happens when docker-compose is not installed
var ErrDockerComposeNotFound = errors.New("docker-compose doesn't seem to be installed, install it first and retry")

// IsDockerComposeNotFoundError tells whether the given error is checker.ErrDockerComposeNotFound
func IsDockerComposeNotFoundError(err error) bool {
	return err.Error() == ErrDockerComposeNotFound.Error()
}

// ErrDockerNotRunning happens when docker daemon is not running
var ErrDockerNotRunning = errors.New("docker daemon doesn't seem to be running, run it first and retry")

// IsDockerNotRunningError tells whether the given error is checker.ErrDockerNotRunning
func IsDockerNotRunningError(err error) bool {
	return err.Error() == ErrDockerNotRunning.Error()
}
