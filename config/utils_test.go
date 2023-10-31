package config

import "testing"

func TestPrintConfig(t *testing.T) {
	formats := [...]string{"json", "toml", "yaml", "unknown"}
	for _, format := range formats {
		PrintTemplateConfig(format)
	}
	t.Log("OK")
}
