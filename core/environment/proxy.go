package environment

import (
	"kool-dev/kool/core/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/compose-spec/compose-go/template"
)

func initProxy(envStorage EnvStorage, workDir string) {
	var config *parser.KoolYaml
	var err error
	for _, name := range []string{"kool.yml", "kool.yaml"} {
		if config, err = parser.ParseKoolYaml(filepath.Join(workDir, name)); err == nil {
			break
		}
	}
	if err != nil || config.Proxy == nil {
		return
	}

	domain, err := template.Substitute(config.Proxy.Domain, os.LookupEnv)
	if err != nil {
		return
	}
	domain = strings.TrimSuffix(strings.TrimSpace(domain), ".")
	if domain == "" {
		return
	}

	envStorage.Set("KOOL_PROXY_DOMAIN", domain)
	host := domain
	if envStorage.IsTrue("KOOL_WORKSPACE") {
		host = envStorage.Get("KOOL_WORKSPACE_NAME") + ".workspace." + domain
	}
	envStorage.Set("KOOL_PROXY_HOST", host)
}
