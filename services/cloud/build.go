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

func BuildPushImageForDeploy(service string, config *DeployConfigService, deploy *api.DeployCreateResponse) (err error) {
	if config.Build == nil {
		err = errors.New("service " + service + " has no build configuration")
		return
	}

	var (
		env       = environment.NewEnvStorage()
		isVerbose = env.IsTrue("KOOL_VERBOSE")
		sh        = shell.NewShell()
		image     = fmt.Sprintf("%s/%s:%s-%s", deploy.Config.ImagePrefix, deploy.Config.ImageRepository, service, deploy.Config.ImageTag)
		output    string
	)

	// create a default io.Reader for stdin
	var in = bytes.NewBuffer([]byte{})
	sh.SetInStream(in)

	dockerBuild := builder.NewCommand("docker", "build", "-t", image, "--platform", "linux/amd64")

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
		if isVerbose {
			fmt.Println(output)
		}
		return
	}

	dockerPush := builder.NewCommand("docker", "push", image)
	if err = sh.Interactive(dockerPush); err != nil {
		return
	}

	dockerLogout := builder.NewCommand("docker", "logout")
	if output, err = sh.Exec(dockerLogout); err != nil {
		if isVerbose {
			fmt.Println(output)
		}
		return
	}

	return
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
