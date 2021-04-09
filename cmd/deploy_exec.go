package cmd

import (
	"fmt"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var authTempPath = "/tmp"

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService

	kubectl, kool builder.Command

	env     environment.EnvStorage
	apiExec api.ExecCall
}

// NewDeployExecCommand initializes new kool deploy Cobra command
func NewDeployExecCommand(deployExec *KoolDeployExec) *cobra.Command {
	return &cobra.Command{
		Use:   "exec SERVICE COMMAND [--] [ARG...]",
		Short: "Execute a command inside a running service container deployed to Kool Cloud",
		Long: `After deploying your application to Kool Cloud using 'kool deploy',
execute a COMMAND inside the specified SERVICE container (similar to an SSH session).
Must use a KOOL_API_TOKEN environment variable for authentication.`,
		Args: cobra.MinimumNArgs(1),
		Run:  DefaultCommandRunFunction(deployExec),

		DisableFlagsInUseLine: true,
	}
}

// NewKoolDeployExec creates a new pointer with default KoolDeployExec service dependencies
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
		builder.NewCommand("kubectl"),
		builder.NewCommand("kool"),
		environment.NewEnvStorage(),
		api.NewDefaultExecCall(),
	}
}

// Execute runs the deploy exec logic - integrating with Deploy API
func (e *KoolDeployExec) Execute(args []string) (err error) {
	var (
		domain  string
		service string
		resp    *api.ExecResponse
	)

	if len(args) == 0 {
		err = fmt.Errorf("KoolDeployExec.Execute: required at least one argument")
		return
	}

	service = args[0]
	args = args[1:]

	if url := e.env.Get("KOOL_API_URL"); url != "" {
		api.SetBaseURL(url)
	}

	if domain = e.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	e.apiExec.Body().Set("domain", domain)
	e.apiExec.Body().Set("service", service)

	if resp, err = e.apiExec.Call(); err != nil {
		return
	}

	if resp.Token == "" {
		err = fmt.Errorf("failed to generate access credentials to cloud deploy")
		return
	}

	CAPath := filepath.Join(authTempPath, ".kool-cluster-CA")
	defer func() {
		if err := os.Remove(CAPath); err != nil {
			e.Warning("failed to clear up temporary file; error:", err.Error())
		}
	}()
	if err = os.WriteFile(CAPath, []byte(resp.CA), os.ModePerm); err != nil {
		return
	}

	e.kubectl.AppendArgs("--server", resp.Server)
	e.kubectl.AppendArgs("--token", resp.Token)
	e.kubectl.AppendArgs("--namespace", resp.Namespace)
	e.kubectl.AppendArgs("--certificate-authority", CAPath)
	e.kubectl.AppendArgs("exec", "-i")
	if e.IsTerminal() {
		e.kubectl.AppendArgs("-t")
	}
	e.kubectl.AppendArgs(resp.Path, "--")
	if len(args) == 0 {
		args = []string{"bash"}
	}
	e.kubectl.AppendArgs(args...)

	if e.LookPath(e.kubectl) == nil {
		// the command is available on current PATH, so let's use it
		err = e.Interactive(e.kubectl)
		return
	}

	// we do not have 'kubectl' on current path... let's use a container!
	e.kool.AppendArgs(
		"docker", "--",
		"-v", fmt.Sprintf("%s:%s", CAPath, CAPath),
		"kooldev/toolkit:full",
		e.kubectl.Cmd(),
	)
	e.kool.AppendArgs(e.kubectl.Args()...)

	err = e.Interactive(e.kool)
	return
}
