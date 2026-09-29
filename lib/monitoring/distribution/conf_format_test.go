package distribution

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

type exampleConfig struct {
	Distribution Conf `yaml:"distribution" toml:"distribution" json:"distribution"`
}

func TestConfigExamplesAreEquivalent(t *testing.T) {
	files := []struct {
		name      string
		unmarshal func([]byte, any) error
	}{
		{name: "config.example.yaml", unmarshal: yaml.Unmarshal},
		{name: "config.example.toml", unmarshal: toml.Unmarshal},
		{name: "config.example.json", unmarshal: json.Unmarshal},
	}

	var expected Conf
	for index, item := range files {
		data, err := os.ReadFile(item.name)
		if err != nil {
			t.Fatalf("读取 %s: %v", item.name, err)
		}

		var example exampleConfig
		if err := item.unmarshal(data, &example); err != nil {
			t.Fatalf("解析 %s: %v", item.name, err)
		}
		if err := example.Distribution.Validate(); err != nil {
			t.Fatalf("校验 %s: %v", item.name, err)
		}

		if index == 0 {
			expected = example.Distribution
			continue
		}
		if !reflect.DeepEqual(example.Distribution, expected) {
			t.Fatalf("%s 与 YAML 示例不一致\n实际值: %#v\n期望值: %#v", item.name, example.Distribution, expected)
		}
	}
}
