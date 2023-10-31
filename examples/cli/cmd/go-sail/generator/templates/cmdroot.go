package templates

// CmdRootTpl 根指令模板
var CmdRootTpl = `package cmd

import (
	"github.com/spf13/cobra"
)

var cfgPath string

//RootCMD 根级(主)命令
var RootCMD = &cobra.Command{
	Use: "{{ .AppName }}",
}
`
