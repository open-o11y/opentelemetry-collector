package labelsprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/consumer/pdata"
	v11 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	v1 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
)

type labelMetricProcessor struct {
	cfg *Config
}

func newLabelMetricProcessor(cfg *Config) (*labelMetricProcessor, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &labelMetricProcessor{cfg: cfg}, nil
}

func validateConfig(cfg *Config) error {
	for _, elem := range cfg.Labels {
		if elem.Key == "" || elem.Value == "" {
			return fmt.Errorf("Labels Processor configuration contains an empty key or value")
		}
	}

	keys := make(map[string]bool)
	for _, elem := range cfg.Labels {
		_, value := keys[elem.Key]
		if value {
			return fmt.Errorf("Labels Processor configuration contains duplicate keys")
		}
		keys[elem.Key] = true
	}

	return nil
}

func (lp *labelMetricProcessor) ProcessMetrics(_ context.Context, md pdata.Metrics) (pdata.Metrics, error) {

	otlpMetrics := pdata.MetricsToOtlp(md)

	for _, otlpMetric := range otlpMetrics {
		for _, instrMetric := range otlpMetric.GetInstrumentationLibraryMetrics() {
			for _, metric := range instrMetric.GetMetrics() {

				// Multiple types of Data Points exists, and each of them must be handled differently
				if metric.GetIntSum() != nil {
					intDataPoints := metric.GetIntSum().GetDataPoints()
					handleIntDataPoints(intDataPoints, lp)
				} else if metric.GetIntGauge() != nil {
					intDataPoints := metric.GetIntGauge().GetDataPoints()
					handleIntDataPoints(intDataPoints, lp)
				} else if metric.GetDoubleGauge() != nil {
					doubleDataPoints := metric.GetDoubleGauge().GetDataPoints()
					handleDoubleDataPoints(doubleDataPoints, lp)
				} else if metric.GetDoubleSum() != nil {
					doubleDataPoints := metric.GetDoubleSum().GetDataPoints()
					handleDoubleDataPoints(doubleDataPoints, lp)
				} else if metric.GetIntHistogram() != nil {
					intHistogramDataPoints := metric.GetIntHistogram().GetDataPoints()
					handleIntHistogramDataPoints(intHistogramDataPoints, lp)
				} else if metric.GetDoubleHistogram() != nil {
					doubleHistogramDataPoints := metric.GetDoubleHistogram().GetDataPoints()
					handleDoubleHistogramDataPoints(doubleHistogramDataPoints, lp)
				}

			}
		}
	}

	return md, nil
}

func handleIntDataPoints(intDataPoints []*v1.IntDataPoint, lp *labelMetricProcessor) {
	for _, dataPoint := range intDataPoints {
		upsertLabels(&dataPoint.Labels, lp)
	}
}

func handleDoubleDataPoints(doubleDataPoints []*v1.DoubleDataPoint, lp *labelMetricProcessor) {
	for _, dataPoint := range doubleDataPoints {
		upsertLabels(&dataPoint.Labels, lp)
	}
}

func handleIntHistogramDataPoints(intHistogramDataPoints []*v1.IntHistogramDataPoint, lp *labelMetricProcessor) {
	for _, dataPoint := range intHistogramDataPoints {
		upsertLabels(&dataPoint.Labels, lp)
	}
}

func handleDoubleHistogramDataPoints(doubleHistogramDataPoints []*v1.DoubleHistogramDataPoint, lp *labelMetricProcessor) {
	for _, dataPoint := range doubleHistogramDataPoints {
		upsertLabels(&dataPoint.Labels, lp)
	}

}

func upsertLabels(labels *[]*v11.StringKeyValue, lp *labelMetricProcessor) {
	for _, label := range lp.cfg.Labels {
		var updated bool = false
		for _, elem := range *labels {
			if elem.Key == label.Key {
				elem.Value = label.Value
				updated = true
				break
			}
		}
		if !updated {
			*labels = append(*labels, &v11.StringKeyValue{Key: label.Key, Value: label.Value})
		}
	}
}
