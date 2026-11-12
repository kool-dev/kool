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

func TestJavaPresetsStandalone(t *testing.T) {
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
	SetSource(os.DirFS(filepath.Join(cwd, "../..")).(SourceFS))
	automate.SetRecipesSource(os.DirFS(filepath.Join(cwd, "../..")).(automate.RecipeSource))
	t.Setenv("JAVA_VERSION", "")
	t.Setenv("JAVA_DB", "")
	t.Setenv("JAVA_CACHE", "")
	for _, preset := range []string{"java", "spring-boot"} {
		for _, version := range []string{"17", "21", "25", "7"} {
			t.Run(preset+"/"+version, func(t *testing.T) {
				if err := os.Chdir(t.TempDir()); err != nil {
					t.Fatal(err)
				}
				p := NewParser()
				p.PrepareExecutor(&shell.FakeShell{}, automate.Answers{
					"JAVA_VERSION": version, "JAVA_DB": "None", "JAVA_CACHE": "None",
				}, true)
				err := p.Install(preset)
				if version == "7" {
					if err == nil || !strings.Contains(err.Error(), "invalid option") {
						t.Fatalf("expected unsupported JDK error, got %v", err)
					}
					if _, err := os.Stat("docker-compose.yml"); !os.IsNotExist(err) {
						t.Fatal("invalid JDK wrote configuration")
					}
					return
				}
				if err != nil {
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
					Networks map[string]interface{} `yaml:"networks"`
				}
				if err := yaml.Unmarshal(data, &compose); err != nil {
					t.Fatal(err)
				}
				if compose.Services["app"].Image != "gradle:9-jdk"+version {
					t.Fatalf("wrong image: %s", compose.Services["app"].Image)
				}
				for _, network := range []string{"kool_local", "kool_global"} {
					if _, ok := compose.Networks[network]; !ok {
						t.Fatalf("undefined network %s", network)
					}
				}
			})
		}
	}
}
