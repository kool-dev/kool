package cmd

import (
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"regexp"

	"github.com/spf13/cobra"
)

// KoolShareFlags holds the flags for the kool stop command
type KoolShareFlags struct {
	Service   string
	Subdomain string
}

// KoolShare holds handlers and functions to implement the stop command logic
type KoolShare struct {
	DefaultKoolService
	Flags *KoolShareFlags

	env environment.EnvStorage

	status *KoolStatus
	share  builder.Command
}

func init() {
	var (
		share    = NewKoolShare()
		shareCmd = NewShareCommand(share)
	)

	rootCmd.AddCommand(shareCmd)
}

// NewKoolShare creates a new handler for sharing local environment with default dependencies
func NewKoolShare() *KoolShare {
	defaultKoolService := newDefaultKoolService()
	return &KoolShare{
		*defaultKoolService,
		&KoolShareFlags{"app", ""},
		environment.NewEnvStorage(),
		NewKoolStatus(),
		builder.NewCommand("docker", "run", "--rm", "--init"),
	}
}

// validSubdomain runs the stop logic with incoming arguments.
func (s *KoolShare) validSubdomain(subdomain string) bool {
	return regexp.MustCompile("^[a-z]+[a-z0-9]+$").MatchString(subdomain)
}

// Execute runs the stop logic with incoming arguments.
func (s *KoolShare) Execute(args []string) (err error) {
	var isRunning bool

	if isRunning, _, _, err = s.status.getServiceInfo(s.Flags.Service); err != nil {
		return
	}

	if !isRunning {
		err = fmt.Errorf("service %s is not running, please check kool status and use --service flag to set which service to share", s.Flags.Service)
		return
	}

	s.share.AppendArgs("--network", s.env.Get("KOOL_GLOBAL_NETWORK"))
	s.share.AppendArgs("beyondcodegmbh/expose-server:1.4.1", "share")
	s.share.AppendArgs(s.Flags.Service)
	s.share.AppendArgs("--server-host", "kool.live")

	if s.Flags.Subdomain != "" {
		if !s.validSubdomain(s.Flags.Subdomain) {
			err = fmt.Errorf("invalid subdomain '%s'", s.Flags.Subdomain)
			return
		}

		s.share.AppendArgs("--subdomain", s.Flags.Subdomain)
	}

	err = s.Interactive(s.share)
	return
}

// NewShareCommand initializes new kool stop command
func NewShareCommand(share *KoolShare) (shareCmd *cobra.Command) {
	shareCmd = &cobra.Command{
		Use:   "share",
		Short: "Live share your local environment through an HTTP tunnel with anyone, anywhere.",
		Args:  cobra.NoArgs,
		Run:   DefaultCommandRunFunction(share),
	}

	shareCmd.Flags().StringVarP(&share.Flags.Service, "service", "", "app", "The name of the local service container we want to share.")
	shareCmd.Flags().StringVarP(&share.Flags.Subdomain, "subdomain", "", "", "The subdomain desired for subdomain.kool.dev.")
	return
}
