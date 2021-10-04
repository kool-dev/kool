package commands

import (
	"fmt"
	"kool-dev/kool/core/builder"
	"kool-dev/kool/core/environment"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// KoolShareFlags holds the flags for the kool share command
type KoolShareFlags struct {
	Service   string
	Subdomain string
	Port      uint
}

func (f *KoolShareFlags) parseServiceURI() string {
	if f.Port != 0 {
		return fmt.Sprintf("%s:%d", f.Service, f.Port)
	}

	return f.Service
}

// KoolShare holds handlers and functions to implement the share command logic
type KoolShare struct {
	DefaultKoolService
	Flags *KoolShareFlags

	env environment.EnvStorage

	status *KoolStatus
	share  builder.Command
}

func AddKoolShare(root *cobra.Command) {
	var (
		share    = NewKoolShare()
		shareCmd = NewShareCommand(share)
	)

	root.AddCommand(shareCmd)
}

// NewKoolShare creates a new handler for sharing local environment with default dependencies
func NewKoolShare() *KoolShare {
	defaultKoolService := newDefaultKoolService()
	return &KoolShare{
		*defaultKoolService,
		&KoolShareFlags{"app", "", 0},
		environment.NewEnvStorage(),
		NewKoolStatus(),
		builder.NewCommand("docker", "run", "--rm", "--init"),
	}
}

func (s *KoolShare) validSubdomain(subdomain string) bool {
	return regexp.MustCompile("^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$").MatchString(subdomain)
}

// Execute runs the share logic.
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
	s.share.AppendArgs(s.Flags.parseServiceURI())
	s.share.AppendArgs("--server-host", "kool.live")

	if s.Flags.Subdomain != "" {
		s.Flags.Subdomain = strings.ToLower(s.Flags.Subdomain)
		if !s.validSubdomain(s.Flags.Subdomain) {
			err = fmt.Errorf("invalid subdomain '%s'", s.Flags.Subdomain)
			return
		}

		s.share.AppendArgs("--subdomain", s.Flags.Subdomain)
	}

	err = s.Shell().Interactive(s.share)
	return
}

// NewShareCommand initializes new kool share command
func NewShareCommand(share *KoolShare) (shareCmd *cobra.Command) {
	shareCmd = &cobra.Command{
		Use:   "share",
		Short: "Live share your local environment on the Internet using an HTTP tunnel",
		Args:  cobra.NoArgs,
		RunE:  DefaultCommandRunFunction(share),

		DisableFlagsInUseLine: true,
	}

	shareCmd.Flags().StringVarP(&share.Flags.Service, "service", "", "app", "The name of the local service container you want to share.")
	shareCmd.Flags().StringVarP(&share.Flags.Subdomain, "subdomain", "", "", "The subdomain used to generate your public https://subdomain.kool.live URL.")
	shareCmd.Flags().UintVarP(&share.Flags.Port, "port", "", 0, "The port from the target service that should be shared. If not provided, it will default to port 80.")
	return
}
