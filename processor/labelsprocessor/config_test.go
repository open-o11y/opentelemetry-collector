package labelsprocessor

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestLoadConfig(t *testing.T) {

	factories, err := componenttest.ExampleComponents()
	assert.NoError(t, err)

	factories.Processors[typeStr] = NewFactory()

	os.Setenv("VALUE_1", "first_val")
	os.Setenv("VALUE_2", "second_val")

	cfg, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, cfg.Processors["labels_processor"], &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "label1", Value: "value1"},
			{Key: "label2", Value: "value2"},
		},
	})

	assert.Equal(t, cfg.Processors["labels_processor/env_vars"], &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor/env_vars",
		},
		Labels: []LabelConfig{
			{Key: "label1", Value: "first_val"},
			{Key: "label2", Value: "second_val"},
		},
	})

}
