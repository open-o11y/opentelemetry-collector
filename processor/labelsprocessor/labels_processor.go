package labelsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer/pdata"
	v11 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	v1 "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/metrics/v1"
)

type labelMetricProcessor struct {
	cfg *Config
}

func newLabelMetricProcessor(cfg *Config) (*labelMetricProcessor, error) {
	return &labelMetricProcessor{cfg: cfg}, nil
}

func (lp *labelMetricProcessor) ProcessMetrics(_ context.Context, md pdata.Metrics) (pdata.Metrics, error) {

	otlpMetrics := pdata.MetricsToOtlp(md)

	for _, otlpMetric := range otlpMetrics {
		for _, instrMetric := range otlpMetric.GetInstrumentationLibraryMetrics() {
			for _, metric := range instrMetric.GetMetrics() {

				var intDataPoint []*v1.IntDataPoint
				var doubleDataPoint []*v1.DoubleDataPoint
				var intHistogramDataPoint []*v1.IntHistogramDataPoint
				var doubleHistogramDataPoint []*v1.DoubleHistogramDataPoint

				if metric.GetDoubleGauge() != nil {
					doubleDataPoint = metric.GetDoubleGauge().GetDataPoints()
				} else if metric.GetDoubleHistogram() != nil {
					doubleHistogramDataPoint = metric.GetDoubleHistogram().GetDataPoints()
				} else if metric.GetDoubleSum() != nil {
					doubleDataPoint = metric.GetDoubleSum().GetDataPoints()
				} else if metric.GetIntGauge() != nil {
					intDataPoint = metric.GetIntGauge().GetDataPoints()
				} else if metric.GetIntHistogram() != nil {
					intHistogramDataPoint = metric.GetIntHistogram().GetDataPoints()
				} else if metric.GetIntSum() != nil {
					intDataPoint = metric.GetIntSum().GetDataPoints()
				}

				for _, label := range lp.cfg.Labels {
					for _, dataPoint := range intDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range doubleDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range intHistogramDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
					for _, dataPoint := range doubleHistogramDataPoint {
						deDuplicateAndAppend(&dataPoint.Labels, label.Key, label.Value)
					}
				}

			}
		}
	}

	return md, nil
}

// This processor will always by default update existing label values. Also assumes duplicate labels do not already exist in the metric
func deDuplicateAndAppend(labels *[]*v11.StringKeyValue, key string, value string) {
	// If the key already exists, overwrite it
	for _, elem := range *labels {
		if elem.Key == key {
			elem.Value = value
			return
		}
	}
	// If it does not exist, append it
	*labels = append(*labels, &v11.StringKeyValue{Key: key, Value: value})
}
