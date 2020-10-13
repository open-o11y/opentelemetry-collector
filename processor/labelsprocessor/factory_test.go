package labelsprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcheck"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exportertest"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	pType := factory.Type()
	assert.Equal(t, pType, configmodels.Type("labels_processor"))
}

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.Equal(t, cfg, &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			NameVal: typeStr,
			TypeVal: typeStr,
		},
	})
	assert.NoError(t, configcheck.ValidateConfig(cfg))
}

func TestCreateProcessor(t *testing.T) {

	factory := NewFactory()
	cfg := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "label1", Value: "value1"},
		},
	}

	mp, mErr := factory.CreateMetricsProcessor(context.Background(), component.ProcessorCreateParams{Logger: zap.NewNop()}, cfg, exportertest.NewNopMetricsExporter())
	assert.Equal(t, true, mp != nil)
	assert.Equal(t, mErr, nil)

}
