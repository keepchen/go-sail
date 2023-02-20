package templates

var MainTpl = `package main

import (
	"{{ .AppName }}/cmd"
)

func main() {
	_ = cmd.RootCMD.Execute()
}
`
