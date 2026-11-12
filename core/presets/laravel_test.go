package presets

import (
	"embed"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"kool-dev/kool/core/automate"
	"kool-dev/kool/core/shell"

	"gopkg.in/yaml.v2"
)

func TestLaravelPreservesViteConfiguration(t *testing.T) {
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
	for _, ref := range []string{"php-version", "engine", "javascript_package_manager"} {
		t.Setenv(ref, "")
	}
	configs := []struct {
		name, filename, content string
	}{
		{"laravel13", "vite.config.js", "import { bunny } from 'laravel-vite-plugin/fonts';\nimport tailwindcss from '@tailwindcss/vite';\nexport default { plugins: [laravel({ fonts: [bunny('Instrument Sans', { weights: [400, 500, 600] })] }), tailwindcss()], server: { watch: { ignored: ['**/storage/framework/views/**'] } } };\n"},
		{"laravel12", "vite.config.js", "import tailwindcss from '@tailwindcss/vite';\nexport default { plugins: [laravel(), tailwindcss()], server: { watch: { ignored: ['**/storage/framework/views/**'] } } };\n"},
		{"legacy", "vite.config.js", "export default { plugins: [laravel()], resolve: { alias: { '@': '/resources/js' } } };\n"},
		{"typescript", "vite.config.ts", "export default defineConfig({ plugins: [laravel(), customPlugin()] });\n"},
	}
	commands := map[string][]string{
		"npm":  {"npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3001", "--strictPort"},
		"yarn": {"yarn", "dev", "--host", "0.0.0.0", "--port", "3001", "--strictPort"},
		"bun":  {"bun", "--bun", "run", "dev", "--host", "0.0.0.0", "--port", "3001", "--strictPort"},
	}
	for _, preset := range []string{"laravel", "laravel+octane"} {
		for _, config := range configs {
			for manager, command := range commands {
				t.Run(preset+"/"+config.name+"/"+manager, func(t *testing.T) {
					dir := t.TempDir()
					if err := os.Chdir(dir); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(config.filename, []byte(config.content), 0644); err != nil {
						t.Fatal(err)
					}
					p := NewParser()
					p.PrepareExecutor(&shell.FakeShell{}, automate.Answers{
						"php-version": "PHP 8.4", "engine": "Swoole", "javascript_package_manager": manager,
						"Which database service do you want to use": "None", "Which cache service do you want to use": "None",
					}, true)
					if err := p.Install(preset); err != nil {
						t.Fatal(err)
					}
					data, err := os.ReadFile(config.filename)
					if err != nil || string(data) != config.content {
						t.Fatalf("Vite configuration changed: %q (%v)", data, err)
					}
					files, err := filepath.Glob("vite.config.*")
					if err != nil || len(files) != 1 || files[0] != config.filename {
						t.Fatalf("unexpected Vite configuration/backup files: %v (%v)", files, err)
					}
					data, err = os.ReadFile("docker-compose.yml")
					if err != nil {
						t.Fatal(err)
					}
					var compose struct {
						Services struct {
							Node struct {
								Command []string `yaml:"command"`
								Ports   []string `yaml:"ports"`
							} `yaml:"node"`
						} `yaml:"services"`
					}
					if err := yaml.Unmarshal(data, &compose); err != nil {
						t.Fatal(err)
					}
					if node := compose.Services.Node; !reflect.DeepEqual(node.Command, command) || !reflect.DeepEqual(node.Ports, []string{"3001:3001"}) {
						t.Fatalf("Vite command/port mismatch: %+v", node)
					}
				})
			}
		}
	}
}
