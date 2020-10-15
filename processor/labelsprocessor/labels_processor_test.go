package labelsprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"

	otlpmetrics "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
)

func TestValidateConfig(t *testing.T) {
	missingValConfig := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster"},
			{Key: "__replica__", Value: "r1"},
		},
	}

	emptyValConfig := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster", Value: ""},
			{Key: "__replica__", Value: "r1"},
		},
	}

	duplicateKeyConfig := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster", Value: "c1"},
			{Key: "__replica__", Value: "r1"},
			{Key: "cluster", Value: "c2"},
		},
	}

	validCfg := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster", Value: "c1"},
			{Key: "__replica__", Value: "r1"},
		},
	}

	assert.Error(t, validateConfig(missingValConfig))
	assert.Error(t, validateConfig(emptyValConfig))
	assert.Error(t, validateConfig(duplicateKeyConfig))
	assert.Nil(t, validateConfig(validCfg))
}

func TestAttachLabels(t *testing.T) {
	cfg := &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: "labels_processor",
			NameVal: "labels_processor",
		},
		Labels: []LabelConfig{
			{Key: "cluster", Value: "c1"},
			{Key: "__replica__", Value: "r1"},
		},
	}

	// Set up tests
	tests := []struct {
		name                   string
		config                 *Config
		incomingMetric         []*otlpmetrics.ResourceMetrics
		expectedOutgoingMetric []*otlpmetrics.ResourceMetrics
	}{
		{
			name:                   "test_add_labels_on_all_metrics_types",
			config:                 cfg,
			incomingMetric:         metricAllDataTypes,
			expectedOutgoingMetric: expectedMetricAllDataTypes,
		},
		{
			name:                   "test_adding_duplicate_labels",
			config:                 cfg,
			incomingMetric:         metricDuplicateLabels,
			expectedOutgoingMetric: expectedMetricDuplicateLabels,
		},
	}

	// Execute tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory()
			tmn := &testMetricsConsumer{}
			rmp, err := factory.CreateMetricsProcessor(context.Background(), component.ProcessorCreateParams{}, cfg, tmn)
			require.NoError(t, err)
			assert.Equal(t, true, rmp.GetCapabilities().MutatesConsumedData)

			sourceMetricData := pdata.MetricsFromOtlp(tt.incomingMetric)
			wantMetricData := pdata.MetricsFromOtlp(tt.expectedOutgoingMetric)
			err = rmp.ConsumeMetrics(context.Background(), sourceMetricData)
			require.NoError(t, err)
			assert.EqualValues(t, wantMetricData, tmn.md)
		})
	}

}

type testMetricsConsumer struct {
	md pdata.Metrics
}

func (tmn *testMetricsConsumer) ConsumeMetrics(_ context.Context, md pdata.Metrics) error {
	// simply store md in struct so that we can compare the result to the input
	tmn.md = md
	return nil
}
