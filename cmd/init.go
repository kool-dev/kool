package cmd

import "github.com/gookit/color"

func init() {
	initCmd := NewPresetCommand(NewKoolPreset())
	initCmd.Use = "init [PRESET]"
	initCmd.Short = "[DEPRECATED] Proxies preset command"
	initCmd.Deprecated = color.New(color.Yellow).Sprint("use the \"preset\" command instead.")

	rootCmd.AddCommand(initCmd)
}
