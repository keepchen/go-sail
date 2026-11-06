package codec

import (
	"fmt"

	"github.com/DataDog/sketches-go/ddsketch"
)

// DDSketch wraps DataDog's mergeable quantile Sketch and its accuracy configuration.
type DDSketch struct {
	value    *ddsketch.DDSketch
	accuracy float64
}

// NewDDSketch creates an empty Sketch with the specified relative accuracy.
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

// Add inserts a value into the Sketch.
func (d *DDSketch) Add(value float64) error {
	if err := d.value.Add(value); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> add ddsketch value: %w", err)
	}
	return nil
}

// Merge merges all samples from another compatible Sketch.
func (d *DDSketch) Merge(other *DDSketch) error {
	if other == nil {
		return nil
	}
	if err := d.value.MergeWith(other.value); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> merge ddsketch: %w", err)
	}
	return nil
}

// MergeEncoded decodes and merges a payload whose mapping configuration is compatible with the current instance.
func (d *DDSketch) MergeEncoded(payload []byte) error {
	if err := d.value.DecodeAndMergeWith(payload); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> decode and merge ddsketch: %w", err)
	}
	return nil
}

// Quantiles returns approximate values for the specified quantiles.
func (d *DDSketch) Quantiles(values []float64) ([]float64, error) {
	result, err := d.value.GetValuesAtQuantiles(values)
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> calculate ddsketch quantiles: %w", err)
	}
	return result, nil
}

// Encode serializes the Sketch, including the index mapping for compatibility checks.
func (d *DDSketch) Encode() []byte {
	payload := make([]byte, 0, 512)

	// The first version includes Mapping so the decoder can verify compatibility.
	d.value.Encode(&payload, false)

	return payload
}

// Clear removes all samples while retaining the Sketch configuration.
func (d *DDSketch) Clear() {
	d.value.Clear()
}

// Empty indicates whether the Sketch contains no samples.
func (d *DDSketch) Empty() bool {
	return d.value.IsEmpty()
}

// Accuracy returns the configured relative accuracy.
func (d *DDSketch) Accuracy() float64 {
	return d.accuracy
}
