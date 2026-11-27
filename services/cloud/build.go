package cloud

import (
	"bytes"
	"errors"
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"kool-dev/kool/core/shell"
	"kool-dev/kool/services/cloud/api"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func BuildPushImageForDeploy(service string, config *DeployConfigService, deploy *api.DeployCreateResponse, platform string) (err error) {
	if config.Build == nil {
		err = errors.New("service " + service + " has no build configuration")
		return
	}

	var (
		env    = environment.NewEnvStorage()
		sh     = shell.NewShell()
		image  = fmt.Sprintf("%s/%s:%s-%s", deploy.Config.ImagePrefix, deploy.Config.ImageRepository, service, deploy.Config.ImageTag)
		output string
	)

	// create a default io.Reader for stdin
	var in = bytes.NewBuffer([]byte{})
	sh.SetInStream(in)

	dockerBuild := builder.NewCommand("docker", "build", "-t", image, "--platform", platform)

	if folder, isStr := (*config.Build).(string); isStr {
		// this should be a simple build with a context folder
		// docker build -t <image>:<tag> <folder>
		dockerBuild.AppendArgs(parseContext(folder, env.Get("PWD")))
	} else {
		// it's not a string, so it should be a map...
		var buildConfig *DeployConfigBuild
		if buildConfig, err = parseBuild(*config.Build); err != nil {
			return
		}

		if buildConfig.Dockerfile != nil {
			dockerBuild.AppendArgs("-f", *buildConfig.Dockerfile)
		}

		if buildConfig.Args != nil {
			for k, v := range *buildConfig.Args {
				dockerBuild.AppendArgs("--build-arg", fmt.Sprintf("%s=%s", k, parseDeployEnvs(v, deploy.Deploy.Environment.Env)))
			}
		}

		if buildConfig.Context != nil {
			dockerBuild.AppendArgs(parseContext(*buildConfig.Context, env.Get("PWD")))
		}
	}

	if err = sh.Interactive(dockerBuild); err != nil {
		return
	}

	if _, err = in.Write([]byte(deploy.Docker.Password)); err != nil {
		return
	}

	// login & push...
	dockerLogin := builder.NewCommand("docker", "login", "-u", deploy.Docker.Login, "--password-stdin", deploy.Config.ImagePrefix)
	if output, err = sh.Exec(dockerLogin); err != nil {
		err = dockerLoginError(deploy.Config.ImagePrefix, output, err)
		return
	}

	if deploy.Config.ImagePrefix != "" {
		defer dockerLogout(sh, dockerLogoutCommand(deploy.Config.ImagePrefix))
	}

	dockerPush := builder.NewCommand("docker", "push", image)
	err = sh.Interactive(dockerPush)

	return
}

// dockerLogoutCommand logs out of the given registry only - a bare
// docker logout would log the user out of Docker Hub instead
func dockerLogoutCommand(registry string) builder.Command {
	return builder.NewCommand("docker", "logout", registry)
}

// dockerLogout runs the logout after the push; by then the image is pushed
// (or the push failed already), so a failed logout is only a warning
func dockerLogout(sh shell.Shell, logout builder.Command) {
	if _, err := sh.Exec(logout); err != nil {
		sh.Warning("failed to log out of the Kool.dev Cloud registry: ", err)
	}
}

// dockerLoginError tells which registry docker login failed for, along
// with what docker printed - not only its exit status.
func dockerLoginError(registry, output string, err error) error {
	if output = strings.TrimSpace(output); output != "" && !strings.Contains(err.Error(), output) {
		return fmt.Errorf("docker login to %s failed: %s (%w)", registry, output, err)
	}

	return fmt.Errorf("docker login to %s failed: %w", registry, err)
}

// parseContext parses the build context from the build configuration
// changing . to the current working directory
func parseContext(context string, cwd string) (parsed string) {
	context = strings.Trim(context, " ")
	parsed = context

	if strings.HasPrefix(context, "..") {
		// relative path
		parsed = filepath.Join(cwd, context)
	} else if strings.HasPrefix(context, ".") {
		// relative path
		parsed = filepath.Join(cwd, strings.TrimPrefix(context, "."))
	}

	return
}

func parseBuild(build interface{}) (config *DeployConfigBuild, err error) {
	var (
		b []byte
	)

	config = &DeployConfigBuild{}

	if b, err = yaml.Marshal(build); err != nil {
		return
	}

	err = yaml.Unmarshal(b, config)
	return
}

func parseDeployEnvs(i interface{}, env interface{}) string {
	// workaround to allow for escaping $ with a double $$
	_ = os.Setenv("$", "$")

	var value = os.ExpandEnv(fmt.Sprintf("%s", i))

	if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
		// we have something to replace!
		if recs, ok := env.(map[string]interface{}); ok {
			for k, v := range recs {
				if strings.Contains(value, "{{"+k+"}}") {
					value = strings.ReplaceAll(value, "{{"+k+"}}", fmt.Sprintf("%v", v))
				}
			}
		}
	}

	return value
}
