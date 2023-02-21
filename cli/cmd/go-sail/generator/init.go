package generator

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/keepchen/go-sail/cli/cmd/go-sail/generator/templates"
)

func registerFiles4Gen(opts Options) []File {
	return []File{
		//~/workdir/appName/go.mod
		{Path: fmt.Sprintf("%s/go.mod", opts.AppName), Template: templates.GoModTpl},
		//~/workdir/appName/main.go
		{Path: fmt.Sprintf("%s/main.go", opts.AppName), Template: templates.MainTpl},
		//~/workdir/appName/Makefile
		{Path: fmt.Sprintf("%s/Makefile", opts.AppName), Template: templates.MakefileTpl},
		//~/workdir/appName/Dockerfile
		{Path: fmt.Sprintf("%s/Dockerfile", opts.AppName), Template: templates.DockerfileTpl},
		//~/workdir/appName/cmd
		{Path: fmt.Sprintf("%s/cmd", opts.AppName), Template: ""},
		//~/workdir/appName/cmd/root.go
		{Path: fmt.Sprintf("%s/cmd/root.go", opts.AppName), Template: templates.CmdRootTpl},
		//~/workdir/appName/cmd/serviceName.go
		{Path: fmt.Sprintf("%s/cmd/%s.go", opts.AppName, opts.ServiceName), Template: templates.CmdSrvTpl},
		//~/workdir/appName/logs
		{Path: fmt.Sprintf("%s/logs", opts.AppName), Template: ""},
		//~/workdir/appName/logs/.gitignore
		{Path: fmt.Sprintf("%s/logs/.gitignore", opts.AppName), Template: templates.GitIgnoreAllTpl},
		//~/workdir/appName/pkg/serviceName.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/%s.go", opts.AppName, opts.ServiceName, opts.ServiceName), Template: templates.AppRoorSrvTpl},
		//~/workdir/appName/pkg/serviceName/config
		{Path: fmt.Sprintf("%s/pkg/app/%s/config", opts.AppName, opts.ServiceName), Template: ""},
		//~/workdir/appName/pkg/serviceName/app/config/config.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/config.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigTpl},
		//~/workdir/appName/pkg/serviceName/app/config/config.sample.toml
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/config.sample.toml", opts.AppName, opts.ServiceName), Template: templates.AppConfigTomlTpl},
		//~/workdir/appName/pkg/app/serviceName/config/parser.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/parser.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigParserTpl},
		//~/workdir/appName/pkg/app/serviceName/config/utils.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/utils.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigUtilsTpl},
		//~/workdir/appName/pkg/app/common
		{Path: fmt.Sprintf("%s/pkg/common", opts.AppName), Template: ""},
		//~/workdir/appName/pkg/app/constants
		{Path: fmt.Sprintf("%s/pkg/constants", opts.AppName), Template: ""},
		//~/workdir/appName/pkg/app/lib
		{Path: fmt.Sprintf("%s/pkg/lib", opts.AppName), Template: ""},
		//~/workdir/appName/pkg/app/utils
		{Path: fmt.Sprintf("%s/pkg/utils", opts.AppName), Template: ""},
	}
}

func Init(goVersion, workDir, appName, serviceName string) {
	if len(goVersion) == 0 {
		goVersion = runtime.Version()
		log.Printf("[!] Go version is empty,use current version: %s\n", goVersion)
	}
	if len(workDir) == 0 {
		workDir, _ = os.Getwd()
		log.Printf("[!] Work directory is empty,use current directory: %s\n", workDir)
	}
	if len(appName) == 0 {
		log.Println("[x] App name is empty,process exited.")
		return
	}
	if len(serviceName) == 0 {
		serviceName = "demo"
		log.Println("[!] Service name is empty,use default name: demo")
	}
	opts := Options{
		GoVersion:   goVersion,
		WorkDir:     workDir,
		AppName:     appName,
		ServiceName: serviceName,
	}
	err := New(opts).Generate(registerFiles4Gen(opts))
	if err != nil {
		log.Printf("[!] Process finished,error occurred: %s", err.Error())
	} else {
		log.Printf("[√] Proccess finished,create application success! :)")
		log.Printf(`
-------------------------------------------------------------------------------------------
Start service steps:
[1/3] cd %s
[2/3] go mod tidy
[3/3] go run main.go %s -c pkg/app/%s/config/config.sample.toml

Enjoy it. :)
`, workDir, serviceName, serviceName)
	}
}
