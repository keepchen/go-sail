package main

import (
	"github.com/keepchen/go-sail/cli/cmd/go-sail/generator"
	"github.com/spf13/cobra"
)

func main() {
	rootCMD.AddCommand(generatorCMD())
	_ = rootCMD.Execute()
}

var (
	goVersion   string
	workDir     string
	appName     string
	serviceName string
)

// rootCMD 根级(主)命令
var rootCMD = &cobra.Command{
	Use: "go-sail",
}

func generatorCMD() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen",
		Short: "代码生成",
		Run: func(cmd *cobra.Command, args []string) {
			//启动时要执行的操作写在这里
			generator.Gen(goVersion, workDir, appName, serviceName)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.PersistentFlags().StringVarP(&goVersion, "go version", "v", "", "go版本")
	cmd.PersistentFlags().StringVarP(&workDir, "directory", "d", ".", "工程路径(默认为当前文件夹)")
	cmd.PersistentFlags().StringVarP(&appName, "app name", "n", "", "应用名称")
	cmd.PersistentFlags().StringVarP(&serviceName, "service name", "s", "", "服务名称")

	return cmd
}
