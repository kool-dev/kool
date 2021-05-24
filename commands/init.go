package commands

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func AddKoolInit(root *cobra.Command) {
	initCmd := NewPresetCommand(NewKoolPreset())
	initCmd.Use = "init [PRESET]"
	initCmd.Short = "[DEPRECATED] Proxies preset command"
	initCmd.Deprecated = color.New(color.Yellow).Sprint("use the \"preset\" command instead.")

	root.AddCommand(initCmd)
}
