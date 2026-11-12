package presets

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/builder"
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
				database := "None"
				if version == "21" {
					database = "MongoDB"
				}
				p.PrepareExecutor(&shell.FakeShell{}, automate.Answers{
					"JAVA_VERSION": version, "JAVA_DB": database, "JAVA_CACHE": "None",
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
					Volumes  map[string]interface{} `yaml:"volumes"`
				}
				if err := yaml.Unmarshal(data, &compose); err != nil {
					t.Fatal(err)
				}
				if compose.Services["app"].Image != "gradle:9-jdk"+version {
					t.Fatalf("wrong image: %s", compose.Services["app"].Image)
				}
				if database == "MongoDB" {
					if _, ok := compose.Volumes["database"]; !ok {
						t.Fatal("MongoDB refers to an undefined database volume")
					}
					if compose.Services["mongodb"].Image == "" {
						t.Fatal("MongoDB service missing")
					}
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

func TestJavaCreationScriptsKeepUserInputInEnvironment(t *testing.T) {
	t.Setenv("CREATE_DIRECTORY", "my app")
	t.Setenv("JAVA_VERSION", "21")
	t.Setenv("SPRING_NAME", `Demo App; $(touch should-not-run)`)
	for _, recipe := range []string{"create-java", "create-spring-boot"} {
		data, err := os.ReadFile(filepath.Join("../..", "recipes", recipe+".yml"))
		if err != nil {
			t.Fatal(err)
		}
		var set automate.ActionSet
		if err := yaml.Unmarshal(data, &set); err != nil {
			t.Fatal(err)
		}
		scripts := set.Actions[len(set.Actions)-1].Scripts
		mkdir, err := builder.ParseCommand(scripts[0])
		if err != nil || len(mkdir.Args()) != 2 || mkdir.Args()[1] != "my app" {
			t.Fatalf("directory name split: %v, %v", mkdir, err)
		}
		command, err := builder.ParseCommand(scripts[len(scripts)-1])
		if err != nil {
			t.Fatal(err)
		}
		args := command.Args()
		if len(args) < 3 || args[len(args)-2] != "-c" || !strings.Contains(args[len(args)-1], `"$(printenv CREATE_DIRECTORY)"`) {
			t.Fatalf("expected environment-backed shell command: %v", args)
		}
		if recipe == "create-spring-boot" {
			if strings.Contains(command.String(), "should-not-run") || !strings.Contains(command.String(), "--env SPRING_NAME") {
				t.Fatalf("application name interpolated into shell syntax: %s", command.String())
			}
		}
	}
}
