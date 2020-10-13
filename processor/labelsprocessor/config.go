package labelsprocessor

import "go.opentelemetry.io/collector/config/configmodels"

// Config defines configuration for Labels processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
	Labels                         []LabelConfig `mapstructure:"labels"`
}

// LabelConfig defines configuration for provided labels
type LabelConfig struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}
