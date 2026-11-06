package cmd

import (
	"github.com/spf13/cobra"
)

var cfgPath string

//RootCMD root (main) command
var RootCMD = &cobra.Command{
	Use: "go-sail",
}
