package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type versionImpl struct {
}

type IVersion interface {
	// Print prints version information
	Print(fields VersionInfoFields)
}

var vei IVersion = versionImpl{}

// Version instantiates version utilities
func Version() IVersion {
	return vei
}

// VersionInfoFields version information fields
type VersionInfoFields struct {
	AppName   string //Application name
	Branch    string //Version branch
	Version   string //Version number
	Revision  string //commit id
	BuildDate string //Build time
	GoVersion string //Go version
}

const versionInfoTmpl = `
{{.program}}, version: {{.version}} 
(branch: {{.branch}}; revision: {{.revision}})
  build date:   {{.buildDate}}
  go version:   {{.goVersion}}
`

// Print prints version information
func (versionImpl) Print(fields VersionInfoFields) {
	m := map[string]string{
		"program":   fields.AppName,
		"version":   fields.Version,
		"branch":    fields.Branch,
		"revision":  fields.Revision,
		"buildDate": fields.BuildDate,
		"goVersion": fields.GoVersion,
	}

	tmpl, err := template.Must(template.New("version"), nil).Parse(versionInfoTmpl)
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	if exeErr := tmpl.Execute(&buf, m); exeErr != nil {
		panic(exeErr)
	}

	fmt.Println(strings.TrimSpace(buf.String()))
}
