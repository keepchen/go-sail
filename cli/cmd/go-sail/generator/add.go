package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/keepchen/go-sail/cli/cmd/go-sail/generator/templates"
)

func registerFiles4Add(opts Options) []File {
	return []File{
		//~/workdir/appName/cmd/serviceName.go
		{Path: fmt.Sprintf("%s/cmd/%s.go", opts.AppName, opts.ServiceName), Template: templates.CmdSrvTpl},
		//~/workdir/appName/pkg/serviceName.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/%s.go", opts.AppName, opts.ServiceName, opts.ServiceName), Template: templates.AppRoorSrvTpl},
		//~/workdir/appName/pkg/serviceName/config
		{Path: fmt.Sprintf("%s/pkg/app/%s/config", opts.AppName, opts.ServiceName), Template: ""},
		//~/workdir/appName/pkg/app/serviceName/config/config.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/config.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigTpl},
		//~/workdir/appName/pkg/app/serviceName/config/config.sample.toml
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/config.sample.toml", opts.AppName, opts.ServiceName), Template: templates.AppConfigTomlTpl},
		//~/workdir/appName/pkg/app/serviceName/config/parser.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/parser.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigParserTpl},
		//~/workdir/appName/pkg/app/serviceName/config/utils.go
		{Path: fmt.Sprintf("%s/pkg/app/%s/config/utils.go", opts.AppName, opts.ServiceName), Template: templates.AppConfigUtilsTpl},
	}
}

func Add(goVersion, workDir, appName, serviceName string) {
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
	if _, err := os.Stat(filepath.Join(workDir, appName)); os.IsNotExist(err) {
		log.Println("[x] App work directory is not exist,process exited.")
		log.Println("---------------------------------------------------")
		log.Println("Try `go-sail init` command first.")
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
	err := New(opts).Generate(registerFiles4Add(opts))
	if err != nil {
		log.Printf("[!] Process finished,error occurred: %s", err.Error())
	} else {
		log.Printf("[√] Proccess finished,add service success! :)")
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
