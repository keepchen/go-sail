package codec

import (
	"fmt"

	"github.com/DataDog/sketches-go/ddsketch"
)

// DDSketch 封装 DataDog 的可合并分位数 Sketch 及其精度配置。
type DDSketch struct {
	value    *ddsketch.DDSketch
	accuracy float64
}

// NewDDSketch 按指定相对精度创建空 Sketch。
func NewDDSketch(relativeAccuracy float64) (*DDSketch, error) {
	value, err := ddsketch.NewDefaultDDSketch(relativeAccuracy)
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> create ddsketch: %w", err)
	}

	return &DDSketch{
		value:    value,
		accuracy: relativeAccuracy,
	}, nil
}

// Add 向 Sketch 插入一个数值。
func (d *DDSketch) Add(value float64) error {
	if err := d.value.Add(value); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> add ddsketch value: %w", err)
	}
	return nil
}

// Merge 将另一个兼容 Sketch 中的所有样本合并进来。
func (d *DDSketch) Merge(other *DDSketch) error {
	if other == nil {
		return nil
	}
	if err := d.value.MergeWith(other.value); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> merge ddsketch: %w", err)
	}
	return nil
}

// MergeEncoded 解码并合并映射配置与当前实例兼容的载荷。
func (d *DDSketch) MergeEncoded(payload []byte) error {
	if err := d.value.DecodeAndMergeWith(payload); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> decode and merge ddsketch: %w", err)
	}
	return nil
}

// Quantiles 返回指定分位点对应的近似值。
func (d *DDSketch) Quantiles(values []float64) ([]float64, error) {
	result, err := d.value.GetValuesAtQuantiles(values)
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> calculate ddsketch quantiles: %w", err)
	}
	return result, nil
}

// Encode 序列化 Sketch，并携带索引映射以便进行兼容性检查。
func (d *DDSketch) Encode() []byte {
	payload := make([]byte, 0, 512)

	// 第一版携带 Mapping，保证解码端能够验证兼容性。
	d.value.Encode(&payload, false)

	return payload
}

// Clear 清除全部样本，但保留 Sketch 配置。
func (d *DDSketch) Clear() {
	d.value.Clear()
}

// Empty 表示 Sketch 是否不含任何样本。
func (d *DDSketch) Empty() bool {
	return d.value.IsEmpty()
}

// Accuracy 返回配置的相对精度。
func (d *DDSketch) Accuracy() float64 {
	return d.accuracy
}
