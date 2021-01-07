package cmd

import (
	// "k8s.io/cli-runtime/pkg/genericclioptions"
	// kubectlExec "k8s.io/kubectl/pkg/cmd/exec"

	kubectl "k8s.io/kubectl/pkg/cmd"
	// "k8s.io/kubectl/pkg/polymorphichelpers"

	"github.com/spf13/cobra"
)

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService
}

// NewDeployExecCommand initializes new kool deploy Cobra command
func NewDeployExecCommand(deployExec *KoolDeployExec) *cobra.Command {
	return &cobra.Command{
		Use:   "exec",
		Short: "Executes a command in your deployed application on Kool cloud",
		Run:   DefaultCommandRunFunction(deployExec),
	}
}

// NewKoolDeployExec creates a new pointer with default KoolDeploy service
// dependencies.
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
	}
}

// Execute runs the deploy logic.
func (e *KoolDeployExec) Execute(args []string) (err error) {
	e.Println("kool deploy exec - start")

	cmd := kubectl.NewKubectlCommand(e.InStream(), e.OutStream(), e.ErrStream())

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6Im5RYXNwLUFRR0Y2UVlQa2ZnX3ZQc1pESHAyRkNINTdMV2FxVnJSM2NMazAifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrb29sIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Imtvb2wtdG9rZW4tbW5wZzQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoia29vbCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjVkNDBlYWU3LWNhMjEtNDZhMC1iZWZlLTEzNWY3YTk1NzdlZSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprb29sOmtvb2wifQ.i7_OuD0ByRr7vmK6XjQ42YE0M1hoom5Mr8PhLaDhuyipk3n64u7dy-OI-EoJHbW8C0jUu84e55HHVPq37KowwokXYAMV0r2XM83AiX072E39LES8q2o_XgCVgMDiVW6k4qAMR8Xt-IP9jB7Bfb6_Sr_VzqHxlyvIUoPjtzDYX4rCnTZrMmmRKsn1JWi0GzQ5_GxfbayRlEik45bfL-ztk1tMTafG5UGe3piqJ63iQ02tlH1_A0Hv8lcv7uqGD9MbPD5xmslkEmjWct1Haw3ik3y9gexexX7PJXW_INS78g20G1swNLciQg9cSiy05E0xeZguj2cg0fHJJ6EGYKZZbQ"

	cmd.SetArgs([]string{"--token", token, "-n", "kool", "exec", "-it", "deployment/kool", "--", "bash"})
	err = cmd.Execute()

	// streams := genericclioptions.IOStreams{
	// 	In:     e.InStream(),
	// 	Out:    e.OutStream(),
	// 	ErrOut: e.ErrStream(),
	// }

	// f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	// cmd := kubectlExec.NewCmdExec(f, streams)

	// ex := &kubectlExec.ExecOptions{
	// 	StreamOptions: kubectlExec.StreamOptions{
	// 		IOStreams: streams,
	// 	},

	// 	Executor: &kubectlExec.DefaultRemoteExecutor{},
	// }

	// ex.ResourceName = "deployment/kool"
	// ex.Command = []string{"php", "-v"}
	// ex.EnforceNamespace = true
	// ex.Namespace = "kool"
	// ex.ExecutablePodFn = polymorphichelpers.AttachablePodForObjectFn
	// ex.GetPodTimeout = time.Minute * 100

	return
}
