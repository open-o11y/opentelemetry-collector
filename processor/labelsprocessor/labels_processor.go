package labelsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer/pdata"
)

type labelMetricProcessor struct {
	cfg *Config
}

func newLabelMetricProcessor(cfg *Config) (*labelMetricProcessor, error) {
	return &labelMetricProcessor{cfg: cfg}, nil
}

func (lp *labelMetricProcessor) ProcessMetrics(_ context.Context, md pdata.Metrics) (pdata.Metrics, error) {
	return md, nil
}
