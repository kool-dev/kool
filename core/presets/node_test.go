package presets

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/shell"

	"gopkg.in/yaml.v2"
)

func TestNodePresetRuntimeAndHelpers(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	previous := source
	t.Cleanup(func() {
		SetSource(previous)
		automate.SetRecipesSource(embed.FS{})
		if err := os.Chdir(cwd); err != nil {
			t.Error(err)
		}
	})
	root := os.DirFS(filepath.Join(cwd, "../..")).(SourceFS)
	SetSource(root)
	automate.SetRecipesSource(root)
	for _, ref := range []string{"javascript_package_manager", "php-version", "engine"} {
		t.Setenv(ref, "")
	}
	tests := []struct {
		preset   string
		services []string
		managers []string
	}{
		{"nodejs", []string{"app"}, []string{"npm", "yarn", "bun"}},
		{"expressjs", []string{"app"}, []string{"npm", "yarn", "bun"}},
		{"nextjs", []string{"app"}, []string{"npm", "yarn", "bun"}},
		{"nuxtjs", []string{"app"}, []string{"npm", "yarn", "bun"}},
		{"adonis", []string{"app"}, []string{"npm", "yarn", "bun"}},
		{"nestjs", []string{"app"}, []string{"npm"}},
		{"nest+next", []string{"backend", "frontend"}, []string{"npm"}},
		{"laravel", []string{"node"}, []string{"npm", "yarn", "bun"}},
		{"laravel+octane", []string{"node"}, []string{"npm", "yarn", "bun"}},
	}
	for _, test := range tests {
		for _, manager := range test.managers {
			t.Run(test.preset+"/"+manager, func(t *testing.T) {
				if err := os.Chdir(t.TempDir()); err != nil {
					t.Fatal(err)
				}
				p := NewParser()
				p.PrepareExecutor(&shell.FakeShell{}, automate.Answers{
					"javascript_package_manager": manager, "php-version": "PHP 8.4", "engine": "Swoole",
					"Which database service do you want to use": "None", "Which cache service do you want to use": "None",
				}, true)
				if err := p.Install(test.preset); err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile("docker-compose.yml")
				if err != nil {
					t.Fatal(err)
				}
				var compose struct {
					Services map[string]struct {
						Image string `yaml:"image"`
					} `yaml:"services"`
				}
				if err := yaml.Unmarshal(data, &compose); err != nil {
					t.Fatal(err)
				}
				want := "kooldev/node:24"
				if manager == "bun" && test.preset != "nodejs" {
					want = "oven/bun:1"
				}
				for _, service := range test.services {
					if got := compose.Services[service].Image; got != want {
						t.Errorf("service %s image = %q, want %q", service, got, want)
					}
				}
				data, err = os.ReadFile("kool.yml")
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(data), "kooldev/node:20") {
					t.Fatal("generated helpers still use Node 20")
				}
				if manager != "bun" && !strings.Contains(string(data), "kooldev/node:24") {
					t.Fatal("generated npm/Yarn helpers do not use Node 24")
				}
			})
		}
	}
}
