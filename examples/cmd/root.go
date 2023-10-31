package cmd

import (
	"github.com/spf13/cobra"
)

var cfgPath string

//RootCMD 根级(主)命令
var RootCMD = &cobra.Command{
	Use: "go-sail",
}
