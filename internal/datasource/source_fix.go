package datasource

import "context"

var (
	_ Source = &FixSource{}
	_ Config = &FixSourceConfig{}
)

// FixSource defines a source with fixed data.
type FixSource struct {
	data uint64
}

// NewFixSource returns a new fix source.
func NewFixSource(data uint64) *FixSource {
	return &FixSource{
		data: data,
	}
}

// GetName returns the source name.
func (s FixSource) GetName() string {
	return "FixSource"
}

// GetData returns the fixed data.
func (s FixSource) GetData(ctx context.Context) (uint64, error) {
	return s.data, nil
}

// FixSourceConfig defines the configuration for the fix source.
type FixSourceConfig struct {
	SourceType SourceType `mapstructure:"source_type" toml:"source_type"`
	Data       uint64     `mapstructure:"data"        toml:"data"`
}

// Validate validates the fix source configuration.
func (c FixSourceConfig) Validate() error {
	return nil
}

// NewSource creates a new FixSource object.
func (c FixSourceConfig) NewSource() (Source, error) {
	return NewFixSource(c.Data), nil
}
