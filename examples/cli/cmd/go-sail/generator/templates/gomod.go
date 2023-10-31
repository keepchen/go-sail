package templates

var GoModTpl = `module {{ .AppName }}

go {{ .GoVersion }}
`
